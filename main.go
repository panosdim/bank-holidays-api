package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getOrthodoxEaster(year int) time.Time {
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := ((d + e + 114) % 31) + 1
	// Orthodox Easter is calculated based on the Julian calendar, then converted to Gregorian
	easterJulian := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	easterGregorian := easterJulian.AddDate(0, 0, 13) // Add 13 days for Julian to Gregorian conversion
	return easterGregorian
}

func getHolidays(year int) []time.Time {
	holidays := make([]time.Time, 0)
	holidays = append(holidays, time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)) // New Year's Day
	holidays = append(holidays, time.Date(year, time.January, 6, 0, 0, 0, 0, time.UTC)) // Christmas Day (January 6th)

	easter := getOrthodoxEaster(year)
	holidays = append(holidays, easter.AddDate(0, 0, -48)) // Clean Monday (48 days before Easter)

	holidays = append(holidays, time.Date(year, time.March, 25, 0, 0, 0, 0, time.UTC)) // Annunciation (March 25th)
	holidays = append(holidays, easter.AddDate(0, 0, -2))                     // Good Friday
	holidays = append(holidays, easter.AddDate(0, 0, 1))                      // Easter Monday

	// Check if 1st May is moved
	greatWeekStart := easter.AddDate(0, 0, -7)
	greatWeekEnd := easter.AddDate(0, 0, 1)
	firstMay := time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC)
	if firstMay.Weekday() == time.Saturday || firstMay.Weekday() == time.Sunday {
		holidays = append(holidays, firstMay.AddDate(0, 0, int(time.Monday-firstMay.Weekday())))
	} else if firstMay.After(greatWeekStart) && firstMay.Before(greatWeekEnd) {
		holidays = append(holidays, easter.AddDate(0, 0, 2))
	} else {
		holidays = append(holidays, firstMay)
	}

	holidays = append(holidays, easter.AddDate(0, 0, 50))                      // Pentecost
	holidays = append(holidays, time.Date(year, time.August, 15, 0, 0, 0, 0, time.UTC))  // Dominion of the Theotokos (August 15th)
	holidays = append(holidays, time.Date(year, time.October, 28, 0, 0, 0, 0, time.UTC))  // Saint Demetrius' Day (October 26th)
	holidays = append(holidays, time.Date(year, time.December, 25, 0, 0, 0, 0, time.UTC)) // Christmas Day (December 25th)
	holidays = append(holidays, time.Date(year, time.December, 26, 0, 0, 0, 0, time.UTC)) // Boxing Day (December 26th)

	return holidays
}

func handleGetHolidays(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	if yearParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter 'year'")
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid year format")
		return
	}

	holidays := getHolidays(year)
	holidayJSON, err := json.Marshal(holidays)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling holidays data to JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(holidayJSON)
}

func main() {
	http.HandleFunc("/holidays", handleGetHolidays)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
