package cal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func Calendar() {
	r, err := http.Get("https://us.econoday.com/byweek.asp?cust=us")
	if err != nil {
		log.Fatal("Couldn't get calendar data.")
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	events, times := regexp.MustCompile(`econoevents[\s]{0,1}[a-zA-Z]{0,5}[a-zA-Z+]{0,1}\"\>\<[a-zA-Z\s\=\"\.\?\d\&\/]+\#top\"\>(.+?)\<`), regexp.MustCompile(`\d+\:\d+\s[AP][M]`)
	e, t := events.FindAllStringSubmatch(string(body), -1), times.FindAllStringSubmatch(string(body), -1)

	for i := range e {
		str := fmt.Sprintf("%s %s", t[i][0], e[i][1])
		fmt.Println(str)
	}

}