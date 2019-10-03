package tr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type symbol struct {
	Symbol string `json:"symbol"`
	Rank   int    `json:"sortOrder"`
}

type stsymbols struct {
	Aliases        []string `json:"aliases,omitempty"`
	ID             int      `json:"id,omitempty"`
	IsFollowing    bool     `json:"is_following,omitempty"`
	Symbol         string   `json:"symbol"`
	Title          string   `json:"title,omitempty"`
	WatchlistCount int      `json:"watchlist_count"`
}

func Trending() {
	// var yahooBool, stBool = make(chan bool), make(chan bool)
	var yahooBool = make(chan bool)
	data := []string{}

	go func() {
		y := yahoo()

		for i := range y {
			data = append(data, y[i].Symbol)
		}

		yahooBool <- true
	}()

	<-yahooBool

	fmt.Printf("\n")
	yahoo_print(data)

	// <- stBool
}

func Greed(ch *chan []string) {
	data := cnn()
	*ch <- data
}

func cnn() []string {
	data := []string{}
	r, err := http.Get("https://money.cnn.com/data/fear-and-greed/")
	if err != nil {
		data = append(data, "N/a")
		return data
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	greedNow, greedPrevClose := regexp.MustCompile(`Greed\sNow\:\s(.+?)\<\/li`), regexp.MustCompile(`Greed\sPrevious\sClose\:\s(.+?)\<\/li`)
	greedPrevWk, greedPrevMo, greedPrevYr := regexp.MustCompile(`Greed\s1\sWeek\sAgo\:\s(.+?)\<\/li`), regexp.MustCompile(`Greed\s1\sMonth\sAgo\:\s(.+?)\<\/li`), regexp.MustCompile(`Greed\s1\sYear\sAgo\:\s(.+?)\<\/li`)

	data = append(data, greedNow.FindStringSubmatch(string(body))[1], greedPrevClose.FindStringSubmatch(string(body))[1],
		greedPrevWk.FindStringSubmatch(string(body))[1], greedPrevMo.FindStringSubmatch(string(body))[1],
		greedPrevYr.FindStringSubmatch(string(body))[1])

	return data
}

func St() {
	var msgs []map[string]interface{}
	var ticker []stsymbols
	var tickers []string

	data := []string{}
	stmap := make(map[string]interface{})

	r, err := http.Get("https://api.stocktwits.com/api/2/streams/trending.json")
	if err != nil {
		data = append(data, "N/a")
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	_ = json.Unmarshal(body, &stmap)
	m, _ := json.Marshal(stmap["messages"])
	_ = json.Unmarshal(m, &msgs)

	for _, v := range msgs {
		j, _ := json.Marshal(v["symbols"])
		_ = json.Unmarshal(j, &ticker)
		fmt.Println(ticker[0].Symbol)
		tickers = append(tickers, ticker[0].Symbol)
	}
}

func yahoo() []symbol {
	tickers := []symbol{}

	r, err := http.Get("https://finance.yahoo.com/trending-tickers/")
	if err != nil {
		tickers = append(tickers, symbol{"N/a", 0})
		return tickers
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	yahoo := regexp.MustCompile(`{"trending_tickers":{"positions":(.+?),"name"`)
	_ = json.Unmarshal([]byte(yahoo.FindStringSubmatch(string(body))[1]), &tickers)

	return tickers
}

func yahoo_print(_yahoo []string) {
	var sym, sym2, sym3 string

	for i := 0; i < len(_yahoo)/3; i++ {
		if len(_yahoo[i]) > 4 && !(strings.Contains(_yahoo[i], "^")) {
			sym = fmt.Sprintf("%s", _yahoo[i])
		} else {
			sym = fmt.Sprintf("%s\t", _yahoo[i])
		}

		if len(_yahoo[i+10]) > 4 {
			sym2 = fmt.Sprintf("%s", _yahoo[i+10])
		} else {
			sym2 = fmt.Sprintf("%s\t", _yahoo[i+10])
		}

		if len(_yahoo[i+20]) > 4 && !(strings.Contains(_yahoo[i+20], "^")) {
			sym3 = fmt.Sprintf("%s", _yahoo[i+20])
		} else {
			sym3 = fmt.Sprintf("%s\t", _yahoo[i+20])
		}

		fmt.Printf("%v %s\t%v %s\t%v %s\n", i+1, sym, i+1+10, sym2, i+1+20, sym3)
	}
}
