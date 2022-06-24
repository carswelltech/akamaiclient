package edgegrid

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	akamaiTimeURL string = `http://time.akamai.com/?iso`
	akamaiTimeFmt string = `2006-01-02T15:04:05Z`
)

var (
	useLocalTime bool = true
)

func UseLocalTime(v bool) {
	useLocalTime = v
	switch {
	case useLocalTime:
		fmt.Fprintln(os.Stderr, "akamaiclient edgegrid.makeEdgeTimeStamp() using system clock")
	default:
		fmt.Fprintln(os.Stderr, "akamaiclient edgegrid.makeEdgeTimeStamp() using remote time fix")
	}
}

func makeEdgeTimeStamp_TimeFix() string {

	resp, err := http.Get(akamaiTimeURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	t, err := time.Parse(akamaiTimeFmt, string(body))
	if err != nil {
		t = time.Now().In(time.FixedZone("GMT", 0))
		panic(err)
	}
	t = t.In(time.FixedZone("GMT", 0))
	return fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d+0000",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
