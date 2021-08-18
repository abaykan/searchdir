package main

import (
	h "searchdir/helpers"
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

func show_banner() {
	fmt.Printf(`
     ┌─┐┌─┐┌─┐┬─┐┌─┐┬ ┬┌┬┐┬┬─┐
     └─┐├┤ ├─┤├┬┘│  ├─┤ │││├┬┘
     └─┘└─┘┴ ┴┴└─└─┘┴ ┴─┴┘┴┴└─
            github.com/abaykan

`)
}

func gassss(linknya string, eks string, randomAgent bool) {
	var file, err = os.OpenFile("db/dicc.txt", os.O_RDONLY, 0644)
	if h.IsError(err) {
		return
	}
	defer file.Close()

	var text = make([]byte, 133285)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if h.IsError(err) {
				break
			}
		}
		if n == 0 {
			break
		}
	}
	if h.IsError(err) {
		return
	}

	fmt.Println("Starting~")
	pecah := strings.Split(string(text), "\n")
	s := spinner.New(spinner.CharSets[4], 50*time.Millisecond)
	for i := 0; i < (len(pecah) - 1); i++ {
		s.Start()
		time.Sleep(time.Millisecond)

		if linknya[len(linknya)-1:] == "/" {
			linknya = linknya[:len(linknya)-1]
		}

		if pecah[i][0:1] != "/" {
			pecah[i] = h.VarFormat("{{.}}"+pecah[i], "/")
		}

		// ada %EXT%
		if strings.Index(pecah[i], "%EXT%") > -1 {
			// eks ada koma
			if strings.Index(eks, ",") > -1 {
				pecaheks := strings.Split(eks, ",")
				for o := 0; o < len(pecaheks); o++ {
					pathwithekstensi := strings.Replace(pecah[i], "%EXT%", pecaheks[o], 1)
					fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pathwithekstensi))
					if h.IsError(err) {
						return
					}
					if randomAgent {
						h.Rikues(fixurl.String(), true)
					} else {
						h.Rikues(fixurl.String(), false)
					}
				}
			} else {
				pecah[i] = strings.Replace(pecah[i], "%EXT%", eks, 1)
				fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pecah[i]))
				if h.IsError(err) {
					return
				}
				if randomAgent {
					h.Rikues(fixurl.String(), true)
				} else {
					h.Rikues(fixurl.String(), false)
				}
			}

		} else {
			fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pecah[i]))
			if h.IsError(err) {
				return
			}
			if randomAgent {
				h.Rikues(fixurl.String(), true)
			} else {
				h.Rikues(fixurl.String(), false)
			}
		}
	}
	s.Stop()
}

func main() {
	show_banner()

	target := flag.String("u", "", "URL Target (Required)")
	eks := flag.String("e", "", "Extension: php,css,js,etc. (Required)")
	random_agent := flag.Bool("random-agent", false, "Use random user-agent")

	flag.Parse()

	if !h.FlagPassed("u") {
		fmt.Println("[!] URL Target is missing. Use -h to show usage.")
		os.Exit(0)
	}
	if !h.FlagPassed("e") {
		fmt.Println("[!] Extensions is missing. Use -h to show usage.")
		os.Exit(0)
	}
	if !h.ValidUrl(*target) {
		fmt.Println("[!] Target isn't a valid URL. Use -h to show usage.")
		os.Exit(0)
	}

	fmt.Printf("Target: %s \n", *target)
	if strings.Index(*eks, ",") > -1 {
		fmt.Printf("Extensions: %s \n\n", *eks)
	} else {
		fmt.Printf("Extension: %s \n\n", *eks)
	}

	if *random_agent {
		gassss(*target, *eks, true)
	} else {
		gassss(*target, *eks, false)
	}
}
