package main

import "net/http"

func main() {}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValiedAPIKey(r.URL.Query().Get("Key")) {
			//respondErr(w, r, http.StatusUnauthorized, "不正なAPIキーです")
			return
		}
		fn(w, r)
	}
}

func isValiedAPIKey(key string) bool {
	return key == "abc123"
}
