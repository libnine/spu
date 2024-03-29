package fin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func LUexp() {
	r, err := http.Get("https://www.marketbeat.com/ipos/lockup-expirations")
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	html_re := regexp.MustCompile(`ratingstable\'[\s\S]+Priced\<\/th\>\<\/tr\>\<\/thead\>\<tbody\>(.+?)\<\/tbody\>`)
	html := html_re.FindStringSubmatch(string(body))

	tr_re := regexp.MustCompile(`\<tr>(.+?)\<\/tr\>`)
	tr := tr_re.FindAllStringSubmatch(html[1], -1)

	var name_fmt string
	i := 0
	fmt.Println("Ticker\tName\t\t\t\tLast\tExpiration\tShares\t\tIPO\tOffer Size\tDate Priced")
	for n := range tr {
		td_re := regexp.MustCompile(`\<td\>(.+?)\<\/td\>`)
		td := td_re.FindAllStringSubmatch(tr[n][1], -1)

		re_name := regexp.MustCompile(`(.+?)\s\(\<a\shref[\s\S]+\)`)
		name := re_name.FindStringSubmatch(td[0][1])

		re_sym := regexp.MustCompile(`\(\<a\shref\=\'\/[\s\S]+\>(.+?)\<`)
		sym := re_sym.FindStringSubmatch(td[0][1])

		t, _ := time.Parse("1/02/2006", strings.TrimSpace(td[2][1]))

		if t.Sub(time.Now().Local()).Seconds() < 0 {
			continue
		}

		switch x := strings.TrimSpace(name[1]); {
		case len(x) <= 7:
			name_fmt = fmt.Sprintf("%s\t\t\t\t", x)
		case len(x) > 7 && len(x) <= 15:
			name_fmt = fmt.Sprintf("%s\t\t\t", x)
		case len(x) > 15 && len(x) <= 22:
			name_fmt = fmt.Sprintf("%s\t\t", x)
		case len(x) > 22 && len(x) <= 30:
			name_fmt = fmt.Sprintf("%s\t", x)
		case len(x) > 30:
			name_fmt = fmt.Sprintf("%s", x[:30])
		}

		fmt.Printf("%s\t%s%s\t%s\t%s\t%s\t%s\t%s\n", sym[1], name_fmt, td[1][1], td[2][1], td[3][1], td[4][1], td[5][1], td[6][1])
		if i == 25 {
			break
		}
		i++
	}
}
