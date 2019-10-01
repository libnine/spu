package cal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
		temp = append(temp, fmt.Sprintf("%s %s", t[i][0], e[i][1]))

		if len(weekset) == 4 && i == len(e)-1 {
			weekset = append(weekset, temp)
			break
		}

		x, err := strconv.Atoi(strings.Split(t[i][0], ":")[0])
		if err != nil {
			log.Fatal("Error populating calendar data.")
		}

		y, err := strconv.Atoi(strings.Split(t[i+1][0], ":")[0])
		if err != nil {
			log.Fatal("Error populating calendar data.")
		}

		if x > y && strings.Split(t[i+1][0], " ")[1] == "AM" ||
			strings.Split(t[i][0], " ")[1] == "PM" && strings.Split(t[i+1][0], " ")[1] == "AM" {
			weekset = append(weekset, temp)
			temp = nil
		}
	}

	for a := range weekset {
		fmt.Println(weekset[a])
	}

	dn := int(time.Now().Weekday())

	if dn == 6 || dn == 7 {
		dn = 1
	}

	fmt.Printf("\n%s\n", time.Now().Weekday())
	for n := range weekset[dn] {
		fmt.Printf("\n%s", weekset[dn][n])
	}

	fmt.Printf("\n")
}
