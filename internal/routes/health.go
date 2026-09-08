package routes

import (
	"backend/internal/utils"
	"net/http"
)

func handelHealth(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, "Api is Running Successfully!!", true, nil)
}
