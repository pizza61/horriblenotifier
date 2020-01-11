package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/PuerkitoBio/goquery"
)

type Settings struct {
	Config   Config     `json:"config"`
	Schedule []Horrible `json:"schedule"`
}

type Horrible struct {
	Title string `json:"title"`
	Time  string `json:"time"`
}

// Server  init HorribleNotifier configuration panel
func (n *Notificator) Server() {
	fs := http.FileServer(rice.MustFindBox("static").HTTPBox())
	http.Handle("/", fs)

	http.HandleFunc("/nyaa", func(w http.ResponseWriter, r *http.Request) {
		guids, ok := r.URL.Query()["guid"]

		if !ok || len(guids[0]) < 1 {
			log.Println("GUID not found")
			http.NotFound(w, r)
		}

		guid := guids[0]
		for _, v := range n.notifications {
			if v.guid == guid {
				http.Redirect(w, r, v.magnet, 301)
			}
		}
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET": // Get settings
			// HorribleSubs Schedule
			res, err := http.Get("https://horriblesubs.info/release-schedule")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			defer res.Body.Close()

			if res.StatusCode != 200 {
				w.WriteHeader(http.StatusInternalServerError)
			}

			document, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			/*** Schedule ***/
			schedule := []Horrible{} /***/
			/*** Schedule ***/

			ts := strings.NewReplacer("’", "'", "–", "-")
			document.Find(".schedule-page-item").Each(func(i int, el *goquery.Selection) {
				title := ts.Replace(el.Find("td:first-child").Text())
				time := el.Find(".schedule-time").Text()

				schedule = append(schedule, Horrible{title, time})
			})

			config := n.GetConfig()
			response := Settings{config, schedule}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		case "POST": // Set settings
			var ns Config

			err := json.NewDecoder(r.Body).Decode(&ns)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = n.SetConfig(ns)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "OK")
		}

	})

	http.ListenAndServe(":3939", nil)
}
