package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("userId")
	if userIDStr == "" {
		http.Error(w, "userId header missing", http.StatusBadRequest)
		return
	}
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	accessToken, _ := jwtutil.GenerateToken(int64(userID), 365*24*time.Hour)

	respMap := map[string]any{
		"accessToken":  accessToken,
		"refreshToken": "dummy-refresh-token",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respMap)
}
