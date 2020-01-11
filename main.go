package main

import (
	"errors"
	"log"
	"net"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
)

type Notificator struct {
	notifications []Notification
	lastGUID      string
}

type Notification struct {
	magnet  string
	guid    string
	title   string
	episode string
}

func main() {
	systray.Run(ready, exit)
}

func ready() {
	err := porttest()
	if err != nil {
		dialog.Message("%s", "HorribleNotifier is already running. Exiting.").Title("HorribleNotifier").Error()
		exit()
	}
	initIcon()
	n := new(Notificator)

	n.Notificator()
	n.Server()
}

func initIcon() {
	iconBox := rice.MustFindBox("icons")

	file, err := iconBox.Bytes("hn-128.ico")
	if err != nil {
		log.Fatalln("couldn't find icon")
	}

	systray.SetIcon(file)
	systray.SetTitle("HorribleNotifier")
	systray.SetTooltip("HorribleNotifier")

	itemConfig := systray.AddMenuItem("Settings", "Settings")
	systray.AddSeparator()
	itemExit := systray.AddMenuItem("Quit", "Quit HorribleNotifier")

	go func() {
		for {
			select {
			case <-itemExit.ClickedCh:
				exit()
			case <-itemConfig.ClickedCh:
				open.Run("http://localhost:3939")
			}
		}
	}()
}

func porttest() error {
	ln, err := net.Listen("tcp", ":3939")

	if err != nil {
		return errors.New("Already running")
	}

	ln.Close()
	return nil
}

func exit() {
	log.Println("Exiting")
	systray.Quit()
	os.Exit(0)
}
