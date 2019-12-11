package main

import "net/http"

func main() {
	http.HandleFunc("/healthcheck", healthcheck)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
