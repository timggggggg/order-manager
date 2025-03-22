package service

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type Storage interface {
	handlers.AcceptStorage
	handlers.ReturnStorage
	handlers.IssueStorage
	handlers.WithdrawStorage
	handlers.ListOrderStorage
	handlers.ListReturnStorage
	handlers.ListHistoryStorage
}

type ILogPipeline interface {
	LogStatusChange(TS time.Time, ID int64, statusFrom, statusTo models.OrderStatus)
	Shutdown()
}

type App struct {
	storage     Storage
	logPipeline ILogPipeline
}

func NewApp(storage Storage, logPipeline ILogPipeline) *App {
	return &App{storage, logPipeline}
}

type Handler interface {
	Execute(w http.ResponseWriter, r *http.Request)
}

func (app *App) Run() {
	hs := map[string]Handler{
		"accept":       handlers.NewAcceptOrder(app.storage, app.logPipeline),
		"return":       handlers.NewReturnOrder(app.storage, app.logPipeline),
		"issue":        handlers.NewIssueOrder(app.storage, app.logPipeline),
		"withdraw":     handlers.NewWithdrawOrder(app.storage, app.logPipeline),
		"list_order":   handlers.NewListOrder(app.storage, app.logPipeline),
		"list_return":  handlers.NewListReturn(app.storage, app.logPipeline),
		"list_history": handlers.NewListHistory(app.storage, app.logPipeline),
	}
	router := mux.NewRouter()

	router.HandleFunc("/orders/create", hs["accept"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/return", hs["return"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/issue", hs["issue"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/withdraw", hs["withdraw"].Execute).Methods(http.MethodDelete)
	router.HandleFunc("/orders/user", hs["list_order"].Execute).Methods(http.MethodGet)
	router.HandleFunc("/orders/returns", hs["list_return"].Execute).Methods(http.MethodGet)
	router.HandleFunc("/orders", hs["list_history"].Execute).Methods(http.MethodGet)

	router.Use(authMiddleware)

	log.Println("Server is running on port 9000...")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()

		serverUser := os.Getenv("SERVER_USER")
		serverPassword := os.Getenv("SERVER_PASSWORD")

		if !ok || user != serverUser || password != serverPassword {
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
