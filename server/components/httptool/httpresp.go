package httptool

import (
	"net/http"
)




func OKResponse(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func SeeOtherResponse(w http.ResponseWriter, leaderAddr string) {
	w.Header().Add("leader-id", leaderAddr)
	w.WriteHeader(http.StatusSeeOther)
}

func BadResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func ForbiddenResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func ErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Add("message", message)
	w.WriteHeader(http.StatusInternalServerError)
}


