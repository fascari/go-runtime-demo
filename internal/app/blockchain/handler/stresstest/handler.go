package stresstest

import (
	"net/http"

	"go-runtime-demo/internal/app/blockchain/usecase/stresstest"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/stress"

type Handler struct {
	useCase stresstest.UseCase
}

type InputPayload struct {
	Allocations int `json:"allocations"`
	Goroutines  int `json:"goroutines"`
}

func NewHandler(useCase stresstest.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodPost)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := httpjson.ReadJSON(r, &payload); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Allocations <= 0 {
		payload.Allocations = 10
	}

	if payload.Goroutines <= 0 {
		payload.Goroutines = 1
	}

	result := h.useCase.Execute(r.Context(), payload.Allocations, payload.Goroutines)
	httpjson.WriteJSON(w, http.StatusOK, result)
}
