package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("データベースに接続します...")
	mongoInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Username: "root",
		Password: "example",
	}
	db, err := mgo.DialWithInfo(mongoInfo)
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("データベース接続を閉じます...")
		db.Clone()
	}()
	//pollData := db.DB("boolots").C("polls")

	var countsLock sync.Mutex
	var counts map[string]int
	log.Println("NSQに接続します...")
	//nsq votes topicの講読
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}
	//講読データの処理
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		//カウンタを操作するのはgoroutineで1つのみ
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))
	//講読先のサーバ
	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}
}
