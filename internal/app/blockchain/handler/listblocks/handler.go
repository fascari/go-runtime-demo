package listblocks

import (
	"net/http"

	"go-runtime-demo/internal/app/blockchain/usecase/listblocks"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/blocks"

type Handler struct {
	useCase listblocks.UseCase
}

func NewHandler(useCase listblocks.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodGet)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	blocks := h.useCase.Execute(r.Context())
	httpjson.WriteJSON(w, http.StatusOK, blocks)
}
