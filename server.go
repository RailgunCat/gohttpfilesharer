package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"time"
	"strconv"
)

var filesPrefix = "/file/"

var workingDirectory string = "./"
var port int = 3000

func main() {
	http.HandleFunc(filesPrefix,  logHandler(func(w http.ResponseWriter, r *http.Request) {
		log.Println("/GET " + r.URL.Path)
		http.ServeFile(w, r, r.URL.Path[len(filesPrefix):])
	}))

	http.HandleFunc("/", logHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(listFiles()))
	}))

	log.Printf("try start listening at port %d with work directory %s", port, workingDirectory)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Println(err)
	}
}

func listFiles() string {
	template := "<html><body><p>Files list</p><ul>%s<ul></body></html>"

	files, _ := ioutil.ReadDir(workingDirectory)
	res := ""
	for _, f := range files {
		if (!f.IsDir()) {
			res += fmt.Sprintf("<li><a href=\"%s\">%s\t%d bytes</a></li>", filesPrefix + f.Name(), f.Name(), f.Size())
		}
	}

	return fmt.Sprintf(template, res)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(">>%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		fn(w, r)
		log.Printf("<< served for %s", time.Since(start))
	}
}