package etc

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Data struct {
	Pub     Publication `xml:"publication"`
	PubDate *time.Time  `xml:"publication_date"`
	Title   string      `xml:"title"`
	Kws     string      `xml:"keywords"`
}

type Image struct {
	Loc string `xml:"loc"`
	Cap string `xml:"caption"`
}

type Publication struct {
	Name string `xml:"name"`
	Lang string `xml:"language"`
}

type Urls struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []Url    `xml:"url"`
}

type Url struct {
	Loc     string `xml:"loc"`
	Link    string `xml:"link"`
	Details Data   `xml:"news"`
	Img     Image  `xml:"image"`
}

func News(newsrc string) {
	if newsrc == "mw" {
		mw := Urls{}

		r, err := http.Get("https://www.marketwatch.com/mw_news_sitemap.xml")
		if err != nil {
			log.Fatal("Couldn't get MarketWatch data.")
		}

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		err = xml.Unmarshal(body, &mw)

		for n := range mw.URLs {
			fmt.Printf("%v %s\n", mw.URLs[n].Details.PubDate.Local().Format(time.RFC1123), mw.URLs[n].Details.Title)
			if mw.URLs[n].Details.PubDate.Local().Day() != time.Now().Day() || n == 10 {
				return
			}
		}
	}

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

	} else if newsrc == "rtrs" {
		rtrs := regexp.MustCompile(`(http:\/\/feeds\.reuters\.com\/~r\/reuters\/businessNews\/.+)\"\starget\=\"\_blank\".+tab-link\"\>(.+?)<\/a>`)
		rtrsParse := rtrs.FindAllStringSubmatch(string(body), -1)

		fmt.Printf("\nReuters\n\n")
		for n := range rtrsParse {
			fmt.Printf("%s\n%s\n\n", rtrsParse[n][2], rtrsParse[n][1])
		}

	} else if newsrc == "wsj" {
		wsj := regexp.MustCompile(`(https\:\/\/www.wsj.com\/articles\/.+)\"\starget\=\"\_blank\".+link\">(.+?)\<\/a>`)
		wsjParse := wsj.FindAllStringSubmatch(string(body), -1)

		fmt.Printf("\nWSJ\n\n")
		for n := range wsjParse {
			fmt.Printf("%s\n%s\n\n", wsjParse[n][2], wsjParse[n][1])
		}
	}

	return
}
