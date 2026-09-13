package routes

import (
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"net/http"
)

func handelRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), false, nil)
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.JSON(w, http.StatusBadRequest, "name or email or password is required", false, nil)
		return
	}

	exitUser, _ := models.GetUserByEmail(req.Email)
	if exitUser != nil {
		utils.JSON(w, http.StatusConflict, "This Email is in used!", false, nil)
		return
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), false, nil)
		return
	}

	user, err := models.CreateUser(req.Name, req.Email, hashPassword)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), false, nil)
		return
	}
	utils.JSON(w, http.StatusCreated, "User Register Successfully", true, user)
}
