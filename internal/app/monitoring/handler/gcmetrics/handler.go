package gcmetrics

import (
	"net/http"

	"go-runtime-demo/internal/app/monitoring/usecase/gcmetrics"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/gc/metrics"

type Handler struct {
	useCase gcmetrics.UseCase
}

func NewHandler(useCase gcmetrics.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodGet)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	result := h.useCase.Execute(r.Context())
	httpjson.WriteJSON(w, http.StatusOK, result)
}
