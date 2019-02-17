package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Atbash struct {
	Name string `json:"name"`
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	fmt.Println("GET params were:", request.URL.Query())
	text := request.URL.Query().Get("text")
	var cipherText bytes.Buffer
	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			cipherText.WriteRune(rune('A') + rune('Z') - char)
		} else if char >= 'a' && char <= 'z' {
			cipherText.WriteRune(rune('a') + rune('z') - char)
		} else if char >= 'А' && char <= 'Я' {
			cipherText.WriteRune(rune('А') + rune('Я') - char)
		} else if char >= 'а' && char <= 'я' {
			cipherText.WriteRune(rune('а') + rune('я') - char)
		} else {
			cipherText.WriteRune(char)
		}
	}
	fmt.Println(cipherText.String())
	atbash := Atbash{cipherText.String()}
	resp, _ := json.Marshal(atbash)
	writer.Write(resp)
}

func main() {
	http.HandleFunc("/atbash", indexHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
