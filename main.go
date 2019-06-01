package main

import (
	"log"
	"net/http"
	"time"
	"crypto/tls"
	"sync"
	"io/ioutil"
	"bytes"

	"os"
	_ "github.com/heroku/x/hmetrics/onload"
)

func otherSite() (string) {
	return "https://elmacards.herokuapp.com/"
	//return "http://127.0.0.1:8080/"
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!doctype html><html><body><p>MainPage</p></body></html>`))
}

func getCat(w http.ResponseWriter, r *http.Request, bot *int,
	botLock *bool, mu *sync.Mutex, savings *[]byte) {
	out, err := ioutil.ReadAll(r.Body)
	if err == nil {
		log.Println("I have got: '" + string(out)[:50] + "'!")
		*savings = out
	} else {
		log.Println(err.Error())
	}
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
		go sendBear(w, r, bot, botLock, mu, savings)
		mu.Lock()
	}
	mu.Unlock()
	log.Println("getmess-OK ", *bot)
}

func sendBear(w http.ResponseWriter, r *http.Request, bot *int,
	botLock *bool, mu *sync.Mutex, savings*[]byte) {
	mu.Lock()
	if !(*botLock) {
		mu.Unlock()
		log.Println("No pls!")
		return
	}
	*botLock = false
	mu.Unlock()
	for {
		//time.Sleep(3 * time.Minute)
		time.Sleep(10 * time.Second)
		mu.Lock()
		data := bytes.NewReader(*savings)
		mu.Unlock()
		req, err := http.NewRequest(http.MethodDelete,
			otherSite() + "getbot", data)
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
	byt := []byte("(-BLOCK-)\n\nAndrY(-ELEM-)andry(-ELEM-)(-STRING-)(-BLOCK-)\n\n1922-06-31 20:07:42(-ELEM-)AndrY(-ELEM-)Wait, the database is not ready yet...(-STRING-)(-BLOCK-)\n\nSite is broken. Soon it will be fixed.(-BLOCK-)")
	savings := &byt
	
	http.HandleFunc("/sendbot", func(w http.ResponseWriter, r *http.Request) {
		sendBear(w, r, bot, botLock, mu, savings)
	})

	http.HandleFunc("/getbot", func(w http.ResponseWriter, r *http.Request) { 
		getCat(w, r, bot, botLock, mu, savings)
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
