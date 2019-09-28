package fin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func Finviz(ticker string, ch *chan []string) {
	data := []string{}

	url := fmt.Sprintf("https://finviz.com/quote.ashx?t=%s", ticker)
	r, err := http.Get(url)
	if err != nil {
		log.Fatal("Finviz request failed.")
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	chg, er := regexp.MustCompile(`Change\<\/td\>.+\"\>(.+?)\<`), regexp.MustCompile(`Earnings\<\/td\>\<td\swidth=\".+\"\>\<b\>(.+?)\<`)
	a, b := chg.FindStringSubmatch(string(body)), er.FindStringSubmatch(string(body))

	if len(a) == 0 {
		a = append(a, "", "N/a")
	}

	if len(b) == 0 {
		b = append(b, "", "N/a\t")
	}

	data = append(data, a[1], b[1])

	*ch <- data
}
