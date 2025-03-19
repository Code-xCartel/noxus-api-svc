package friends

import (
	"fmt"
	"github.com/Code-xCartel/noxus-api-svc/types/auth"
	"github.com/Code-xCartel/noxus-api-svc/utils"
	"net/http"
)

func Router(router *http.ServeMux, store *Store) {

	router.HandleFunc("GET /friends/search/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			user, err := store.SearchNoxId(noxId)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			_ = utils.WriteJson(w, http.StatusOK, user)
		})

	router.HandleFunc("GET /friends",
		func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			friends, err := store.GetFriends(claims.NoxID, Accepted)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			_ = utils.WriteJson(w, http.StatusOK, friends)
		})

	router.HandleFunc("GET /friends/pending",
		func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			pending, err := store.GetFriends(claims.NoxID, Pending)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			_ = utils.WriteJson(w, http.StatusOK, pending)
		})

	router.HandleFunc("POST /friends/add/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			fErr := store.AddFriendByNoxId(claims.NoxID, u.NoxID)
			if fErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, fErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusCreated, utils.GenericJSON("request sent"))
		})

	router.HandleFunc("DELETE /friends/remove/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			delErr := store.DeleteFriendByNoxId(claims.NoxID, u.NoxID)
			if delErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, delErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusOK, utils.GenericJSON("friend removed"))
		})

	router.HandleFunc("PUT /friends/accept/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			fErr := store.ActionOnFriendByNoxId(claims.NoxID, u.NoxID, []Status{Pending}, Accepted)
			if fErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, fErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusCreated, utils.GenericJSON("request accepted"))
		})

	router.HandleFunc("PUT /friends/reject/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			actionErr := store.ActionOnFriendByNoxId(claims.NoxID, u.NoxID, []Status{Pending}, Rejected)
			if actionErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, actionErr)
				return
			}
			delErr := store.DeleteFriendByNoxId(claims.NoxID, u.NoxID)
			if delErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, delErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusCreated, utils.GenericJSON("request rejected"))
		})

	router.HandleFunc("PUT /friends/block/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			actionErr := store.ActionOnFriendByNoxId(claims.NoxID, u.NoxID, []Status{Accepted, Pending}, Blocked)
			if actionErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, actionErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusCreated, utils.GenericJSON("friend blocked"))
		})

	router.HandleFunc("PUT /friends/unblock/{noxId}",
		func(w http.ResponseWriter, r *http.Request) {
			noxId := r.PathValue("noxId")
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			u, err := store.authStore.GetUserByNoxID(noxId)
			if err != nil || noxId == claims.NoxID {
				utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid user"))
				return
			}
			actionErr := store.ActionOnFriendByNoxId(claims.NoxID, u.NoxID, []Status{Blocked}, Accepted)
			if actionErr != nil {
				utils.WriteError(w, http.StatusInternalServerError, actionErr)
				return
			}
			_ = utils.WriteJson(w, http.StatusCreated, utils.GenericJSON("friend unblocked"))
		})

	router.HandleFunc("GET /friends/blocked",
		func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value("claims").(*auth.CustomClaims)
			blocked, err := store.GetFriends(claims.NoxID, Blocked)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			_ = utils.WriteJson(w, http.StatusOK, blocked)
		})

}
