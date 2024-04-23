package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"rate-limiter/api/types"
	"rate-limiter/pkg/constants"
)

func HandleError(w http.ResponseWriter, reason string, userID string, timeBeforeNextReq string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	log.Printf("user: %-10v | success: %-5v | req made in window: %-5v | time before next req: %-3v", userID, false, constants.MaxReqAllowedPerWindow, timeBeforeNextReq)

	json.NewEncoder(w).Encode(types.Response{Success: false, Reason: reason})
}
