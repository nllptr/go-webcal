package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Month is a collection of data used in the month view template.
type Month struct {
	Dates [][]time.Time
	Weeks []int
}

func monthHandler(w http.ResponseWriter, r *http.Request) {
	dateString := r.URL.Path[len("/month/"):]
	parsed, err := time.Parse("2006/01/02", dateString)
	if err != nil {
		log.Println("could not parse date (using current date instead)")
		parsed = time.Now()
	}
	log.Println(parsed)
	startDate := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, parsed.Location())
	startDate = startDate.AddDate(0, 0, -int((startDate.Weekday()+6)%7))
	endDate := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, parsed.Location())
	endDate = endDate.AddDate(0, 1, -1)
	endDate = endDate.AddDate(0, 0, int(6-(endDate.Weekday()+6)%7))
	endDate = endDate.AddDate(0, 0, 1)
	month := [][]time.Time{}
	weekNumbers := []int{}
	week := []time.Time{}
	date := startDate
	for date.Before(endDate) {
		week = append(week, date)
		date = date.AddDate(0, 0, 1)
		if len(week) == 7 {
			month = append(month, week)
			_, wNo := date.ISOWeek()
			weekNumbers = append(weekNumbers, wNo)
			week = []time.Time{}
		}
	}
	t, err := template.ParseFiles("monthView.html")
	if err != nil {
		log.Println("could not parse template: ", err)
	}
	t.Execute(w, Month{month, weekNumbers})
}

func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/month/", monthHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on port 8080...")
}
