package stats

import (
	"net/http"

	"go-runtime-demo/internal/app/monitoring/usecase/stats"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/stats"

type Handler struct {
	useCase stats.UseCase
}

func NewHandler(useCase stats.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r *mux.Router, h Handler) {
	r.HandleFunc(Path, h.Handle).Methods(http.MethodGet)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	runtimeStats := h.useCase.Execute(r.Context())
	httpjson.WriteJSON(w, http.StatusOK, runtimeStats)
}
