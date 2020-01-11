package main

import (
	"log"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/toast.v1"
)

// Notificator ...
func (n *Notificator) Notificator() {
	// Get refresh rate from configuration file
	config := n.GetConfig()
	// Check HorribleSubs rss feed every X minutes
	ticker := time.NewTicker(time.Duration(config.Refresh) * time.Minute)
	noti()
	fp := gofeed.NewParser()
	go func() {
		for {
			// get current config because it could have changed
			config = n.GetConfig()
			//file, _ := os.Open("tests/test.xml")
			//feed, err := fp.Parse(file)
			feed, err := fp.ParseURL("http://www.horriblesubs.info/rss.php?res=" + config.Quality)
			if err != nil {
				log.Println("Failed to fetch RSS feed")
				break
			}
			log.Println("Fetch")
			for i, v := range feed.Items {
				exp := regexp.MustCompile(`^\[HorribleSubs\] (?P<title>.*?) - (?P<episode>\d{1,4}) \[(360p|480p|720p|1080p)\](?:\.mkv)$`)
				match := exp.FindStringSubmatch(v.Title)
				// 1 - title, 2 - episode
				if len(n.lastGUID) > 0 {
					if v.GUID == n.lastGUID {
						if i > 0 {
							n.lastGUID = feed.Items[i-1].GUID
						} else {
							n.lastGUID = v.GUID
						}
						break
					} else { // notify
						if config.SubscribedAll || find(config.Subscriptions, match[1]) {
							log.Println("Notification")
							n.Add(match[1], match[2], v.Link, v.GUID)
						}
					}
				} else { // first run
					n.lastGUID = v.GUID
					break
				}
			}
			<-ticker.C
		}
	}()
}

func noti() {
	ntf := toast.Notification{
		AppID: "HorribleNotifier",
		Title: "HorribleNotifier is now running",
		Actions: []toast.Action{
			{"protocol", "Open settings", "http://localhost:3939"},
		},
	}
	log.Println("dos")
	err := ntf.Push()
	if err != nil {
		log.Println("Failed to send notification")
	}
}

// Add ...
func (n *Notificator) Add(title string, ep string, magnet string, guid string) {
	release := Notification{guid: guid, magnet: magnet, title: title, episode: ep}
	n.notifications = append(n.notifications, release)

	release.Notify()
}

// Notify ...
func (n *Notification) Notify() {
	notification := toast.Notification{
		AppID:   "HorribleNotifier",
		Title:   n.title,
		Message: "Episode " + n.episode + " released",
	}
	notification.ActivationArguments = "http://localhost:3939/nyaa?guid=" + n.guid
	err := notification.Push()
	if err != nil {
		log.Println("Failed to send notification")
	}
}

func find(arr []string, search string) bool {
	for _, el := range arr {
		if el == search {
			return true
		}
	}
	return false
}
