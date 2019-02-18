package main

import (
	"bytes"
	"encoding/json"
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

const lowerCaseAlphabet = "abcdefghijklmnopqrstuvwxyz"
const upperCaseAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Encrypts a plaintext string by shifting each character with the provided key.
func EncryptPlaintext(plaintext string, key int) string {
	return rotateText(plaintext, key)
}

// Decrypts a ciphertext string by reverse shifting each character with the provided key.
func DecryptCiphertext(ciphertext string, key int) string {
	return rotateText(ciphertext, -key)
}

// Takes a string and rotates each character by the provided amount.
func rotateText(inputText string, rot int) string {
	rot %= 26
	rotatedText := []byte(inputText)

	for index, byteValue := range rotatedText {
		if byteValue >= 'a' && byteValue <= 'z' {
			rotatedText[index] = lowerCaseAlphabet[(int((26+(byteValue-'a')))+rot)%26]
		} else if byteValue >= 'A' && byteValue <= 'Z' {
			rotatedText[index] = upperCaseAlphabet[(int((26+(byteValue-'A')))+rot)%26]
		}
	}
	return string(rotatedText)
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
	query := request.URL.Query()
	text := query.Get("text")
	key, _ := strconv.Atoi(query.Get("key"))
	encrypt, _ := strconv.ParseBool(query.Get("encrypt"))
	if encrypt {
		caesarText := EncryptPlaintext(text, key)
		resp, _ := json.Marshal(CaesarText{caesarText})
		writer.Write(resp)
	} else {
		caesarText := DecryptCiphertext(text, key)
		resp, _ := json.Marshal(CaesarText{caesarText})
		writer.Write(resp)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
