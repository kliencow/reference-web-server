# Golang / Gorilla Reference Server

## Purpose
I wrote this thing because I was getting wound around the axel trying to figure out all the little subtle parts to golang webservers while, and at the same time, trying to learn Gorilla session management and Gorilla router / handler chaining. To simplify things, I just wrote a quick few-trick-poney to explore the space. 

## Implementation Details
The implementation here has some good parts and bad parts. I like that it's pretty simple. Learning about subrouters was great! However, relying on subrouters is part of my greif here. I would like to spend some more time to figure out how to do this without using subrouters to protect some of the site, but not others. In some respects, it is elegant. The client gets to know by url which pages are secure and which are not. 

## The learned bits:
### Router Chaining
This turned out to be pretty simple. There are two part's of Gorilla's mux system. The `subRouter` and the `Use` functions. I wanted to be able to do two common things: log requests and secure my secret page. 

In the code below I create a main router that handles my login and forbidden pages since these need to be open to the world. The subrouter is created from the main router and is attached to it. The reference to it is then used to decorate the router. Since the subrouter is part of the main router, the `logHander` will get called for both while the `authHandler` will only be called for the stuff in the `closedHander`. 
```go
mainRouter := mux.NewRouter()
mainRouter.HandleFunc("/login", login)
mainRouter.HandleFunc("/forbidden", forbidden)
mainRouter.Use(logHander)

closedRouter := mainRouter.PathPrefix("/auth").Subrouter()
closedRouter.HandleFunc("/secret", secret)
closedRouter.Use(authHandler)
```
After this, just attach the mainRouter to you server and you're good to go!

### Router Chaining Complications
It turns out that routing isn't quite so simple. The trick here is the `http.ResponseWriter` which is devilishly simple and easy to get into trouble with. Here's the sticky bit: writing to the header and the writer in the wrong order will cause the redirection to stop working and start barfing ugly `http: superfluous response.WriteHeader call` errors to start printing. 

### Redirecting
There are two flavors of redirection that I tried: `http.RedirectHandler` and `http.Redirect`. The former would be used for chaining, I suspect. The latter is what I ended up using for simplicity. It's importand to note that if you have already written to the `http.ResponseWriter` this will fail. The reason is that if you have not set the header by the time you write, go will set the status code to 200 for you when you do write. So, when you pass in the `http.StatusFound` 302 to redirect the client browser, it will fail. The 
```go
http.Redirect(w, r, "/forbidden", http.StatusFound)
```

### Sessions
This are silly simple, but there are a couple of tricks here as well. The good news is that the gotchas are related to the writing parts out at the wrong times.

Here are the three parts to authentication:

Logging in:
```go
session, _ := store.Get(r, cookieName)

session = "walts-rad-login"

session.Save(r, w)
```

Using Session Information

```go
session, _ := store.Get(r, cookieName)

name := session.Values[cookieValKey]

```
Checking Session Information

```go
session, _ := store.Get(r, cookieName)

if _, ok := session.Values[cookieValKey]; !ok {
  http.Redirect(w, r, "/forbidden", http.StatusFound)
  return
}
```
Logging out will be left as an excersize for the reader. 

### Session Tricks
The trick with sessions comes from the bit `session.Save(r, w)`. This part has to come at the right point or the session won't save. You need to save the session before you start write out to the `http.ResponseWriter`. That's the trick! Done and Done.
