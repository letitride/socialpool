package main

import (
	"flag"
	"github.com/tylerstillwater/graceful"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "エンドポイントのアドレス")
		mongo = flag.String("mongo", "localhost", "MongoDBのアドレス")
	)
	flag.Parse()
	log.Println("MongoDBに接続します", *mongo)
	mongoInfo := &mgo.DialInfo{
		Addrs:    []string{*mongo},
		Username: "root",
		Password: "example",
	}
	db, err := mgo.DialWithInfo(mongoInfo)
	if err != nil {
		log.Fatalln("MongoDBへの接続に失敗しました:", err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Webサーバを開始します:", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("停止します")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValiedAPIKey(r.URL.Query().Get("Key")) {
			respondErr(w, r, http.StatusUnauthorized, "不正なAPIキーです")
			return
		}
		fn(w, r)
	}
}

func isValiedAPIKey(key string) bool {
	return key == "abc123"
}

func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
