package main

import (
	"log"
	"net/http"
	"time"
	"crypto/tls"

	"os"
	_ "github.com/heroku/x/hmetrics/onload"
)

func site() (string) {
	return "https://cdracamle.herokuapp.com/"
	// return "/"
}

func getCat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!doctype html><html><body><p>TEST!</p></body></html>`))
}

func sendBear(w http.ResponseWriter, r *http.Request) {
	for "321" = "321" {
		time.Sleep(10 * time.Second)
		req, err := http.NewRequest(http.MethodDelete,
			"https://elmacards.herokuapp.com/getbot", nil)
		if err == nil {
			tr := &http.Transport{
        			TLSClientConfig: &tls.Config{
        	    			InsecureSkipVerify: true,
        			},
    			}
    			client := &http.Client{
        			Transport: tr,
        			Timeout:   20 * time.Second,
    			}
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
}

func main() {
	http.HandleFunc("/sendbot", sendBear)
	http.HandleFunc("/getbot", getCat)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	http.ListenAndServe(":"+port, nil)

	log.Println("starting server at :8080")
	//http.ListenAndServe(":8080", nil)
}
