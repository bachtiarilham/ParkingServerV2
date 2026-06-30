package response

import (
	"encoding/json"
	"net/http"
)

func JSON(
	w http.ResponseWriter,
	status int,
	resp APIResponse,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}

func Success(
	w http.ResponseWriter,
	status int,
	message string,
	data any,
) {
	JSON(
		w,
		status,
		APIResponse{
			Success: true,
			Message: message,
			Data:    data,
		},
	)
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
) {
	JSON(
		w,
		status,
		APIResponse{
			Success: false,
			Message: message,
		},
	)
}
