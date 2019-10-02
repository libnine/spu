package tr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type symbol struct {
	Symbol string `json:"symbol"`
	Rank   int    `json:"sortOrder"`
}

func Trending(ch *chan []string) {
	data := []string{}
	y := yahoo()

	for i := range y {
		data = append(data, y[i].Symbol)
	}

	*ch <- data
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

// func st() []string {
// 	data := []string{}
// 	r, err := http.Get("https://api.stocktwits.com/api/2/streams/trending.json")
// 	if err != nil {
// 		data = append(data, "N/a")
// 		return data
// 	}

// 	return data
// }

func yahoo() []symbol {
	tickers := []symbol{}

	r, err := http.Get("https://finance.yahoo.com/trending-tickers/")
	if err != nil {
		tickers = append(tickers, symbol{"N/a", 0})
		// return tickers
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	yahoo := regexp.MustCompile(`{"trending_tickers":{"positions":(.+?),"name"`)
	_ = json.Unmarshal([]byte(yahoo.FindStringSubmatch(string(body))[1]), &tickers)

	return tickers
}
