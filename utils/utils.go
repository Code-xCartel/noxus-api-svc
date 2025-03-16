package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var Validate = validator.New()

func ParseJson(r *http.Request, payload any) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	if err == io.EOF {
		return fmt.Errorf("request body is empty")
	}
	return err
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	_ = WriteJson(w, status, map[string]string{"error": err.Error()})
}

func GenericJSON(message string) map[string]interface{} {
	return map[string]interface{}{"detail": message}
}

func ParseAndValidate(w http.ResponseWriter, r *http.Request, payload any) error {
	// get json
	if err := ParseJson(r, &payload); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return err
	}
	// validate the request
	if err := Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return err
	}
	return nil
}

var excludedURLS = []string{
	"/auth/login",
	"/auth/register",
}

func Contains[T comparable](arr []T, match T) bool {
	for _, v := range arr {
		if match == v {
			return true
		}
	}
	return false
}

func IsExcludedFromAuth(url string) bool {
	return Contains(excludedURLS, url)
}

func FilterArray[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
