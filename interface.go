package rsvpgo

import (
	"bufio"
	"fmt"
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

func rsvpget_twipla_parse(data string) (rsvp *Rsvp) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
	return nil
}

func rsvpget_twipla(url string) (rsvp *Rsvp) {
	if rsvpget_twipla_checkurl(url) == false {
		return nil
	}

	if resp, err := http.Get(url + "/.ics"); err == nil {
		resp.Body.Close()
		return nil
	} else {
		defer resp.Body.Close()
		br := bufio.NewReader(resp.Body)

		for line, err := br.ReadString('\n'); err != nil; line, err = br.ReadString('\n') {
			fmt.Printf("%s", line)
		}
	}
	return nil
}
