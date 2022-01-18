package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

var (
	conn   net.Conn
	reader io.ReadCloser
)

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDBにダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("データベース接続が閉じられました")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) {
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	for vote := range votes {
		pub.Publish("votes", []byte(vote))
	}
	log.Println("Pubisher: 停止中です")
	pub.Stop()
	log.Println("Publisher: 停止しました")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		cancel()
		log.Println("停止します...")
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("MongoDBへのダイヤルに失敗しました: ", err)
	}
	defer closedb()

	votes := make(chan string)
	go twitterStream(ctx, votes)
	publishVotes(votes)
}
