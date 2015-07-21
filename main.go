package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	lineBreak  = "--------------------------------------------------------------------------------"
	portFlag   = flag.Int("port", 5555, "Port number")
	headerFlag = flag.Bool("header", true, "Print headers")
	bodyFlag   = flag.Bool("body", false, "Print body")
)

func main() {
	flag.Parse()
	h := handlers()
	addr := fmt.Sprintf(":%d", *portFlag)
	err := http.ListenAndServe(addr, wrap(h...))
	if err != nil {
		fmt.Println(err)
	}
}

func handlers() []http.HandlerFunc {
	handlers := []http.HandlerFunc{}
	if *headerFlag {
		handlers = append(handlers, writeHeaders)
	}
	if *bodyFlag {
		handlers = append(handlers, writeBody)
	}
	return handlers
}

func wrap(funcs ...http.HandlerFunc) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		for _, f := range funcs {
			f(w, r)
		}
		fmt.Println(lineBreak)
	}
	return http.HandlerFunc(handler)
}

func writeHeaders(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		values := strings.Join(v, ", ")
		fmt.Printf("[%s] = %s\n", k, values)
	}
}

func writeBody(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ERR: ", err)
		return
	}
	fmt.Printf("%s\n", b)
	defer r.Body.Close()
}
