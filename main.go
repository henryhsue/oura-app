package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	durationDays, err := strconv.Atoi(os.Getenv("DURATION_DAYS"))
	if err != nil {
		durationDays = 1
	}
	start, end := genDatesByDaysOfDuration(durationDays)

	// process user info
	body, err := makeRequest("/v1/userinfo")
	check(err)
	var userInfo *UserInfo
	err = json.Unmarshal(body, &userInfo)
	check(err)
	log.Printf("User Info: %#v\n", userInfo)

	// process sleep data
	body, err = makeRequest(fmt.Sprintf("/v1/sleep?start=%s&end=%s", start, end))
	check(err)
	var sleepSummaries SleepSummaries
	err = json.Unmarshal(body, &sleepSummaries)
	check(err)
	log.Printf("Sleep Summary: %#v\n", sleepSummaries)
}

func genDatesByDaysOfDuration(days int) (startStr string, endStr string) {
	duration := time.Hour * 24 * time.Duration(int64(days))
	end := time.Now()
	endStr = end.Format("2006-01-02")
	start := end.Add(-duration)
	startStr = start.Format("2006-01-02")
	return startStr, endStr
}

type UserInfo struct {
	Age    int32
	Weight float32
	Height float32
	Gender string
	Email  string
}

type SleepSummaries struct {
	Sleep []SleepSummary // by date
}

type SleepSummary struct {
	SummaryDate       string    `json:"summary_date"`
	PeriodID          int       `json:"period_id"`
	IsLongest         int       `json:"is_longest"`
	Timezone          int       `json:"timezone"`
	BedtimeStart      time.Time `json:"bedtime_start"`
	BedtimeEnd        time.Time `json:"bedtime_end"`
	Score             int       `json:"score"`
	ScoreTotal        int       `json:"score_total"`
	ScoreDisturbances int       `json:"score_disturbances"`
	ScoreEfficiency   int       `json:"score_efficiency"`
	ScoreLatency      int       `json:"score_latency"`
	ScoreRem          int       `json:"score_rem"`
	ScoreDeep         int       `json:"score_deep"`
	ScoreAlignment    int       `json:"score_alignment"`
	Total             int       `json:"total"`
	Duration          int       `json:"duration"`
	Awake             int       `json:"awake"`
	Light             int       `json:"light"`
	Rem               int       `json:"rem"`
	Deep              int       `json:"deep"`
	OnsetLatency      int       `json:"onset_latency"`
	Restless          int       `json:"restless"`
	Efficiency        int       `json:"efficiency"`
	MidpointTime      int       `json:"midpoint_time"`
	HrLowest          int       `json:"hr_lowest"`
	HrAverage         float64   `json:"hr_average"`
	Rmssd             int       `json:"rmssd"`
	BreathAverage     float64   `json:"breath_average"`
	TemperatureDelta  float64   `json:"temperature_delta"`
	Hypnogram5Min     string    `json:"hypnogram_5min"`
	Hr5Min            []int     `json:"hr_5min"`
	Rmssd5Min         []int     `json:"rmssd_5min"`
}

func makeRequest(path string) ([]byte, error) {
	// create request
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.URL.Scheme = "https"
	req.URL.Host = "api.ouraring.com"
	token, err := ioutil.ReadFile("bearer.token")
	check(err)
	req.Header.Add("Authorization", string(token))

	// process request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

// Trigger is the payload of a Pub/Sub event.
type Trigger struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func TriggerRun(ctx context.Context, m Trigger) error {
	main()
	return nil
}
