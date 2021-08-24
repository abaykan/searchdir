package main

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net/url"
	"os"
	h "searchdir/helpers"
	"strconv"
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

func gassss(linknya string, eks string, opsi []string) {
	randomagent, _ := strconv.ParseBool(opsi[0])

	pecah_exclude := []string{}
	all_exclude := []string{"404"}

	if strings.Index(opsi[1], "-") > -1 {
		pecah_exclude = strings.Split(opsi[1], "-")
		awal, _ := strconv.Atoi(pecah_exclude[0])
		akhir, _ := strconv.Atoi(pecah_exclude[1])
		for i := awal; i <= akhir; i++ {
			all_exclude = append(all_exclude, strconv.Itoa(i))
		}
	} else if strings.Index(opsi[1], ",") > -1 {
		for _, excl := range strings.Split(opsi[1], ",") {
			all_exclude = append(all_exclude, excl)
		}
	} else if len(opsi[1]) == 3 {
		all_exclude = append(all_exclude, opsi[1])
	} else if opsi[1] == "gausah" {
		all_exclude = all_exclude
	} else {
		fmt.Println("Weird Status Code. Please re-check your input.")
		os.Exit(0)
	}

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

		if strings.Index(pecah[i], "%EXT%") > -1 {
			if strings.Index(eks, ",") > -1 {
				pecaheks := strings.Split(eks, ",")
				for o := 0; o < len(pecaheks); o++ {
					pathwithekstensi := strings.Replace(pecah[i], "%EXT%", pecaheks[o], 1)
					fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pathwithekstensi))
					if h.IsError(err) {
						return
					}
					if randomagent {
						h.Rikues(fixurl.String(), true, all_exclude)
					} else {
						h.Rikues(fixurl.String(), false, all_exclude)
					}
				}
			} else {
				pecah[i] = strings.Replace(pecah[i], "%EXT%", eks, 1)
				fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pecah[i]))
				if h.IsError(err) {
					return
				}
				if randomagent {
					h.Rikues(fixurl.String(), true, all_exclude)
				} else {
					h.Rikues(fixurl.String(), false, all_exclude)
				}
			}

		} else {
			fixurl, err := url.Parse(h.VarFormat(linknya+"{{.}}", pecah[i]))
			if h.IsError(err) {
				return
			}
			if randomagent {
				h.Rikues(fixurl.String(), true, all_exclude)
			} else {
				h.Rikues(fixurl.String(), false, all_exclude)
			}
		}
	}
	s.Stop()
}

func main() {
	h.SetupCloseHandler()
	show_banner()

	target := flag.String("u", "", "URL Target (Required)")
	eks := flag.String("e", "", "Extension: php,css,js,etc. (Required)")
	random_agent := flag.Bool("random-agent", false, "Use random user-agent")
	exclude_status_code := flag.String("es", "", "Exclude Status Code. Ex: \"400\", \"400,403\" or \"400-500\"")

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

	opsi := make([]string, 2)

	if *random_agent {
		opsi[0] = "true"
	} else {
		opsi[0] = "false"
	}

	if h.FlagPassed("es") {
		opsi[1] = *exclude_status_code
	} else {
		opsi[1] = "gausah"
	}

	gassss(*target, *eks, opsi)
}
