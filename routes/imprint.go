package routes

import (
	"net/http"
"github.com/kcoderhtml/hackatime/
	"github.com/kcoderhtml/hackatime/"
	"github.com/kcoderhtml/hackatime/ml/hackatime/config"
	"github.com/kcoderhtml/hackatime/ckatime/models"
	"github.com/kcoderhtml/hackatime/models/view"
	"github.com/kcoderhtml/hackatime/services"
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
