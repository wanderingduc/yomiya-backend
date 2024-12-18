package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"yomiya/backend/api/auth"
	"yomiya/backend/api/responses"
)

func CheckToken(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response responses.JSONResponse
		h := r.Header
		tokenString := h["Authorization"][0][7:]
		client := r.RemoteAddr
		log.Printf("Authenticating [%s]...", client)

		err := auth.CheckToken(tokenString)
		if err != nil {
			log.Printf("[%s] failed to authenticate", client)
			errResponse := responses.JSONError{
				Err: "Invalid token",
			}
			response = responses.JSONResponse{
				Success: false,
				Data:    errResponse,
				Meta:    nil,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Printf("Authenticated [%s]", client)
		next.ServeHTTP(w, r)
	})

}
