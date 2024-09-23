package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hackclub/hackatime/helpers"

	conf "github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/middlewares"
	v1 "github.com/hackclub/hackatime/models/compat/wakatime/v1"
	routeutils "github.com/hackclub/hackatime/routes/utils"
	"github.com/hackclub/hackatime/services"
)

type UsersHandler struct {
	config        *conf.Config
	userSrvc      services.IUserService
	heartbeatSrvc services.IHeartbeatService
}

func NewUsersHandler(userService services.IUserService, heartbeatService services.IHeartbeatService) *UsersHandler {
	return &UsersHandler{
		userSrvc:      userService,
		heartbeatSrvc: heartbeatService,
		config:        conf.Get(),
	}
}

func (h *UsersHandler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middlewares.NewAuthenticateMiddleware(h.userSrvc).Handler)
		r.Get("/compat/wakatime/v1/users/{user}", h.Get)
	})
}

// @Summary Retrieve the given user
// @Description Mimics https://wakatime.com/developers#users
// @ID get-wakatime-user
// @Tags wakatime
// @Produce json
// @Param user path string true "User ID to fetch (or 'current')"
// @Security ApiKeyAuth
// @Success 200 {object} v1.UserViewModel
// @Router /compat/wakatime/v1/users/{user} [get]
func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	wakapiUser, err := routeutils.CheckEffectiveUser(w, r, h.userSrvc, "current")
	if err != nil {
		return // response was already sent by util function
	}

	user := v1.NewFromUser(wakapiUser)
	if hb, err := h.heartbeatSrvc.GetLatestByUser(wakapiUser); err == nil {
		user = user.WithLatestHeartbeat(hb)
	} else {
		conf.Log().Request(r).Error("error occurred", "error", err)
	}

	helpers.RespondJSON(w, r, http.StatusOK, v1.UserViewModel{Data: user})
}
