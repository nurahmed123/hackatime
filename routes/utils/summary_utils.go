package utils

import (
	"net/http"
	"strings"

	"github.com/hackclub/hackatime/helpers"
	"github.com/hackclub/hackatime/models"
	"github.com/hackclub/hackatime/models/types"
	"github.com/hackclub/hackatime/services"
)

func LoadUserSummary(ss services.ISummaryService, r *http.Request) (*models.Summary, error, int) {
	summaryParams, err := helpers.ParseSummaryParams(r)
	if err != nil {
		return nil, err, http.StatusBadRequest
	}
	return LoadUserSummaryByParams(ss, summaryParams)
}

func LoadUserSummaryByParams(ss services.ISummaryService, params *models.SummaryParams) (*models.Summary, error, int) {
	var retrieveSummary types.SummaryRetriever = ss.Retrieve
	if params.Recompute {
		retrieveSummary = ss.Summarize
	}

	summary, err := ss.Aliased(
		params.From,
		params.To,
		params.User,
		retrieveSummary,
		params.Filters,
		params.Recompute,
	)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	summary.FromTime = models.CustomTime(summary.FromTime.T().In(params.User.TZ()))
	summary.ToTime = models.CustomTime(summary.ToTime.T().In(params.User.TZ()))

	return summary, nil, http.StatusOK
}

func FilterColors(all map[string]string, haystack models.SummaryItems) map[string]string {
	subset := make(map[string]string)
	for _, item := range haystack {
		if c, ok := all[strings.ToLower(item.Key)]; ok {
			subset[strings.ToLower(item.Key)] = c
		}
	}
	return subset
}
