package vol

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func Vol(ticker string, ch *chan []string) {
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
