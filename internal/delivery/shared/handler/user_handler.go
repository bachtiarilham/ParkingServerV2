package handler

import (
	"errors"
	dto "modulegue/internal/delivery/shared/dto"
	"modulegue/internal/middleware"   // Import middleware untuk userID
	"modulegue/internal/usecase/user" // Import usecase user
	"modulegue/pkg/response"
	"net/http"
)

type UserHandler struct {
	getCurrentProfileUC *user.GetProfileUseCase
}

func NewUserHandler(getCurrentProfileUC *user.GetProfileUseCase) *UserHandler {
	return &UserHandler{
		getCurrentProfileUC: getCurrentProfileUC,
	}
}

// Endpoint: GET /api/v2/linespot/users/me
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari JWT context
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Panggil usecase
	input := user.GetProfileInput{
		UserID: userID,
	}
	result, err := h.getCurrentProfileUC.Execute(r.Context(), input)
	if err != nil {
		// Log error jika perlu
		if errors.Is(err, user.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	// Mapping: UseCase Output -> DTO
	resp := dto.UserProfileDto{
		UserId:     result.UserID,
		Nik:        result.Nik,
		FullName:   result.FullName,
		Phone:      result.Phone,
		Email:      result.Email,
		Username:   result.Username,
		Role:       result.RoleID, // Gunakan RoleID dari result
		IsVerified: result.IsVerified,
	}

	response.Success(w, http.StatusOK, "User profile retrieved", resp)
}
