package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"time"
)

/**
	Иерархия html документа, по которой производится поиск

	div class=post bg1 or class=post bg2
	div class=inner
	div class=postbody
	div class=content <<- здесь лежит нужное сообщение и иногда лишний html код (например, картинка)
 */

var App Config

type Config struct {
	USR     string
	PASS    string
	Delay   int64
	MailArr []string
}

var alreadyParsedPosts []string

func main() {

	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	byte, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(byte, &App)
	fmt.Println(App)

	sendEmail("Атаншен-сервер запущен\n Хорошего дня, удачи в поисках мест на визу, у тебя все получится, неси добро в своем сердце, люби друзей и маму :3 (папу забудь)")

	setTimeOut(func() {
		htmlDoc, err := goquery.NewDocument("https://forum.awd.ru/viewtopic.php?f=326&t=326384&start=999999999999") // start=99999999999999 написано для того, чтобы скрипт всегда попадал на последнюю страницу форума
		if err != nil {
			fmt.Printf("Err %v\n", err)
			sendEmail("Проблемы с парсингом, нужно проверить сервер!")
		}

		htmlDoc.Find("div").Each(onDivFound) //начинаем парсинг
	})

	//"https://forum.awd.ru/viewtopic.php?f=326&t=326384&start=15900" <<- на этой странице "есть места", можно использовать для проверки парсера

}

func onDivFound(_ int, selection *goquery.Selection) {
	class, ok := selection.Attr("class")
	if ok && (class == "post bg1" || class == "post bg2") {
		onPostBgFound(selection)
	}
}

func onPostBgFound(selection *goquery.Selection) {
	children := selection.Children()
	children.Each(func(_ int, selection *goquery.Selection) {
		class, ok := selection.Attr("class")
		if ok && class == "inner" {
			onInnerFound(selection)
		}
	})
}

func onInnerFound(selection *goquery.Selection) {
	children := selection.Children()
	children.Each(func(_ int, selection *goquery.Selection) {
		class, ok := selection.Attr("class")
		if ok && class == "postbody" {
			onPostbodyFound(selection)
		}
	})
}

func onPostbodyFound(selection *goquery.Selection) {
	children := selection.Children()
	children.Each(func(_ int, selection *goquery.Selection) {
		class, ok := selection.Attr("class")
		if ok && class == "content" {
			onContentFound(selection)
		}
	})
}

func onContentFound(selection *goquery.Selection) {
	html, _ := selection.Html()
	htmlLowerCase := strings.ToLower(html)
	contains := strings.Contains(htmlLowerCase, "мест нет")
	if !contains {
		contains = strings.Contains(htmlLowerCase, "нет мест")
		if !contains {
			//значит есть места или кто-то решил написать на форуме что-то особенное (что маловероятно)
			isItSent := false
			vitalMessage := strings.TrimSpace(html)
			for i := range alreadyParsedPosts {
				if vitalMessage == alreadyParsedPosts[i] {
					isItSent = true
				}
			}
			if !isItSent {
				alreadyParsedPosts = append(alreadyParsedPosts, vitalMessage)
				fmt.Println("АТАНШЕН")
				msg := fmt.Sprintf("Найдены места! \n\n\n %s", vitalMessage)
				sendEmail(msg)
			}
		}
	}
}

func sendEmail(msg string) {
	auth := smtp.PlainAuth("", App.USR, App.PASS, "smtp.gmail.com")
	for _, mail := range App.MailArr {
		fmt.Println(mail)
		err := smtp.SendMail("smtp.gmail.com:587", auth, App.USR, []string{mail}, []byte(msg))
		if err != nil {
			fmt.Printf("Err when send email : %v\n", err)
		}
	}
}

func setTimeOut(handler func()) {
	for {
		handler()
		time.Sleep(3 * time.Minute)
	}
}
