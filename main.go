package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("totes-secret-key")
	store = sessions.NewCookieStore(key)
)

const (
	cookieName   = "Mmmmm!Cookies!!"
	cookieValKey = "some-val"
	timeout      = 10 * time.Second
	host         = "localhost:8000"
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cookieName)

	name := session.Values[cookieValKey]

	fmt.Fprintf(w, "The cookie's secret is: %v\n", name)
	fmt.Fprintln(w, "Hi!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cookieName)

	session.Values[cookieValKey] = "walts-rad-login"
	session.Save(r, w)

	fmt.Fprint(w, "You're Logged in!")
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "FORBIDDEN!")
}

func logHander(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, cookieName)

		if v, ok := session.Values[cookieValKey]; !ok {
			fmt.Printf("%v\n", v)
			fmt.Printf("%v\n", ok)
			http.Redirect(w, r, "/forbidden", http.StatusFound)
			return
		} else {
			fmt.Printf("%v\n", v)
			fmt.Printf("%v\n", ok)
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/login", login)
	mainRouter.HandleFunc("/forbidden", forbidden)
	mainRouter.Use(logHander)

	closedRouter := mainRouter.PathPrefix("/auth").Subrouter()
	closedRouter.HandleFunc("/secret", secret)
	closedRouter.Use(authHandler)

	srv := &http.Server{
		Handler:      mainRouter,
		Addr:         host,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
