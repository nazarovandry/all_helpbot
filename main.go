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

func sendmess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!doctype html><html><body><p>TEST!</p></body></html>`))
	//if r.Method == http.MethodGet {
		req, err := http.NewRequest(http.MethodGet,
			"https://elmacards.herokuapp.com/events", nil)
		if err == nil {
			client := &http.Client{Timeout:	2 * time.Second}
			_, err := client.Do(req)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("done")
			}
		} else {
			log.Println(err.Error())
		}
		//time.Sleep(9 * time.Minute)
	//}
}

func main() {
	http.HandleFunc("/", sendmess)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	/*req, err := http.NewRequest(http.MethodGet,
		"https://elmacards.herokuapp.com/events", nil)
	if err == nil {
		client := &http.Client{Timeout:	2 * time.Second}
		_, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}*/
	
	http.ListenAndServe(":"+port, nil)

	log.Println("starting server at :8080")
	//http.ListenAndServe(":8080", nil)
}
