package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/bmatcuk/doublestar"
)

const (
	IdeCmd = "atom"
)

func CheckFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var (
		interval float64
		wd       string
		glob     string
		re2      string
	)

	flag.Float64Var(&interval, "interval", 5, "Refresh interval in minutes")
	flag.StringVar(&wd, "path", cwd, "Working directory")

	flag.StringVar(&glob, "glob", "**/*.{js,go}", "Glob for file paths")
	flag.StringVar(&re2, "re2", "", "Re2 regexp for file contents. Example: \"(?i)caseinsensitivesubstring\" ")

	flag.Parse()

	fmt.Println(glob)
	matches, err := doublestar.Glob(glob)
	CheckFatalError(err)

	if len(re2) > 0 {
		var (
			m []string
		)

		r, err := regexp.Compile(re2)
		CheckFatalError(err)

		for _, filename := range matches {
			content, err := ioutil.ReadFile(filename)
			CheckFatalError(err)
			if r.Match(content) {
				m = append(m, filename)
			}
		}

		matches = m
	}

	total := len(matches)

	log.Println("Interval:", interval, "minutes")
	log.Println("Working directory:", wd)
	log.Println("Glob:", glob)
	if len(re2) > 0 {
		log.Println("Re2 regexp:", re2)
	}
	log.Println("Files founded:", total)

	tickerInterval := time.Duration(interval * float64(time.Minute))
	ticker := time.NewTicker(tickerInterval)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _ = range ticker.C {
		name := matches[r.Intn(total)]
		cmd := exec.Command(IdeCmd, name)
		err := cmd.Start()
		if err != nil {
			log.Panicln(err)
		}

		log.Println("->", name)
	}
}
