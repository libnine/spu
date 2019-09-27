package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

func main() {
	args := os.Args[1:]
	var syms []string
	var data []string
	var wg sync.WaitGroup

	for _, v := range args {
		if len(v) > 5 {
			fmt.Printf("%s: invalid ticker format.\n", v)
			continue
		}

		syms = append(syms, v)
	}

	if len(syms) != 0 {
		for _, v := range syms {
			wg.Add(1)
			go func(v string, wg *sync.WaitGroup) {
				ch := make(chan []string)
				ch_ := make(chan []string)

				go vol(v, &ch)
				go finviz(v, &ch_)

				_vol := <-ch
				_finviz := <-ch_

				dump := fmt.Sprintf("%s\t%s\t\t%s\t%s\t\t%s\t\t%s\t\t%s", strings.ToUpper(v), _finviz[0], _finviz[1], _vol[1], _vol[0], _vol[2], _vol[3])

				data = append(data, dump)
				wg.Done()
			}(v, &wg)
		}
	}

	wg.Wait()
	pretty(data)
}

func vol(ticker string, ch *chan []string) {
	data := []string{}

	url := fmt.Sprintf("http://www.optionslam.com/earnings/stocks/%s", strings.ToUpper(ticker))
	r, err := http.Get(url)
	if err != nil {
		log.Fatal("Options request failed.")
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	iv, short := regexp.MustCompile(`Current 7 Day Implied Movement:(?s)\s+\</td\>(?s)\s+\<td\>(?s)\s+(.+?)\s+`), regexp.MustCompile(`Short Interest:\s(.+?)\s+`)
	ivweekly, exp := regexp.MustCompile(`Implied Move Weekly\:\s+\<\/td\>\s+\<td>\s+(.+?)\s+`), regexp.MustCompile(`Expires on\:\s(.+?)\s+\<\/font\>`)

	a, b, c, d := iv.FindStringSubmatch(string(body)), short.FindStringSubmatch(string(body)), ivweekly.FindStringSubmatch(string(body)), exp.FindStringSubmatch(string(body))

	if len(a) == 0 {
		a = append(a, "", "N/a")
	}

	if len(b) == 0 {
		b = append(b, "", "N/a")
	}

	if len(c) == 0 {
		c = append(c, "", "N/a")
	}

	if len(d) == 0 {
		d = append(d, "", "N/a")
	}

	if b[1] == "None" {
		b[1] = "N/a"
	} else if !(b[1] == "N/a") {
		b[1] = fmt.Sprintf("%s%%", b[1])
	}

	data = append(data, a[1], b[1], c[1], d[1])

	*ch <- data
}

func finviz(ticker string, ch *chan []string) {
	data := []string{}

	url := fmt.Sprintf("https://finviz.com/quote.ashx?t=%s", ticker)
	r, err := http.Get(url)
	if err != nil {
		log.Fatal("Finviz request failed.")
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	chg, er := regexp.MustCompile(`Change\<\/td\>.+\"\>(.+?)\<`), regexp.MustCompile(`Earnings\<\/td\>\<td\swidth=\".+\"\>\<b\>(.+?)\<`)
	data = append(data, chg.FindStringSubmatch(string(body))[1], er.FindStringSubmatch(string(body))[1])

	*ch <- data
}

func pretty(dump []string) {
	sort.Strings(dump)
	fmt.Println("\n\tChange %\tEarnings\tShort Interest\t7 Day Implied\tWeekly Implied\tWeekly Expiration")
	for _, v := range dump {
		fmt.Println(v)
	}
}
