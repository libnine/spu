package quandl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Quandl struct {
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

func Qaaii() {
	q := Quandl{}

	url := fmt.Sprintf("https://www.quandl.com/api/v3/datasets/AAII/AAII_SENTIMENT.json?start_date=%s&end_date=%s&api_key=L-GpxP_AZvDf_67jqgMh",
		time.Now().Local().Add(time.Hour*-48).Format("2006-01-02"), time.Now().Local().Add(time.Hour*-24).Format("2006-01-02"))

	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &q)

	fmt.Printf("\n%s\t%s\n", q.Dataset.Data[0][0], q.Dataset.Name)
	fmt.Printf("\nBullish\tNeutral\tBearish\tBull-Bear Spread\t")
	fmt.Printf("\n%.2f\t%.2f\t%.2f\t%.2f\n\n", q.Dataset.Data[0][1].(float64)*100, q.Dataset.Data[0][2].(float64)*100,
		q.Dataset.Data[0][3].(float64)*100, q.Dataset.Data[0][6].(float64)*100)
}

func Qyc() {
	q := Quandl{}

	url := fmt.Sprintf("https://www.quandl.com/api/v3/datasets/USTREASURY/YIELD.json?start_date=%s&end_date=%s&api_key=L-GpxP_AZvDf_67jqgMh",
		time.Now().Local().Add(time.Hour*-48).Format("2006-01-02"), time.Now().Local().Add(time.Hour*-24).Format("2006-01-02"))

	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &q)

	fmt.Printf("\n%s\tYield Curve\n", q.Dataset.Data[0][0])
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", q.Dataset.Cols[1], q.Dataset.Cols[2], q.Dataset.Cols[3], q.Dataset.Cols[4],
		q.Dataset.Cols[5], q.Dataset.Cols[6], q.Dataset.Cols[7], q.Dataset.Cols[8],
		q.Dataset.Cols[9], q.Dataset.Cols[10], q.Dataset.Cols[11], q.Dataset.Cols[12])
	fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n\n", q.Dataset.Data[0][1], q.Dataset.Data[0][2], q.Dataset.Data[0][3], q.Dataset.Data[0][4],
		q.Dataset.Data[0][5], q.Dataset.Data[0][6], q.Dataset.Data[0][7], q.Dataset.Data[0][8],
		q.Dataset.Data[0][9], q.Dataset.Data[0][10], q.Dataset.Data[0][11], q.Dataset.Data[0][12])
}
