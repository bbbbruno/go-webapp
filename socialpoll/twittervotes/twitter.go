package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/oauth1/oauth"
	"github.com/joeshaw/envdecode"
)

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

func setupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AcessToken     string `env:"SP_TWITTER_ACCESS_TOKEN"`
		AcessSecret    string `env:"SP_TWITTER_ACCESS_SECRET"`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	creds = &oauth.Credentials{
		Token:  ts.AcessToken,
		Secret: ts.AcessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(query url.Values) (*http.Request, error) {
	authSetupOnce.Do(func() {
		setupTwitterAuth()
	})
	const endpoint = "https://stream.twitter.com/1.1/statuses/filter.json"

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(query.Encode()))
	if err != nil {
		return nil, err
	}

	formEnc := query.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	ah := authClient.AuthorizationHeader(creds, "POST", req.URL, query)
	req.Header.Set("Authorization", ah)

	return req, nil
}

type tweet struct {
	Text string
}

func twitterStream(ctx context.Context, votes chan<- string) {
	defer close(votes)
	for {
		log.Println("Twitterに問い合わせます...")
		readFromTwitterWithTimeout(ctx, 1*time.Minute, votes)
		log.Println(" (待機中)")
		select {
		case <-ctx.Done():
			log.Println("Twitterへの問い合わせを終了します...")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func readFromTwitterWithTimeout(ctx context.Context, timeout time.Duration, votes chan<- string) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	readFromTwitter(ctx, votes)
}

func readFromTwitter(ctx context.Context, votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("選択肢の読み込みに失敗しました: ", err)
		return
	}

	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := makeRequest(query)
	if err != nil {
		log.Println("検索のリクエストに失敗しました: ", err)
		return
	}

	client := &http.Client{}
	if deadline, ok := ctx.Deadline(); ok {
		client.Timeout = time.Until(deadline)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("検索のリクエストに失敗しました: ", err)
		return
	}
	done := make(chan struct{})
	defer func() { <-done }()

	defer resp.Body.Close()
	go func() {
		defer close(done)
		log.Println("resp: ", resp.StatusCode)
		if resp.StatusCode != 200 {
			var buf bytes.Buffer
			io.Copy(&buf, resp.Body)
			log.Println("resp body: ", buf.String())
			return
		}
		decoder := json.NewDecoder(resp.Body)
		for {
			var tweet tweet
			if err := decoder.Decode(&tweet); err != nil {
				break
			}
			log.Println("tweet: ", tweet)
			for _, option := range options {
				if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(option)) {
					log.Println("投票: ", option)
					votes <- option
				}
			}
		}
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}
