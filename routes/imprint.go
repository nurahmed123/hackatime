package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	conf "github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
	"github.com/hackclub/hackatime/models/view"
	"github.com/hackclub/hackatime/services"
)

type ImprintHandler struct {
	config       *conf.Config
	keyValueSrvc services.IKeyValueService
}

func NewImprintHandler(keyValueService services.IKeyValueService) *ImprintHandler {
	return &ImprintHandler{
		config:       conf.Get(),
		keyValueSrvc: keyValueService,
	}
}

func (h *ImprintHandler) RegisterRoutes(router chi.Router) {
	router.Get("/imprint", h.GetImprint)
}

func (h *ImprintHandler) GetImprint(w http.ResponseWriter, r *http.Request) {
	if h.config.IsDev() {
		loadTemplates()
	}

	text := "failed to load content"
	if data, err := h.keyValueSrvc.GetString(models.ImprintKey); err == nil {
		text = data.Value
	}

	templates[conf.ImprintTemplate].Execute(w, h.buildViewModel(r).WithHtmlText(text))
}

func (h *ImprintHandler) buildViewModel(r *http.Request) *view.ImprintViewModel {
	return &view.ImprintViewModel{
		SharedViewModel: view.NewSharedViewModel(h.config, nil),
	}
}
