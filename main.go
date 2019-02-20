package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Atbash struct {
	Text string `json:"text"`
}

type CaesarText struct {
	Text string `json:"text"`
}

func caesar(r rune, shift int) rune {
	if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= 'А' && r <= 'Я' || r >= 'а' && r <= 'я' {
		s := int(r) + shift
		if r >= 'A' && r <= 'Z' {

			if s > 'Z' {
				return rune(s - 26)
			} else if s < 'A' {
				return rune(s + 26)
			}
			return rune(s)
		}
		if r >= 'a' && r <= 'z' {
			if s > 'z' {
				return rune(s - 26)
			} else if s < 'a' {
				return rune(s + 26)
			}
		}
		if r >= 'А' && r <= 'Я' {
			if s > 'Я' {
				return rune(s - 32)
			} else if s < 'А' {
				return rune(s + 32)
			}
		}
		if r >= 'а' && r <= 'я' {
			if s > 'я' {
				return rune(s - 32)
			} else if s < 'а' {
				return rune(s + 32)
			}
		}
		return rune(s)
	}
	return rune(r)
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
		caesarText := strings.Map(func(r rune) rune {
			return caesar(r, key)
		}, text)
		resp, _ := json.Marshal(CaesarText{caesarText})
		writer.Write(resp)
	} else {
		caesarText := strings.Map(func(r rune) rune {
			return caesar(r, -key)
		}, text)
		resp, _ := json.Marshal(CaesarText{caesarText})
		writer.Write(resp)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
