package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/adefemi171/githubChat/trace"
	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"
	// "github.com/stretchr/gomniauth/providers/outlook"
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
		// t.templ = template.Must(template.ParseFiles(filepath.Join("../templates/", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["userData"] = objx.MustFromBase64(authCookie.Value)
	}
	// t.templ.Execute(w, nil)
	// Passes request details as data into the Execute method
	// t.templ.Execute(w, r)
	t.templ.Execute(w, data)
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

	var host = flag.String("addr", ":7000", "http service address of the application.")
	flag.Parse() // Parsing the flag
	// Setting up gomniauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		github.New("12e9f8e2dc131f2f19ae", "48ec096a33e5cccb14bbccc605c68837183cff38", "http://localhost:7000/auth/callback/github"),
		google.New("312929473322-k77gg52v8ibpr7jan148gallnuslorjj.apps.googleusercontent.com", "nO3Y7zIMdgMzAdboAEx5dOD2", "http://localhost:7000/auth/callback/google"),
		// outlook.New("key", "secret", "http://localhost:7000/authy/callback/outlook"),
	)
	r := newRoom()
	//To print trace debugging in terminal
	r.tracer = trace.New(os.Stdout)
	// root
	// http.Handle("/", &templateHandler{filename: "test.html"})
	// // //
	// Wrapping templateHandler with the MustAuthy func will
	// allow the execution run through authyHandler firsty,
	// it will run only to templateHandler if the request is authenticated.
	http.Handle("/chat", MustAuthy(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//get the room going
	go r.run()
	log.Println("Starting web server on ", *host)
	fmt.Println("Open browser and redirect to http://localhost:7000")
	//Startting webserver on port 7000
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
