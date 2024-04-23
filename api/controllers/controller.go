package controllers

import (
	"encoding/json"
	// "log"
	"net/http"
	"rate-limiter/api/types"
	"rate-limiter/pkg/constants"
)

func FancyController(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(constants.CtxKey).(types.Response)
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		http.Error(w, "Error encoding response: ", http.StatusInternalServerError)
		return
	}
}
