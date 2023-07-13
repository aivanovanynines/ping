package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/aivanovanynines/ping/internal/models"
	"github.com/aivanovanynines/ping/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server         http.Server
	fetcherUsecase models.FetcherUsecase
}

func New(port int) *App {
	app := &App{
		server: http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		fetcherUsecase: usecase.NewFetcherUsecase(),
	}

	r := chi.NewRouter()
	r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("pong"))
		if err != nil {
			log.Printf("error response ping handler: %v", err)
		}
	})
	r.Post("/fetchApi", app.fetcherHandler)

	app.server.Handler = r

	return app
}

func (a *App) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		err := a.server.Shutdown(ctx)
		if err != nil {
			log.Printf("shutdown server error: %v", err)
		}
	}()

	log.Printf("starting server at address: %s", a.server.Addr)
	err := a.server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (a *App) fetcherHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := fetcherRequest{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid request payload: %s", err), http.StatusBadRequest)
		return
	}

	rules, err := requestData.getRules()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid request payload.rules: %s", err), http.StatusBadRequest)
		return
	}

	targetApi, err := url.Parse(requestData.TargetAPI)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid request.target_api: %s", err), http.StatusBadRequest)
		return
	}

	responseData, err := a.fetcherUsecase.Fetch(request.Context(), targetApi, rules)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseData, err := json.Marshal(&responseData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonResponseData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
