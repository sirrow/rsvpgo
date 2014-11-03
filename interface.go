package rsvpgo

import (
	"encoding/xml"
	"fmt"
	"github.com/laurent22/ical-go"
	"io/ioutil"
	"net/http"
	"strings"
)

type Rsvp struct {
	Url         string
	ServiceName string
	EventName   string
	EventDate   string
	EventVenue  string
}

func RsvpGet(url string) (rsvp *Rsvp) {
	return nil
}

func rsvpget_twipla_checkurl(url string) bool {
	if strings.Contains(url, "http://twipla.jp/events/") {
		return true
	} else {
		return false
	}
}

func rsvpget_twipla_parse(data string) {

}

func rsvpget_twipla(url string) (rsvp *Rsvp) {
	if rsvpget_twipla_checkurl(url) == false {
		return nil
	}
	if resp, err := http.Get(url); err == nil {
		resp.Body.Close()
		return nil
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		ical.ParseCalendarNode

	}
}
