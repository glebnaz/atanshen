package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

/**
	Иерархия html документа, по которой производится поиск

	div class=post bg1 or class=post bg2
	div class=inner
	div class=postbody
	div class=content <<- здесь лежит нужное сообщение и иногда лишний html код (например, картинка)
 */

func main() {
	//"https://forum.awd.ru/viewtopic.php?f=326&t=326384&start=15900" <<- на этой странице "есть места", можно использовать для проверки парсера
	htmlDoc, err := goquery.NewDocument("https://forum.awd.ru/viewtopic.php?f=326&t=326384&start=99999999999999") // start=99999999999999 написано для того, чтобы скрипт всегда попадал на последнюю страницу форума
	if err != nil {
		fmt.Printf("Err %v\n", err)
		//TODO send to user
	}

	htmlDoc.Find("div").Each(onDivFound) //начинаем парсинг
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

			vitalMessage := strings.TrimSpace(html)
			fmt.Println("АТАНШЕН")
			fmt.Println(vitalMessage)
			fmt.Println()
			//TODO send message to user
		}
	}
}
