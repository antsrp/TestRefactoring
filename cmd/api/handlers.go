package main

import (
	"net/http"
	"refactoring/internal/exchange/requests"
	"refactoring/internal/exchange/responses"
	"refactoring/internal/service"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.SugaredLogger
	service *service.Service
}

func createNewHandler(logger *zap.Logger, s *service.Service) *Handler {

	return &Handler{
		logger:  logger.Sugar(),
		service: s,
	}
}

func (h Handler) Routes() chi.Router {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", h.getUsers)
				r.Post("/", h.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.getUser)
					r.Patch("/", h.updateUser)
					r.Delete("/", h.deleteUser)
				})
			})
		})
	})

	return r
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(w, r, responses.ErrInvalidRequest(err))
}

func getID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

func (h Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	store := h.service.GetUsers()

	render.JSON(w, r, store.List)
}

func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}
	u, err := h.service.GetUser(id)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	render.JSON(w, r, *u)
}

func (h Handler) createUser(w http.ResponseWriter, r *http.Request) {

	request := requests.CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	id := h.service.AddUser(request)
	resp := responses.CreateUserResponse{UserID: id}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func (h Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	request := requests.UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	if err := h.service.UpdateUser(id, request); err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
