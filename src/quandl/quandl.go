package quandl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Aaii struct {
	Dataset DS `json:"dataset"`
}

type DS struct {
	ID          string          `json:"id"`
	DScode      string          `json:"dataset_code"`
	DBcode      string          `json:"database_code"`
	Name        string          `json:"name"`
	Desc        string          `json:"description"`
	Refreshed   *time.Time      `json:"refreshed_at"`
	NewestAvail string          `json:"newest_available_date"`
	OldestAvail string          `json:"oldest_available_date"`
	Cols        []string        `json:"column_names"`
	Freq        string          `json:"frequency"`
	Type        string          `json:"type"`
	Premium     string          `json:"premium"`
	Limit       string          `json:"limit"`
	Transform   string          `json:"transform"`
	ColIndex    string          `json:"column_index"`
	Start       string          `json:"start_date"`
	End         string          `json:"end_date"`
	Data        [][]interface{} `json:"data"`
	Collapse    string          `json:"collapse"`
	Order       string          `json:"order"`
	DbID        string          `json:"database_id"`
}

func Q() {
	q := Aaii{}

	r, err := http.Get("https://www.quandl.com/api/v3/datasets/AAII/AAII_SENTIMENT.json?start_date=2019-11-21&end_date=2019-11-21&api_key=L-GpxP_AZvDf_67jqgMh")
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &q)

	fmt.Printf("\n%s\t%s\n", q.Dataset.Data[0][0], q.Dataset.Name)
	fmt.Printf("\nBullish\tNeutral\tBearish\tBull-Bear Spread\t")
	fmt.Printf("\n%.2f\t%.2f\t%.2f\t%.2f\n", q.Dataset.Data[0][1].(float64)*100, q.Dataset.Data[0][2].(float64)*100,
		q.Dataset.Data[0][3].(float64)*100, q.Dataset.Data[0][6].(float64)*100)
}
