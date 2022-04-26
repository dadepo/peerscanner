package scanner

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func CorsScan(ips []string, resChan chan<- map[string][][]string) {
	var wg sync.WaitGroup

	// for each ip
	// add to wait group
	// find cors setting in a go routine
	for _, ip := range ips {
		wg.Add(1)
		fmt.Println("scanning:", ip)
		go findCorsSetting(&wg, resChan, ip)
	}

	// wait till all goroutines are done then close channel
	go func() {
		wg.Wait()
		close(resChan)
	}()
}

func DnsScan(ips []string, resChan chan<- map[string]string) {
	var wg sync.WaitGroup
	// for each ip
	// add to wait group
	// find cors setting in a go routine
	for _, ip := range ips {
		wg.Add(1)
		fmt.Println("scanning DNS:", ip)
		go func(ipp string) {
			defer wg.Done()
			// recover from the error
			defer func() {
				if e := recover(); e != nil {
					log.Fatal(e)
				}
			}()

			addr, err := net.LookupAddr(ipp)
			if err == nil {
				resChan <- map[string]string{ip: addr[0]}
			} else {
				//resChan <- map[string]string{ip: "---"}
				//log.Println(err)
			}
		}(ip)
	}

	// wait till all goroutines are done then close channel
	go func() {
		wg.Wait()
		close(resChan)
	}()
}

func findCorsSetting(wg *sync.WaitGroup, c chan<- map[string][][]string, ip string) {
	defer wg.Done()
	// recover from the error
	defer func() {
		if e := recover(); e != nil {
			log.Fatal(e)
		}
	}()

	client := http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("OPTIONS", fmt.Sprintf("http://%s:%d/api/v0/add", ip, 5001), nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Origin", "*")

	resp, err := client.Do(request)
	if err == nil {
		c <- map[string][][]string{
			ip: {
				resp.Header["Access-Control-Allow-Origin"],
				resp.Header["Headers.Access-Control-Allow-Credentials"],
			},
		}
	} else {
		log.Println(err)
	}
}
