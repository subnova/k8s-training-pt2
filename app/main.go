package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var greetingName *string
var langFile *string
var langMap map[string]string

func hello(w http.ResponseWriter, req *http.Request) {
	lang := req.Header.Get("Accept-Language")
	formatString := langMap[lang]
	if formatString == "" {
		formatString = "Hello %s"
	}

	fmt.Fprintf(w, formatString+"\n", *greetingName)
}

func main() {
	defaultName := os.Getenv("GREETING_NAME")
	if defaultName == "" {
		defaultName = "world"
	}
	greetingName = flag.String("name", defaultName, "name to greet")
	langFile = flag.String("langs", "langs.json", "file containing language templates")
	flag.Parse()

	readLangMap()

	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}

func readLangMap() {
	file, err := os.Open(*langFile)
	if err != nil {
		panic("unable to open language file")
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(bytes), &langMap)
}
