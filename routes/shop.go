package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	conf "github.com/kcoderhtml/hackatime/config"
	"github.com/kcoderhtml/hackatime/middlewares"
	"github.com/kcoderhtml/hackatime/models"
	"github.com/kcoderhtml/hackatime/models/view"
	routeutils "github.com/kcoderhtml/hackatime/routes/utils"
	"github.com/kcoderhtml/hackatime/services"
	"github.com/kcoderhtml/hackatime/utils"
)

type ShopHandler struct {
	config      *conf.Config
	userService services.IUserService
}

func NewShopHandler(userService services.IUserService) *ShopHandler {
	return &ShopHandler{
		config:      conf.Get(),
		userService: userService,
	}
}

func (h *ShopHandler) RegisterRoutes(router chi.Router) {
	r := chi.NewRouter()
	r.Use(
		middlewares.NewAuthenticateMiddleware(h.userService).
			WithRedirectTarget(defaultErrorRedirectTarget()).
			WithRedirectErrorMessage("unauthorized").Handler,
	)
	r.Get("/", h.GetShop)

	router.Mount("/shop", r)
}

func (h *ShopHandler) GetShop(w http.ResponseWriter, r *http.Request) {
	if h.config.IsDev() {
		loadTemplates()
	}

	if err := templates[conf.ShopTemplate].Execute(w, h.buildViewModel(r, w)); err != nil {
		conf.Log().Request(r).Error("failed to get shop page", "error", err)
	}
}

func (h *ShopHandler) buildViewModel(r *http.Request, w http.ResponseWriter) *view.ShopViewModel {
	user := middlewares.GetPrincipal(r)
	if user == nil { // this should actually never occur, because of auth middleware
		w.WriteHeader(http.StatusUnauthorized)
		return h.buildViewModel(r, w).WithError("unauthorized")
	}

	products := []*models.Product{
		{
			Name:        "Sticker Pile",
			Price:       1,
			Description: "We'll send you 3 random stickers! (Available anywhere!)",
			Image:       "https://cloud-c1gqq7ttf-hack-club-bot.vercel.app/0sticker_pile_2.png",
		},
	}

	pageParams := utils.ParsePageParamsWithDefault(r, 1, 24)

	vm := &view.ShopViewModel{
		SharedLoggedInViewModel: view.SharedLoggedInViewModel{
			SharedViewModel: view.NewSharedViewModel(h.config, nil),
			User:            user,
			ApiKey:          user.ApiKey,
		},
		Products:   products,
		PageParams: pageParams,
	}
	return routeutils.WithSessionMessages(vm, r, w)
}
