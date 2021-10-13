package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"unicode"
	"strings"
)

type User struct {
	eapNumber string
	fName string
	trackingNumber string
	homePhone string
    multi1 string
}

type Config struct{
	username string
	password string
	callerID string
}

func loadSheet(ch chan User, file string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, _ := os.Open(file)
	defer f.Close()
	lines, _ := csv.NewReader(f).ReadAll()
	flag := 0
	for _, line := range lines {
		if flag == 0{
			flag++
			continue
		}
		var user User
		user.eapNumber = line[0]
		user.fName = line[1]
		user.trackingNumber = line[2]
		user.homePhone = line[3]
		user.multi1 = line[4]
		ch <- user
	}
	close(ch)
}

func listen(ch chan User, config Config,wg *sync.WaitGroup,) {
	for user := range ch {
		wg.Add(1)
		go makeReq(config, user, wg,)
	}
	wg.Done()
	wg.Wait()
}

func checkAlpha(str string) bool {  
	for _, charVariable := range str {  
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {  
			return false  
		}  
	}  
	return true  
} 

func makeReq(config Config, user User,wg *sync.WaitGroup) {
	defer wg.Done()

	validateNumber := func(number string) string{
		for _, char := range number {
			if unicode.IsLetter(char){
				return ""
			}
		}

		if number == "0400000000" || number == "" || number == "New Number" || number == "400000000" {
			return ""
		}else{
			return number
		}
	}

	var phoneNumber string
	if phoneNumber = validateNumber(user.homePhone); phoneNumber == ""{
		if phoneNumber = validateNumber(user.multi1); phoneNumber == ""{
			fmt.Println(user.homePhone, ", ", user.multi1, ", could not find a valid number")
			return
		}
	}

	if len(phoneNumber) != 9 && len(phoneNumber) != 10 {
		fmt.Println(phoneNumber, " :phone number is ", len(phoneNumber), " digits")
		return
	}

	if user.trackingNumber == ""{
		fmt.Println(phoneNumber, " :tracking number is empty")
		return
	}

	for _, char := range user.trackingNumber {
		if unicode.IsLetter(char){
			fmt.Println(phoneNumber, " :tracking number has letter in it")
			return
		}
	}

	message := fmt.Sprintf("Hi %s. Your %s SIM has been dispatched. From tomorrow track its progress here: auspost.com.au/mypost/track/#/details/%s", user.fName, "redacted brand name", user.trackingNumber)
	url := fmt.Sprintf("https://api.smsbroadcast.com.au/api-adv.php?username=%s&password=%s&to=%s&from=%s&message=%s&ref=112233&maxsplit=5", config.username, config.password, phoneNumber, config.callerID, url.QueryEscape(message))
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	sb := string(body)
	sb = strings.TrimSuffix(sb, "\n")
	if strings.Contains(sb, "BAD") {
		fmt.Println(phoneNumber, " failed, ", sb)
	}else{
		fmt.Println(phoneNumber, "sent successfully, ", sb)
	}
}

func main() {
	var config Config
	config.username = "" //redacted
	config.password = "" //redacted
	config.callerID = "" //redacted
	var wg sync.WaitGroup
	numbers := make(chan User)
	wg.Add(1)
	go loadSheet(numbers, "input.csv", &wg)
	wg.Add(1)
	listen(numbers, config, &wg)
	fmt.Println("done")
	for{}
}
