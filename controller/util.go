package controller

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	hexHash := hex.EncodeToString(hash)

	return hexHash, nil
}

func ComparePassword(password, hashedPassword string) error {
	decodedHash, err := hex.DecodeString(hashedPassword)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(decodedHash, []byte(password))
}

func NewResponse[T any](code int, status bool, message string, data T) Response[T] {
	return Response[T]{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string) Response[*int] {
	return NewResponse[*int](code, false, message, nil)
}

type Response[T any] struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func (r Response[T]) Error() string {
	return fmt.Sprintf("code: %d, status: %t, message: %s", r.Code, r.Status, r.Message)
}

func (r Response[T]) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		panic(err)
	}
}
