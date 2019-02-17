package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Atbash struct {
	Text string `json:"text"`
}

type CaesarText struct {
	Text string `json:"text"`
}

type Caesar struct {
	key int
}

func NewCaesar(key int) *Caesar {
	return &Caesar{key}
}

// Encipher enciphers string using Caesar cipher according to key.
func (c *Caesar) Encipher(text string) string {
	return caesarEncipher(text, c.key)
}

// Decipher deciphers string using Caesar cipher according to key.
func (c *Caesar) Decipher(text string) string {
	return caesarEncipher(text, -c.key)
}

func caesarEncipher(text string, key int) string {
	if key == 0 {
		return text
	}
	return mapAlpha(text, func(i, char int) int {
		return char + key
	})
}

func mod(a int, b int) int {
	return (a%b + b) % b
}

func mapAlpha(text string, f func(i, char int) int) string {
	runes := []rune(text)
	for i, char := range runes {
		if char >= 'A' && char <= 'Z' {
			runes[i] = rune(mod(f(i, int(char-'A')), 26)) + 'A'
		} else if char >= 'a' && char <= 'z' {
			runes[i] = rune(mod(f(i, int(char-'a')), 26)) + 'a'
		}
	}
	return string(runes)
}

func main() {
	http.HandleFunc("/atbash", atbashHandler)
	http.HandleFunc("/caesar", caesarHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func atbashHandler(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
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
	atbash := Atbash{cipherText.String()}
	resp, _ := json.Marshal(atbash)
	writer.Write(resp)
}

func caesarHandler(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	fmt.Println("GET params were:", request.URL.Query())
	text := request.URL.Query().Get("text")
	key := request.URL.Query().Get("key")
	ikey, _ := strconv.Atoi(key)
	caesarText := CaesarText{NewCaesar(ikey).Encipher(text)}
	resp, _ := json.Marshal(caesarText)
	writer.Write(resp)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
