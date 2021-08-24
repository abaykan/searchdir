package helpers

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"
)

const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000
)

func LenReadable(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	if length > TB {
		unit = "TB"
		i = length / TB
		remainder = length - (i * TB)
	} else if length > GB {
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
	} else if length > MB {
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
	} else if length > KB {
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
	} else {
		return strconv.Itoa(length) + "B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else if remainder > KB {
		width = 6
	} else {
		width = 3
	}

	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d%s", i, unit)
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCtrl+C pressed")
		fmt.Println("\r- Stopping Services -")
		os.Exit(0)
	}()
}

func IsError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func VarFormat(s string, v interface{}) string {
	t, b := new(template.Template), new(strings.Builder)
	template.Must(t.Parse(s)).Execute(b, v)
	return b.String()
}

func FlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func get_random_agent() string {
	var file, _ = os.OpenFile("db/user-agents.txt", os.O_RDONLY, 0644)
	defer file.Close()

	var text = make([]byte, 133285)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if IsError(err) {
				break
			}
		}
		if n == 0 {
			break
		}
	}

	pecahagent := strings.Split(string(text), "\n")

	rand.Seed(time.Now().UTC().UnixNano())
	agent := pecahagent[rand.Intn(len(pecahagent)-1)]

	return agent
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Rikues(urlnya string, randomAgent bool, excode []string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", urlnya, nil)
	if IsError(err) {
		return
	}

	if randomAgent {
		user_agent := string(get_random_agent())
		req.Header.Set("User-Agent", user_agent)
	}

	resp, err := client.Do(req)
	if IsError(err) {
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if IsError(err) {
		return
	}

	sb := string(body)

	if Contains(excode, strconv.Itoa(resp.StatusCode)) != true {
		u, err := url.Parse(urlnya)
		if IsError(err) {
			return
		}
		WriteLog(u.Host, strconv.Itoa(resp.StatusCode)+"\t"+LenReadable(len(sb), 1)+"\t"+urlnya)
		fmt.Print("[", time.Now().Format("03:04:05"), "] ", resp.StatusCode, " --- ", LenReadable(len(sb), 1), "\t", urlnya, "\n")

	}
}

func WriteLog(domain string, log string) {
	// bikin dir
	path := "logs/"
	os.MkdirAll(path, 0755)

	// deteksi apakah file sudah ada
	var _, err = os.Stat("logs/" + domain + ".txt")

	// buat file baru jika belum ada
	if os.IsNotExist(err) {
		var file, err = os.Create("logs/" + domain + ".txt")
		if IsError(err) {
			return
		}
		defer file.Close()
	}

	// buka file dengan level akses READ & WRITE
	var file, errr = os.OpenFile("logs/"+domain+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if IsError(errr) {
		return
	}
	defer file.Close()

	// tulis data ke file
	_, err = file.WriteString(log + "\n")
	if IsError(err) {
		return
	}

	// simpan perubahan
	err = file.Sync()
	if IsError(err) {
		return
	}
}
