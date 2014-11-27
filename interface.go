package rsvpgo

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Rsvp struct {
	ServiceName string
	EventId     string
	EventName   string
	EventDate   time.Time
	EventVenue  string
}

type twvt_date struct {
	Start_date string `xml:"start_date"`
	Start_time string `xml:"start_time"`
}

type twvt_location struct {
	Location_name    string `xml:"location_name"`
	Location_address string `xml:"location_address"`
}

type twvt_xml struct {
	Title    string        `xml:"title"`
	Date     twvt_date     `xml:"date"`
	Location twvt_location `xml:"location"`
}

func (r Rsvp) String() string {
	return r.ServiceName + " " + r.EventId + " " + r.EventName + " " + r.EventDate.String() + " " + r.EventVenue
}

func RsvpGet(url string) (rsvp *Rsvp) {
	if r := rsvpget_twipla(url); r != nil {
		return r
	}
	if r := rsvpget_tweetvite(url); r != nil {
		return r
	}
	return nil
}

func strcmp_from_beginning_of_line(a string, b string) bool {
	var shorter int
	if len(a) < len(b) {
		shorter = len(a)
	} else {
		shorter = len(b)
	}

	if a[0:shorter] == b[0:shorter] {
		return true
	}
	return false
}

func tweetvite_getid(url string) (id string) {
	idx := strings.LastIndex(url, "/")
	return url[idx+1:]
}

func rsvpget_tweetvite_checkurl(url string) bool {
	if strcmp_from_beginning_of_line(url, "http://tweetvite.com/event/") {
		return true
	} else if strcmp_from_beginning_of_line(url, "http://twvt.us/") {
		return true
	} else {
		return false
	}
}

func rsvpget_tweetvite(url string) (rsvp *Rsvp) {
	if rsvpget_tweetvite_checkurl(url) == false {
		return nil
	}
	id := tweetvite_getid(url)
	if resp, err := http.Get("http://tweetvite.com/api/1.0/rest/events/event?public_id=" + id + "&format=xml"); err != nil {
		fmt.Printf("%s\n", err.Error)
		return nil
	} else {
		defer resp.Body.Close()
		tv := &Rsvp{}
		tv.EventId = id
		tv.ServiceName = "tweetvite"
		t := twvt_xml{"", twvt_date{"", ""}, twvt_location{"", ""}}
		body, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(body, &t)
		tv.EventName = t.Title
		tv.EventVenue = t.Location.Location_name + " " + t.Location.Location_address

		year := 0
		var mon time.Month = 0
		day := 0
		hour := 0
		min := 0
		sec := 0

		if len(t.Date.Start_date) != 0 {
			year, _ = strconv.Atoi(t.Date.Start_date[0:4])
			month, _ := strconv.Atoi(t.Date.Start_date[5:7])
			mon = time.Month(month)
			day, _ = strconv.Atoi(t.Date.Start_date[8:10])
		}
		if len(t.Date.Start_time) != 0 {
			hour, _ = strconv.Atoi(t.Date.Start_time[0:2])
			min, _ = strconv.Atoi(t.Date.Start_time[3:5])
			sec, _ = strconv.Atoi(t.Date.Start_time[6:8])
		}
		eventdate_local := time.Date(year, mon, day, hour, min, sec, 0, time.Local)
		eventdate := eventdate_local.In(time.UTC)
		tv.EventDate = eventdate

		return tv
	}
	return nil
}

func rsvpget_twipla_checkurl(url string) bool {
	if strcmp_from_beginning_of_line(url, "http://twipla.jp/events/") {
		return true
	} else {
		return false
	}
}

func rsvpget_twipla(url string) (rsvp *Rsvp) {
	if rsvpget_twipla_checkurl(url) == false {
		return nil
	}
	if resp, err := http.Get(url + "/.ics"); err != nil {
		fmt.Printf("%s\n", err.Error)
		resp.Body.Close()
		return nil
	} else {
		defer resp.Body.Close()
		tp := &Rsvp{}
		idx := strings.LastIndex(url, "/")
		tp.EventId = url[idx+1:]
		tp.ServiceName = "twipla"
		br := bufio.NewReader(resp.Body)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			sline := string(line)
			if strcmp_from_beginning_of_line(sline, "LOCATION:") {
				idx := strings.Index(sline, ":")
				tp.EventVenue = sline[idx+1:]
			} else if strcmp_from_beginning_of_line(sline, "DTSTART:") {
				idx := strings.Index(sline, ":")
				datestr := sline[idx+1:]
				if len(datestr) != 0 {
					year, _ := strconv.Atoi(datestr[0:4])
					month, _ := strconv.Atoi(datestr[4:6])
					mon := time.Month(month)
					day, _ := strconv.Atoi(datestr[6:8])
					hour, _ := strconv.Atoi(datestr[9:11])
					min, _ := strconv.Atoi(datestr[11:13])
					sec, _ := strconv.Atoi(datestr[13:15])
					eventdate := time.Date(year, mon, day, hour, min, sec, 0, time.UTC)
					tp.EventDate = eventdate
				}
			} else if strcmp_from_beginning_of_line(sline, "SUMMARY:") {
				idx := strings.Index(sline, ":")
				tp.EventName = sline[idx+1:]
			}
		}
		return tp
	}
	return nil
}
