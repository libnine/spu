package tr

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func Trending(ch *chan []string) {
	return
}

func Greed(ch *chan []string) {
	data := cnn()
	*ch <- data
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
