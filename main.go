package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	doc, err := goquery.NewDocument("https://forum.awd.ru/viewtopic.php?f=326&t=326384&start=99999999999999")
	if err != nil {
		fmt.Printf("Err %v\n", err)
	}

	doc.Find("div").Each(func(index int, s *goquery.Selection) {
		// получаем ссылку из атрибута
		if class, ok := s.Attr("class"); ok && class == "post bg1" {
			ch := s.Children()
			ch.Each(func(index int, s *goquery.Selection) {
				if class, ok := s.Attr("class"); ok && class == "inner" {
					ch := s.Children()
					ch.Each(func(index int, s *goquery.Selection) {
						if class, ok := s.Attr("class"); ok && class == "postbody" {
							ch := s.Children()
							ch.Each(func(index int, s *goquery.Selection) {
								if class, ok := s.Attr("class"); ok && class == "content" {
									html, _ := s.Html()
									//fmt.Println(html)
									htmlLow:=strings.ToLower(html)
									have := strings.Contains(htmlLow,"мест нет")
									if !have {
										have = strings.Contains(htmlLow,"нет мест")
										if !have {
											fmt.Println("AHTUNG")
											fmt.Println(strings.TrimSpace(html))
										}
									}
								}
							})
						}
					})
				}
			})
		}
	})

}
