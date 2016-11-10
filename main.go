package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/bmatcuk/doublestar"
)

const (
	IdeCmd  = "atom"
	ExtList = "js,jsx,css,sass,scss"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var (
		interval float64
		wd       string
	)

	flag.Float64Var(&interval, "interval", 5, "Refresh interval in minutes")
	flag.StringVar(&wd, "path", cwd, "Working directory")

	flag.Parse()

	pattern := path.Join(wd, "**/*.{"+ExtList+"}")
	fmt.Println(pattern)
	matches, err := doublestar.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	total := len(matches)

	log.Println("Interval:", interval, "minutes")
	log.Println("Working directory:", wd)
	log.Println("Files founded:", total)

	tickerInterval := time.Duration(interval * float64(time.Minute))
	ticker := time.NewTicker(tickerInterval)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _ = range ticker.C {
		cmd := exec.Command(IdeCmd, matches[r.Intn(total)])
		err := cmd.Start()
		if err != nil {
			log.Panicln(err)
		}
	}
}
