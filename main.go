package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func loadSheet(ch chan string, file string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, _ := os.Open(file)
	defer f.Close()
	lines, _ := csv.NewReader(f).ReadAll()
	for _, line := range lines {
		ch <- strings.TrimSpace(line[0])
	}
	close(ch)
}

func listen(ch chan string, username string, password string, callerID string, msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		wg.Add(1)
		go makeReq(username, password, num, callerID, msg, wg)
	}
}

func makeReq(username string, password string, callerID string, to string, msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://api.smsbroadcast.com.au/api-adv.php?username=%s&password=%s&to=%s&from=%s&message=%s&ref=112233&maxsplit=5", username, password, callerID, to, url.QueryEscape(msg))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
	body, _ := ioutil.ReadAll(resp.Body)
	sb := string(body)
	fmt.Println(sb)
}

func main() {

	var wg sync.WaitGroup

	username := "" //redacted
	password := "" //redacted
	callerID := "" //redacted
  
	numbers := make(chan string)

	wg.Add(1)
	go loadSheet(numbers, "sheet.csv", &wg)

	wg.Add(1)
	listen(numbers, username, password, callerID, "txt msg here", &wg)

	wg.Wait()
	fmt.Println("done")
}
