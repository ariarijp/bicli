package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

func getURLsFromFile(fileName string) ([]string, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	urls := []string{}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func shorten(login string, apiKey string, longURL string) (string, error) {
	urlFormat := "http://api.bit.ly/v3/shorten?login=%s&apiKey=%s&format=json&longUrl=%s"
	longURL = url.QueryEscape(longURL)
	apiURL := fmt.Sprintf(urlFormat, login, apiKey, longURL)

	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var bitlyResp BitlyResponse
	err = json.Unmarshal(body, &bitlyResp)
	if err != nil {
		var bitlyErrResp BitlyErrorResponse
		err = json.Unmarshal(body, &bitlyErrResp)
		if err != nil {
			return "", err
		}

		return bitlyErrResp.StatusText, nil
	}

	return string(bitlyResp.Data.URL), nil
}

func makeConfigFile(fileName string) error {
	login := scan("Put your Bitly login name: ")
	apiKey := scan("Put your Bitly API Key: ")
	conf := BicliConfig{
		Login:  login,
		APIKey: apiKey,
	}

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(conf)
	if err != nil {
		return err
	}

	ioutil.WriteFile(fileName, buffer.Bytes(), os.ModePerm)

	return nil
}

func readConfigFile(fileName string) (*BicliConfig, error) {
	var conf BicliConfig
	_, err := toml.DecodeFile(fileName, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func scan(msg string) string {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	init := flag.Bool("init", false, "Create config file")
	configFileName := flag.String("conf", "config.toml", "Config file")
	urlsFileName := flag.String("urls", "urls.csv", "Long URL urls file")
	sep := flag.String("sep", ",", "Output separator")
	sleepMsec := flag.Uint("sleep-msec", 1000, "Sleep time for each request")
	flag.Parse()

	if *init {
		err := makeConfigFile(*configFileName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Config file created")

		return
	}

	conf, err := readConfigFile(*configFileName)
	if err != nil {
		panic(err)
	}

	urls, err := getURLsFromFile(*urlsFileName)
	if err != nil {
		panic(err)
	}
	shortURLs := []ShortURL{}

	var wg sync.WaitGroup
	for i, longURL := range urls {
		wg.Add(1)
		go func(i int, longURL string) {
			defer wg.Done()

			_url, err := shorten(conf.Login, conf.APIKey, longURL)
			if err != nil {
				_url = fmt.Sprintf("%v", err)
			}

			shortURL := ShortURL{
				LineNum: i + 1,
				URL:     _url,
				LongURL: longURL,
			}
			shortURLs = append(shortURLs, shortURL)
		}(i, longURL)

		time.Sleep(time.Duration(*sleepMsec) * time.Millisecond)
	}
	wg.Wait()

	sort.Sort(ShortURLs(shortURLs))
	for _, s := range shortURLs {
		fmt.Println(s.ToCSV(*sep))
	}
}
