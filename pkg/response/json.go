package response

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func Render(w http.ResponseWriter, log *zap.Logger, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if resp == nil {
		return
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("encoding error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
