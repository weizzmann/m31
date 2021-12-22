package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"m31/internal/apperror"
	"m31/internal/handlers"
	"m31/pkg/logging"
)

const (
	usersURL       = "/users"
	userURL        = "/users/{id}"
	makeFriends    = "/users/make_friends"
	userFriendsURL = "/users/{id}/friends"
)

type Handler struct {
	Service *Service
	Logger  *logging.Logger
}

func NewHandler(logger *logging.Logger, service *Service) handlers.Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Get(usersURL, apperror.Middleware(h.GetUsers))
	router.Post(usersURL, apperror.Middleware(h.CreateUser))
	router.Get(userURL, apperror.Middleware(h.GetUserByID))
	router.Put(userURL, apperror.Middleware(h.UpdateUser))
	router.Delete(userURL, apperror.Middleware(h.DeleteUser))
	router.Post(makeFriends, apperror.Middleware(h.MakeFriends))
	router.Get(userFriendsURL, apperror.Middleware(h.GetUserFriends))
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.Service.FindAll(context.Background())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(users); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var createUserDTO *CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&createUserDTO); err != nil {
		return apperror.SystemError(err)
	}

	id, err := h.Service.Create(context.Background(), createUserDTO)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(fmt.Sprintf("id: %v", id)); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user, err := h.Service.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var userAge *UpdateUserAgeDTO
	if err := json.NewDecoder(r.Body).Decode(&userAge); err != nil {
		return apperror.SystemError(err)
	}
	user, err := h.Service.UpdateUser(context.Background(), id, userAge)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user, err := h.Service.DeleteUser(context.Background(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}

func (h *Handler) MakeFriends(w http.ResponseWriter, r *http.Request) error {
	var CreateFriendDTO *CreateFriendDTO

	if err := json.NewDecoder(r.Body).Decode(&CreateFriendDTO); err != nil {
		return apperror.SystemError(err)
	}

	message, err := h.Service.MakeFriends(context.Background(), CreateFriendDTO)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(fmt.Sprintf("message: %v", message)); err != nil {
		return apperror.SystemError(err)
	}

	return nil

}

func (h *Handler) GetUserFriends(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	users, err := h.Service.GetUserFriends(context.Background(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(users); err != nil {
		return apperror.SystemError(err)
	}

	return nil
}
