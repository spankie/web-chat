package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/spankie/web-chat/chat"
	"github.com/spankie/web-chat/web"
)

const (
	// DEFAULTPORT holds the default port for the server
	DEFAULTPORT = "8080"
)

// Router struct would carry the httprouter instance, so its methods could be verwritten and replaced with methds with wraphandler
type Router struct {
	*httprouter.Router
}

// Get is an endpoint to only accept requests of method GET
func (r *Router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

// Post is an endpoint to only accept requests of method POST
func (r *Router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

// Put is an endpoint to only accept requests of method PUT
func (r *Router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

// Delete is an endpoint to only accept requests of method DELETE
func (r *Router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

// NewRouter is a wrapper that makes the httprouter struct a child of the router struct
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// basically modifies the request context to include parameters
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)

		// call the next handler. 'cause this handler is done :)
		h.ServeHTTP(w, r)
	}
}

// TODO: How does context work in requests
// Learn casting interface with .(interface{})

func main() {
	log.Println("Hello Web Chat")

	// Initialize the http router ...
	router := NewRouter()

	commonHandlers := alice.New(web.LoggingHandler)

	// log.Println("commonhandlers: ", commonHandlers)

	// start chat server
	go chat.StartServer()

	router.Post("/api/search/friend", commonHandlers.Append(web.AuthHandler).ThenFunc(web.SearchFriend))
	router.Post("/api/add/friend/:id", commonHandlers.Append(web.AuthHandler).ThenFunc(web.AddFriend))
	router.Post("/api/get/friends", commonHandlers.Append(web.AuthHandler).ThenFunc(web.GetFriends))

	router.Post("/api/signup", commonHandlers.ThenFunc(web.Signup))
	router.Post("/api/login", commonHandlers.ThenFunc(web.Login))
	router.Post("/api/logout", commonHandlers.ThenFunc(web.Logout))

	router.Get("/", commonHandlers.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		// SERVE ANGULAR PAGE...
		// http.ServeFile(w, r, "web/templates/index.html")

		// SERVE MITHRIL PAGE...
		http.ServeFile(w, r, "mithril_client/index.html")
	}))

	router.Get("/api/chat", commonHandlers.Append(web.AuthHandler).ThenFunc(chat.Chat))

	// SERVE STATIC FILES FROM THE STATIC FOLDER
	fileServer := http.FileServer(http.Dir("./web/static"))
	router.GET("/assets/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})

	// SERVE MITHRIL FILES FROM THE MITHRIL FOLDER
	mithrilFileServer := http.FileServer(http.Dir("./mithril_client"))
	router.GET("/mithril/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		mithrilFileServer.ServeHTTP(w, r)
	})

	router.Get("/web/:user", commonHandlers.Append(web.AuthHandler).ThenFunc(web.UserArea))

	// Get the port to serve on from environment variable
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No Global port has been defined, using default 8080")
		PORT = "8080"
	}

	// This handler overrides the default servemux
	handler := cors.New(cors.Options{
		//		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)

	log.Println("serving on port :", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
