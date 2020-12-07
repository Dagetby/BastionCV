package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCounter(t *testing.T) {
	var qq = flag.String("sub", "", "Substring to search in the input url")
	fmt.Println(qq)
	q := "Go"
	cases := []string{"/ok", "/err", "/"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			// Count: 4
			fmt.Fprint(w, "Go Go Gophers Gophers go go gophers gophers")
		case "/err":
			// Bad request
			w.WriteHeader(http.StatusBadRequest)
		default:
			// Count: 0
			fmt.Fprint(w, "go go gophers gophers")

		}
	}))
	defer ts.Close()

	for _, val := range cases {
		url := ts.URL + val
		if strings.EqualFold(val, "/ok") {
			result, err := Worker(url, q)

			if result != 4 && err != nil {
				fmt.Printf("Fail OK test URL: %s, expected result: 2, got: %d \n", url, result)
				t.Fail()
				continue
			}
		}

		if strings.EqualFold(val, "/err") {
			result, err := Worker(url, q)
			if result != -1 && err == nil {
				fmt.Printf("Fail errStatusCod test URL: %s, expected result: Error on statusCode, got: %d \n", url, result)
				t.Fail()
				continue
			}
		}

		if strings.EqualFold(val, "/") {
			result, err := Worker(url, q)
			if result != -1 && err != nil {
				fmt.Printf("Fail defaut test URL: %s, expected result: 0, got: %d \n", url, result)
				t.Fail()
				continue
			}
		}

		result, err := Worker("ERROR URL", q)
		if result != -1 && err != nil {
			fmt.Printf("Fail errorURL test URL: %s, expected result: 0, got: %d \n", url, result)
			t.Fail()
			continue
		}

		fmt.Printf("Test: %s  OK\n", url)
	}

}
