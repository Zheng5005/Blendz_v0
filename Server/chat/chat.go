package chat

import (
	"encoding/json"
	"net/http"

	"github.com/Zheng5005/Blendz_v0/stream"
	"github.com/Zheng5005/Blendz_v0/utils"
)

func GetStreamToken(w http.ResponseWriter, r *http.Request)  {
	userId, err := utils.ParseToken(r)
	if err != nil {
		http.Error(w, "No cookie provied", http.StatusUnauthorized)
		return
	}

	token, err := stream.GenerateStreamToken(userId)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
