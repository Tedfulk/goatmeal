package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

func GetLocationInfo() (*IPInfo, error) {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func GetFormattedLocationAndTime() string {
	info, err := GetLocationInfo()
	if err != nil {
		return fmt.Sprintf("Unknown location, Current time: %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	loc, err := time.LoadLocation(info.Timezone)
	if err != nil {
		loc = time.Local
	}

	return fmt.Sprintf("Location: %s, %s, %s, Current time: %s", 
		info.City, info.Region, info.Country, 
		time.Now().In(loc).Format("2006-01-02 15:04:05"))
} 