package gcprofile

import (
	"net/http"

	"go-runtime-demo/internal/app/monitoring/usecase/gcprofile"
	httpjson "go-runtime-demo/pkg/http"

	"github.com/gorilla/mux"
)

const Path = "/gc/profile"

type Handler struct {
	useCase gcprofile.UseCase
}

func NewHandler(useCase gcprofile.UseCase) Handler {
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
	if payload.DurationSeconds <= 0 {
		payload.DurationSeconds = 5
	}
	if payload.ProfileType == "" {
		payload.ProfileType = "heap"
	}

	// Parse profile type
	profileType := gcprofile.ProfileTypeHeap
	switch payload.ProfileType {
	case "cpu":
		profileType = gcprofile.ProfileTypeCPU
	case "goroutine":
		profileType = gcprofile.ProfileTypeGoroutine
	case "allocs":
		profileType = gcprofile.ProfileTypeAllocs
	}

	input := gcprofile.Input{
		DurationSeconds: payload.DurationSeconds,
		ProfileType:     profileType,
	}

	result := h.useCase.Execute(r.Context(), input)
	httpjson.WriteJSON(w, http.StatusOK, result)
}
