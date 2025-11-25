package addblock

import (
	"net/http"

	"go-runtime-demo/internal/app/blockchain/usecase/addblock"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/blocks"

type Handler struct {
	useCase addblock.UseCase
}

func NewHandler(useCase addblock.UseCase) Handler {
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

	block := h.useCase.Execute(r.Context(), payload.Data)
	httpjson.WriteJSON(w, http.StatusCreated, block)
}
