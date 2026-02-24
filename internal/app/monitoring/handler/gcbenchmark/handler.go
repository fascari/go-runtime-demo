package gcbenchmark

import (
	"net/http"

	"go-runtime-demo/internal/app/monitoring/usecase/gcbenchmark"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/gc/benchmark"

type Handler struct {
	useCase gcbenchmark.UseCase
}

func NewHandler(useCase gcbenchmark.UseCase) Handler {
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

	// Set defaults
	if payload.Allocations <= 0 {
		payload.Allocations = 10000
	}
	if payload.SizeKB <= 0 {
		payload.SizeKB = 1
	}

	// Parse pattern
	pattern := gcbenchmark.PatternShortLived
	switch payload.Pattern {
	case "long-lived":
		pattern = gcbenchmark.PatternLongLived
	case "mixed":
		pattern = gcbenchmark.PatternMixed
	}

	input := gcbenchmark.Input{
		Allocations: payload.Allocations,
		SizeKB:      payload.SizeKB,
		Pattern:     pattern,
	}

	result := h.useCase.Execute(r.Context(), input)
	httpjson.WriteJSON(w, http.StatusOK, result)
}
