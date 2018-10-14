package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	logTmpl = "127.0.0.1 - %v [09/May/2018:16:00:39 +0000] \"GET %v HTTP/1.0\" %v 1234\n"
)

var (
	names = []string{
		"james",
		"jill",
		"frank",
		"mary",
		"jake",
		"lindsay",
	}

	sections = []string{
		"/report",
		"/api/user",
		"/api/user/create",
		"/foo",
		"/foo/bar/baz",
	}

	statusCodes = []string{
		"200",
		"500",
	}
)

// This is a test which writes random logs for 3 minutes then stops.
// Good for testing alerting and recovery.
func main() {
	watchFile := flag.String("src", "/var/log/access.log", "file to watch for incoming logs")
	flag.Parse()

	file, err := os.OpenFile(*watchFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	twoMinTimer := time.NewTicker(time.Minute * 3)
	tenSecTimer := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-tenSecTimer.C:
			for i := 0; i < 400; i++ {
				if _, err := file.WriteString(randomLog()); err != nil {
					log.Fatal(err)
				}
			}
		case <-twoMinTimer.C:
			log.Print("End of writes")
			return
		}
	}
}

func randomLog() string {
	index := rand.Int()

	return fmt.Sprintf(
		logTmpl,
		names[index%len(names)],
		sections[index%len(sections)],
		statusCodes[index%len(statusCodes)],
	)
}
