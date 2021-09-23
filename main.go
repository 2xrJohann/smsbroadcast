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
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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

func listen(ch chan string, username string, password string, callerID string, msg string, wg *sync.WaitGroup, outTE *walk.TextEdit, output *string) {
	for num := range ch {
		wg.Add(1)
		go makeReq(username, password, num, callerID, msg, wg, outTE, output)
	}
	wg.Done()
	wg.Wait()
}

func makeReq(username string, password string, callerID string, to string, msg string, wg *sync.WaitGroup, outTE *walk.TextEdit, output *string) {
	defer wg.Done()
	url := fmt.Sprintf("https://api.smsbroadcast.com.au/api-adv.php?username=%s&password=%s&to=%s&from=%s&message=%s&ref=112233&maxsplit=5", username, password, callerID, to, url.QueryEscape(msg))
	
	//commented out url below is the old url, it still works lol
	//url := fmt.Sprintf("https://www.smsbroadcast.com.au/api-adv.php?username=%s&password=%s&to=%s&from=%s&message=%s&ref=112233&maxsplit=5", username, password, callerID, to, url.QueryEscape(msg))
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	sb := string(body)
	sb += "\r\n"
	*output+=sb
}

func main() {

	var wg sync.WaitGroup
	output := ""
	username := "" //redacted
	password := "" //redacted
	callerID := "" //redacted
	numbers := make(chan string)
	var inTE, outTE *walk.TextEdit

	wg.Add(1)
	go loadSheet(numbers, "sheet.csv", &wg)

	MainWindow{
		Title:   "sned txt",
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "sned txt",
				OnClicked: func() {
					wg.Add(1)
					listen(numbers, username, password, callerID, inTE.Text(), &wg, outTE, &output)
					outTE.SetText(output)
				},
			},
		},
	}.Run()
}
