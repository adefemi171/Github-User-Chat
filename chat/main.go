package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// template represents a single template
// struct type responsible for loading, compiling and delivering the template content
type templateHandler struct {
	once     sync.Once //to compile template once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates/", t.filename)))
	})
	// t.templ.Execute(w, nil)
	// Passes request details as data into the Execute method
	t.templ.Execute(w, r)
}

// always specify using ./githubuserChat -addr="192.168.43.195:7000"
func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(`
	// 		<html>
	// 			<head>
	// 				<title> Github User Chat </title>
	// 			</head>
	// 			<body>
	// 				Live Testing
	// 			</body>
	// 		</html>
	// 	`))
	// })

	var addr = flag.String("addr", ":7000", "http service address of the application.")
	flag.Parse() // Parsing the flag
	r := newRoom()
	//To print trace debugging in terminal
	// r.tracer = trace.New(os.Stdout)
	// root
	http.Handle("/", &templateHandler{filename: "test.html"})
	http.Handle("/room", r)
	//get the room going
	go r.run()
	log.Println("Starting web server on ", addr)
	fmt.Println("Open browser and redirect to http://localhost:7000")
	//Startting webserver on port 7000
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
