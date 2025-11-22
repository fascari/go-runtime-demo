package mineparallel

import (
	"net/http"

	"go-runtime-demo/internal/app/blockchain/usecase/mineparallel"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/mine"

type Handler struct {
	useCase mineparallel.UseCase
}

type InputPayload struct {
	Data       string `json:"data"`
	Goroutines int    `json:"goroutines"`
}

func NewHandler(useCase mineparallel.UseCase) Handler {
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

	if payload.Data == "" {
		httpjson.WriteError(w, http.StatusBadRequest, httpjson.ErrMissingValue)
		return
	}

	if payload.Goroutines <= 0 {
		payload.Goroutines = 1
	}

	result := h.useCase.Execute(r.Context(), payload.Data, payload.Goroutines)
	httpjson.WriteJSON(w, http.StatusOK, result)
}
