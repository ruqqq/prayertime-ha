package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/osext"
)

func parseLocationCode() string {
	code := os.Getenv("LOCATION_CODE")
	if code != "" {
		return code
	}

	return "SG"
}

var prayerTimeLocationCode = parseLocationCode()

var ApiUrl = "https://ruqqq.github.io/prayertimes-database/data/" + prayerTimeLocationCode + "/1/"

type Result struct {
	Date  int      `json:"date"`
	Month int      `json:"month"`
	Year  int      `json:"year"`
	Times []string `json:"times"`
}

func (r *Result) FullDate() string {
	return fmt.Sprintf("%d-%d-%d", r.Date, r.Month, r.Year)
}

func getPrayerTimes(datetime time.Time, noCache ...bool) (*Result, error) {
	fmt.Println("=== TRYING CACHE")
	result, err := getPrayerTimesFromCache(datetime)
	if result == nil || (len(noCache) > 0 && noCache[0]) {
		fmt.Println("=== ...CACHE UNAVAILABLE")
		fmt.Println("=== FETCHING PRAYERTIMES")

		result, err = getPrayerTimesFromServer(datetime)
		if err != nil {
			fmt.Printf("=== ERR: %v\n", err)
			return nil, err
		}
	} else {
		fmt.Println("=== ...USING CACHE")
	}

	return result, nil
}

func getPrayerTimesFromServer(datetime time.Time) (*Result, error) {
	year := datetime.Format("2006")
	date := datetime.Format("2-1-2006")

	fmt.Printf("...FROM: %s\n", ApiUrl+year+".json")

	response, err := http.Get(ApiUrl + year + ".json")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(contents))

	path, _ := osext.ExecutableFolder()
	err = ioutil.WriteFile(path+"/"+year+".json", contents, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}

	var results [][]Result

	err = json.Unmarshal(contents, &results)
	if err != nil {
		return nil, err
	}

	for _, month := range results {
		for _, result := range month {
			if result.FullDate() == date {
				return &result, nil
			}
		}
	}

	return nil, nil
}

func getPrayerTimesFromCache(datetime time.Time) (*Result, error) {
	year := datetime.Format("2006")
	date := datetime.Format("2-1-2006")

	path, _ := osext.ExecutableFolder()
	contents, err := ioutil.ReadFile(path + "/" + year + ".json")
	if err != nil {
		return nil, err
	}

	var results [][]Result

	err = json.Unmarshal(contents, &results)
	if err != nil {
		return nil, err
	}

	for _, month := range results {
		for _, result := range month {
			if result.FullDate() == date {
				return &result, nil
			}
		}
	}

	return nil, nil
}
