package cal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Calendar() {
	var temp []string
	var weekset [][]string

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
		detail := fmt.Sprintf("%s %s", t[i][0], e[i][1])
		temp = append(temp, detail)

		if strings.Split(t[i][0], ":")[0] > strings.Split(t[i+1][0], ":")[0] && strings.Split(t[i+1][0], " ")[1] == "AM" ||
			strings.Split(t[i][0], " ")[1] == "PM" && strings.Split(t[i+1][0], " ")[1] == "AM" {
			fmt.Println(temp)
			fmt.Println("YAY!")
			time.Sleep(3 * time.Second)
			weekset = append(weekset, temp)
			temp = temp[:0]
		}
	}

	for n := range weekset {
		fmt.Println(n)
	}

	// str := fmt.Sprintf("%s %s", t[i][0], e[i][1])
	// fmt.Println(str)
}
