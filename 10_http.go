package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!")
	w.Write([]byte("!!!"))
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	//http.ListenAndServe("31.170.160.61:8080", nil)
	http.ListenAndServe(":8080", nil)
}
