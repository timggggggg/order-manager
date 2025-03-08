package service

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.ozon.dev/timofey15g/homework/internal/commands"
	storage "gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
)

type App struct {
	storage *storage.PgFacade
}

func NewApp(storage *storage.PgFacade) *App {
	return &App{storage}
}

type Command interface {
	Execute(w http.ResponseWriter, r *http.Request)
}

func (app *App) Run() {
	cmds := map[string]Command{
		"accept":       commands.NewAcceptOrder(app.storage),
		"return":       commands.NewReturnOrder(app.storage),
		"issue":        commands.NewIssueOrder(app.storage),
		"withdraw":     commands.NewWithdrawOrder(app.storage),
		"list_order":   commands.NewListOrder(app.storage),
		"list_return":  commands.NewListReturn(app.storage),
		"list_history": commands.NewListHistory(app.storage),
	}

	router := mux.NewRouter()

	router.HandleFunc("/orders/create", cmds["accept"].Execute)
	router.HandleFunc("/orders/return", cmds["return"].Execute)
	router.HandleFunc("/orders/issue", cmds["issue"].Execute)
	router.HandleFunc("/orders/withdraw", cmds["withdraw"].Execute)
	router.HandleFunc("/orders/user", cmds["list_order"].Execute)
	router.HandleFunc("/orders/returns", cmds["list_return"].Execute)
	router.HandleFunc("/orders", cmds["list_history"].Execute)

	router.Use(authMiddleware)

	log.Println("Server is running on port 9000...")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// fmt.Print(userName, password)
		if userName != "admin" || password != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("invalid username or password")
			return
		}

		next := logMiddleware(handler)

		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			log.Printf("Request: %s %s, Body: %v", r.Method, r.URL.Path, r.Body)
		}
		next.ServeHTTP(w, r)
	})
}
