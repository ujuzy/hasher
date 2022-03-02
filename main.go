package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	httpPrefix = "http://"
)

func main() {
	goroutinesCount := flag.Int64("parallel", 10, "specifies max parallel requests")
	flag.Parse()

	args := sanitizeArgs(os.Args[1:])
	urlsGroups := splitUrls(args, *goroutinesCount)
	hashes := getHashes(urlsGroups)

	for _, hash := range hashes {
		println(hash)
	}
}

func getHashes(urlsByGroup [][]string) []string {
	result := make([]string, 0)

	hashCh := make(chan []string, len(urlsByGroup))
	errorsCh := make(chan error)
	var wg sync.WaitGroup

	for _, urls := range urlsByGroup {
		wg.Add(1)
		go makeRequests(urls, hashCh, errorsCh, &wg)
	}

	go func() {
		wg.Wait()
		close(hashCh)
		close(errorsCh)
	}()

	for err := range errorsCh {
		println(err.Error())
	}

	for hashes := range hashCh {
		result = append(result, hashes...)
	}

	return result
}

func splitUrls(urls []string, groupsCount int64) [][]string {
	goroutinesData := make([][]string, groupsCount)

	for idx, url := range urls {
		groupNumber := int64(idx) % groupsCount
		goroutinesData[groupNumber] = append(goroutinesData[groupNumber], url)
	}

	return goroutinesData
}

func makeRequests(urls []string, hashesCh chan []string, errors chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	hashes := make([]string, 0, len(urls))
	for _, url := range urls {
		resp, err := makeRequest(url)
		if err != nil {
			errors <- err
			return
		}

		hash := md5.Sum(resp)
		hashes = append(hashes, url+" "+hex.EncodeToString(hash[:]))
	}

	hashesCh <- hashes

	return
}

func makeRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func addUrlPrefixIfNotExists(url string) string {
	if !strings.Contains(url, httpPrefix) {
		return httpPrefix + url
	}

	return url
}

func sanitizeArgs(args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		arg = addUrlPrefixIfNotExists(arg)
		if _, err := url.Parse(arg); err == nil && strings.Contains(arg, ".") {
			result = append(result, arg)
		}
	}
	return result
}
