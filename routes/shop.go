package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	conf "github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/middlewares"
	"github.com/hackclub/hackatime/models/view"
	routeutils "github.com/hackclub/hackatime/routes/utils"
	"github.com/hackclub/hackatime/services"
	"github.com/hackclub/hackatime/utils"
)

type ShopHandler struct {
	config      *conf.Config
	userService services.IUserService
	shopService services.IShopService
}

func NewShopHandler(userService services.IUserService, shopService services.IShopService) *ShopHandler {
	return &ShopHandler{
		config:      conf.Get(),
		userService: userService,
		shopService: shopService,
	}
}

func (h *ShopHandler) RegisterRoutes(router chi.Router) {
	r := chi.NewRouter()
	r.Use(
		middlewares.NewAuthenticateMiddleware(h.userService).
			WithRedirectTarget(defaultErrorRedirectTarget()).
			WithRedirectErrorMessage("unauthorized").Handler,
		h.shopAvailabilityMiddleware,
	)
	r.Get("/", h.GetShop)

	router.Mount("/shop", r)
}

func (h *ShopHandler) shopAvailabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.config.Shop.Enabled {
			user := middlewares.GetPrincipal(r)
			if user != nil {
				http.Redirect(w, r, "/summary", http.StatusSeeOther)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
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

	products, err := h.shopService.GetProducts()
	if err != nil {
		conf.Log().Request(r).Error("failed to get products", "error", err.Error())
		return h.buildViewModel(r, w).WithError("failed to get products")
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
