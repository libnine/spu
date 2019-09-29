package rh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type popular struct {
	Symbol     string `json:"symbol",omitempty`
	Popularity int    `json:"popularity",omitempty`
}

type incdec struct {
	Start  int    `json:"start_popularity",omitempty`
	End    int    `json:"end_popularity",omitempty`
	Diff   int    `json:"popularity_difference",omitempty`
	Symbol string `json:"string",omitempty`
}

func RHpop() {
	var pops []popular
	var sym, sym25 string

	r, err := http.Get("https://robintrack.net/api/most_popular?limit=50&start_index=0")
	if err != nil {
		log.Fatal("Robinhood most popular unavailable.")
	}

	body, _ := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	err = json.Unmarshal(body, &pops)

	fmt.Printf("\n")
	for i := 0; i < 25; i++ {

		if len(pops[i].Symbol) == 5 {
			sym = fmt.Sprintf("%s", pops[i].Symbol)
		} else {
			sym = fmt.Sprintf("%s\t", pops[i].Symbol)
		}

		if len(pops[i+25].Symbol) == 5 {
			sym25 = fmt.Sprintf("%s", pops[i+25].Symbol)
		} else {
			sym25 = fmt.Sprintf("%s\t", pops[i+25].Symbol)
		}

		fmt.Printf("%v %s\t%v\t\t\t%v %s\t\t%v\n", i+1, sym,
			pops[i].Popularity, i+26, sym25, pops[i+25].Popularity)
	}
}

// func RHincdec(id string, tf string) {
// 	switch id {
// 	case "inc":
// 		url := fmt.Sprintf("https://robintrack.net/api/largest_popularity_changes?hours_ago=24&limit=50&percentage=false&min_popularity=50&start_index=0", %s)

// 	}
// }
