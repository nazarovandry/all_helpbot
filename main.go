package main

import (
	"log"
	"net/http"
	"time"

	"os"
	_ "github.com/heroku/x/hmetrics/onload"
)

func site() (string) {
	return "https://cdracamle.herokuapp.com/"
	// return "/"
}

func sendMess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!doctype html><html><body><p>TEST!</p></body></html>`))
	w.WriteHeader(http.StatusOK)
	req, err := http.NewRequest(http.MethodDelete,
		"https://elmacards.herokuapp.com/events", nil)
	if err == nil {
		client := &http.Client{Timeout:	2 * time.Second}
		_, err := client.Do(req)
		if err != nil {
			log.Println("client error: " + err.Error())
		} else {
			log.Println("sanmess-DONE")
		}
	} else {
		log.Println("request error: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", sendMess)
	//http.HandleFunc("/2", getMess)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	http.ListenAndServe(":"+port, nil)

	log.Println("starting server at :8080")
	//http.ListenAndServe(":8080", nil)
}
