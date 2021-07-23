package main

import (
	c "belajar-golang/helpers"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net/http"
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

`)
}

func getDictionary(linknya string) {
	show_banner()
	// buka file
	var file, err = os.OpenFile("db/dicc.txt", os.O_RDONLY, 0644)
	if c.IsError(err) {
		return
	}
	defer file.Close()

	// baca file
	var text = make([]byte, 133285)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if c.IsError(err) {
				break
			}
		}
		if n == 0 {
			break
		}
	}
	if c.IsError(err) {
		return
	}

	pecah := strings.Split(string(text), "\n")
	s := spinner.New(spinner.CharSets[4], 50*time.Millisecond)
	for i := 0; i < (len(pecah) - 1); i++ {
		s.Start()
		time.Sleep(time.Millisecond)
		s.Suffix = " ~processing"
		if linknya[len(linknya)-1:] == "/" {
			linknya = linknya[:len(linknya)-1]
		}
		if pecah[i][0:1] != "/" {
			pecah[i] = c.VarFormat("{{.}}"+pecah[i], "/")
		}
		fixurl, err := url.Parse(c.VarFormat(linknya+"{{.}}", pecah[i]))
		if c.IsError(err) {
			return
		}
		resp, err := http.Get(fixurl.String())
		if c.IsError(err) {
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if c.IsError(err) {
			return
		}
		sb := string(body)
		if resp.StatusCode == 404 {
			continue
		}
		fmt.Print("[", time.Now().Format("03:04:05"), "] ", resp.StatusCode, " --- ", c.LenReadable(len(sb), 1), "\t", fixurl, "\n")
	}
	s.Stop()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go {url}")
		os.Exit(0)
	}
	getDictionary(os.Args[1])
}
