package etc

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func News(newsrc string) {
	r, err := http.Get("https://finviz.com/news.ashx?v=2")
	if err != nil {
		log.Fatal("Couldn't get Bloomberg headlines.")
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	if newsrc == "bbg" {
		bbg := regexp.MustCompile(`bloomberg\.com\/\/news\/articles\/.+return\sfalse;\">\s+.+\s+.+\s+<a\shref\=\"(.+?)\"\starget.+link\"\>(.+?)<\/a>`)
		bbgParse := bbg.FindAllStringSubmatch(string(body), -1)

		fmt.Printf("\nBloomberg\n\n")
		for n := range bbgParse {
			fmt.Printf("%s\n%s\n\n", bbgParse[n][2], bbgParse[n][1])
		}

		return
	} else if newsrc == "rtrs" {
		rtrs := regexp.MustCompile(`(http:\/\/feeds\.reuters\.com\/~r\/reuters\/businessNews\/.+)\"\starget\=\"\_blank\".+tab-link\"\>(.+?)<\/a>`)
		rtrsParse := rtrs.FindAllStringSubmatch(string(body), -1)

		fmt.Printf("\nReuters\n\n")
		for n := range rtrsParse {
			fmt.Printf("%s\n%s\n\n", rtrsParse[n][2], rtrsParse[n][1])
		}

		return
	} else if newsrc == "wsj" {
		wsj := regexp.MustCompile(`(https\:\/\/www.wsj.com\/articles\/.+)\"\starget\=\"\_blank\".+link\">(.+?)\<\/a>`)
		wsjParse := wsj.FindAllStringSubmatch(string(body), -1)

		fmt.Printf("\nWSJ\n\n")
		for n := range wsjParse {
			fmt.Printf("%s\n%s\n\n", wsjParse[n][2], wsjParse[n][1])
		}
	}
}
