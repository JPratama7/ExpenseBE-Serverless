package controller

import (
	"crud/repository"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	userDb, err := ExtractFromRequest[repository.UserRepository](r, "userDb")
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	type Register struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data := new(Register)

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(data)

	if err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid request body").Write(w)
		return

	}

	if _, err = userDb.GetByEmailOrUsername(data.Username, data.Email); err == nil {
		NewErrorResponse(http.StatusBadRequest, "User already exists").Write(w)
		return
	}

	res, err := userDb.Create(data.Username, data.Email, data.Password)
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	NewResponse(http.StatusCreated, true, "User created", res.Id.Hex()).Write(w)

}

func Login(w http.ResponseWriter, r *http.Request) {
	userDb, err := ExtractFromRequest[repository.UserRepository](r, "userDb")
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	type Login struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data := new(Login)

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(data)

	if err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid request body").Write(w)
		return

	}

	user, err := userDb.GetByEmailOrUsername(data.Username, data.Email)
	if err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid email or password").Write(w)
		return
	}

	passHashed, err := HashPassword(data.Password)
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	if err = ComparePassword(passHashed, user.PasswordHash); err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid email or password").Write(w)
		return
	}

	// TODO: Generate Access and Refresh Token

	NewResponse(http.StatusOK, true, "Login successful", user.Id.Hex()).Write(w)
}
