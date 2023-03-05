package dataserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func getLoanInfo(url string) []byte {
	defaultContext := context.Background()
	ctx, cancel := context.WithTimeout(defaultContext, 30*time.Second)
	defer cancel()
	body := "yourBody"
	userInfoRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader([]byte(body)))
	if err != nil {
		log.Fatalln(err)
	}
	client := http.DefaultClient
	res, err := client.Do(userInfoRequest)
	if err != nil {
		log.Fatalln(err)
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(resBytes))
	return resBytes
}

func main() {
	var wg sync.WaitGroup
	requestUrls := []string{"https://www.baidu.com", "https://google.com"}
	wg.Add(len(requestUrls))
	for _, requestUrls := range requestUrls {
		go func(url string) {
			getLoanInfo(url)
			wg.Done()
		}(requestUrls)
	}
	wg.Wait()
}
