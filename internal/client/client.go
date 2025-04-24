package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"gitlab.ozon.dev/timofey15g/homework/internal/handlers"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	desc "gitlab.ozon.dev/timofey15g/homework/pkg/service"
)

type App struct {
	client desc.OrderServiceClient
}

func NewApp(client desc.OrderServiceClient) *App {
	return &App{client: client}
}

type Handler interface {
	Execute(w http.ResponseWriter, r *http.Request)
}

func (app *App) Run() {
	hs := map[string]Handler{
		"accept":       handlers.NewAcceptOrder(app.client),
		"return":       handlers.NewReturnOrder(app.client),
		"issue":        handlers.NewIssueOrder(app.client),
		"withdraw":     handlers.NewWithdrawOrder(app.client),
		"list_order":   handlers.NewListOrder(app.client),
		"list_return":  handlers.NewListReturn(app.client),
		"list_history": handlers.NewListHistory(app.client),
		"renew_task":   handlers.NewRenewTask(app.client),
	}
	router := mux.NewRouter()

	router.HandleFunc("/orders/create", hs["accept"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/return", hs["return"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/issue", hs["issue"].Execute).Methods(http.MethodPost)
	router.HandleFunc("/orders/withdraw", hs["withdraw"].Execute).Methods(http.MethodDelete)
	router.HandleFunc("/orders/user", hs["list_order"].Execute).Methods(http.MethodGet)
	router.HandleFunc("/orders/returns", hs["list_return"].Execute).Methods(http.MethodGet)
	router.HandleFunc("/orders", hs["list_history"].Execute).Methods(http.MethodGet)

	router.HandleFunc("/tasks/reset", hs["renew_task"].Execute).Methods(http.MethodPost)

	router.Use(authMiddleware)

	log.Println("Client server is running on port 9000...")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()

		serverUser := os.Getenv("CLIENT_SERVER_USER")
		serverPassword := os.Getenv("CLIENT_SERVER_PASSWORD")

		if !ok || user != serverUser || password != serverPassword {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("invalid username or password")
			return
		}

		next := logMiddleware(handler)

		next.ServeHTTP(w, r)
	})
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	Body       *bytes.Buffer
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	rw.Body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logPipeline := logpipeline.GetLogPipelineInstance()

		var bodyBuffer bytes.Buffer

		tee := io.TeeReader(r.Body, &bodyBuffer)

		bodyBytes, err := io.ReadAll(tee)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}

		bodyString := string(bodyBytes)

		r.Body = io.NopCloser(&bodyBuffer)

		logPipeline.LogRequest(time.Now(), r.Method, r.URL.Path, bodyString)

		wrappedWriter := &ResponseWriterWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			Body:           &bytes.Buffer{},
		}

		next.ServeHTTP(wrappedWriter, r)

		logPipeline.LogResponse(time.Now(), int64(wrappedWriter.StatusCode), wrappedWriter.Body.String())
	})
}
