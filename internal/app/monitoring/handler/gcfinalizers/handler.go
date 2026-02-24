package gcfinalizers

import (
	"net/http"

	"go-runtime-demo/internal/app/monitoring/usecase/gcfinalizers"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/gc/finalizers"

type Handler struct {
	useCase gcfinalizers.UseCase
}

func NewHandler(useCase gcfinalizers.UseCase) Handler {
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

	if payload.Count <= 0 {
		payload.Count = 100
	}

	input := gcfinalizers.Input{
		Count:     payload.Count,
		TriggerGC: payload.TriggerGC,
	}

	result := h.useCase.Execute(r.Context(), input)
	httpjson.WriteJSON(w, http.StatusOK, result)
}
