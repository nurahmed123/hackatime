package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hackclub/hackatime/helpers"

	conf "github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/middlewares"
	"github.com/hackclub/hackatime/services"
)

type SpecialApiHandler struct {
	config   *conf.Config
	userSrvc services.IUserService
}

func NewSpecialApiHandler(userService services.IUserService) *SpecialApiHandler {
	return &SpecialApiHandler{
		userSrvc: userService,
		config:   conf.Get(),
	}
}

func (h *SpecialApiHandler) RegisterRoutes(router chi.Router) {
	r := chi.NewRouter()
	r.Use(middlewares.NewAuthenticateMiddleware(h.userSrvc).Handler)
	r.Get("/", h.Get)

	router.Mount("/email", r)
}

// @Summary Retrieve a users email
// @ID get-email
// @Tags email
// @Produce json
// @Param user query string false "The user to filter by if using Bearer authentication and the admin token"
// @Security ApiKeyAuth
// @Success 200 {object} models.Email
// @Router /email [get]
func (h *SpecialApiHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, err := h.userSrvc.GetUserById(r.URL.Query().Get("user"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	helpers.RespondJSON(w, r, http.StatusOK, map[string]interface{}{
		"email": user.Email,
	})
}
