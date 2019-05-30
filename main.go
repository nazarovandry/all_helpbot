package main

import (
	"log"
	"net/http"
	"time"
	"crypto/tls"
	"sync"

	"os"
	_ "github.com/heroku/x/hmetrics/onload"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!doctype html><html><body><p>Main Page</p></body></html>`))
}

func getCat(w http.ResponseWriter, r *http.Request, bot *int, botLock *bool, mu *sync.Mutex) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!doctype html><html><body><p>TEST!</p></body></html>`))
	mu.Lock()
	if *bot > 0 {
		*bot = 0
	} else {
		*bot -= 1
	}
	if *bot < -2 {
		*bot = 0
		*botLock = true
		mu.Unlock()
		log.Println("HELP!")
		go sendBear(w, r, bot, mu)
		mu.Lock()
	}
	mu.Unlock()
	log.Println("getmess-OK ", *bot)
}

func sendBear(w http.ResponseWriter, r *http.Request, bot *int, botLock *bool, mu *sync.Mutex) {
	mu.Lock()
	if !(*botLock) {
		mu.Unlock()
		log.Println("No pls!")
		return
	}
	*botLock = false
	for {
		time.Sleep(3 * time.Minute)
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
				mu.Lock()
				*bot = 0
				mu.Unlock()
				log.Println("client error: " + err.Error())
			} else {
				mu.Lock()
				*bot = 1
				mu.Unlock()
				log.Println("sendmess-DONE ", *bot)
			}
		} else {
			mu.Lock()
			*bot = 0
			mu.Unlock()
			log.Println("request error: " + err.Error())
		}
	}
}

func main() {
	botInt := 0
	bot := &botInt
	botLockBool := true
	botLock := &botLockBool
	mu := &sync.Mutex{}
	
	http.HandleFunc("/sendbot", func(w http.ResponseWriter, r *http.Request) {
		sendBear(w, r, bot, botLock, mu)
	})

	http.HandleFunc("/getbot", func(w http.ResponseWriter, r *http.Request) { 
		getCat(w, r, bot, botLock, mu)
	})

	http.HandleFunc("/", mainPage)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.ListenAndServe(":"+port, nil)

	log.Println("starting server at :8081")
	//http.ListenAndServe(":8081", nil)
}
