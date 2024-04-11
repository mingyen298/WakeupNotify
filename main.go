package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	lineUrl = "https://notify-api.line.me/api/notify"
)

var tokens = []string{""}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <magnitude> <seconds>")
		return
	}

	magnitude := strings.Replace(os.Args[1], "+", "強", -1)
	magnitude = strings.Replace(magnitude, "-", "弱", -1)
	seconds := os.Args[2]

	msg := fmt.Sprintf("警告：地區預計震度%s級地震\n預計到達時間:%s秒", magnitude, seconds)
	wg := &sync.WaitGroup{}

	for _, token := range tokens {
		wg.Add(1)
		go notify(wg, msg, token)
	}
	wg.Wait()
}

func notify(wg *sync.WaitGroup, msg string, token string) {

	client := &http.Client{}
	data := "message=" + msg
	req, err := http.NewRequest("POST", lineUrl, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer func() {
		resp.Body.Close()
		wg.Done()
	}()

	fmt.Println("Message sent successfully")

}
