package controller

import (
	"crud/model"
	"crud/model/request"
	"crud/model/response"
	"crud/repository"
	"crud/util"
	"crud/util/token"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	userDb, err := util.ExtractFromRequest[repository.UserRepository](r, "userDb")
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	data := new(request.Auth)

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(data)

	if err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid request body").Write(w)
		return

	}

	if _, err = userDb.GetByEmailOrUsername(r.Context(), data.Username, data.Email); err == nil {
		NewErrorResponse(http.StatusBadRequest, "User already exists").Write(w)
		return
	}

<<<<<<< HEAD
	res, err := userDb.Create(r.Context(), model.User{Username: data.Username, Email: data.Email, PasswordHash: data.Password})
=======
	hashed, err := HashPassword(data.Password)
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	res, err := userDb.Create(r.Context(), model.User{
		Username:     data.Username,
		Email:        data.Email,
		PasswordHash: hashed,
	})

>>>>>>> 8a7e0ed (refactor: move and change few logic)
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	NewResponse(http.StatusCreated, true, "User created", res.Id.Hex()).Write(w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	userDb, err := util.ExtractFromRequest[repository.UserRepository](r, "userDb")
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	data := new(request.Auth)

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(data)

	if err != nil {
		NewErrorResponse(http.StatusBadRequest, "Invalid request body").Write(w)
		return

	}

	user, err := userDb.GetByEmailOrUsername(r.Context(), data.Username, data.Email)
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
	tokenizer, err := util.ExtractFromRequest[token.Token](r, "token")
	if err != nil {
		NewErrorResponse(http.StatusInternalServerError, err.Error()).Write(w)
		return
	}

	enc := tokenizer.Encrypt(
		token.WithClaims("user_id", user.Id.Hex()),
		token.WithClaims("user_name", user.Username),
		token.WithClaims("user_email", user.Email),
	)

	NewResponse(http.StatusOK, true, "Login successful", response.Auth{AccessToken: enc}).Write(w)
}
