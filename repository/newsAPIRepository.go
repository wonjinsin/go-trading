package repository

import (
	"context"
	"fmt"
	"magmar/config"
	"magmar/model"
	"magmar/util"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
)

type newsAPIRepository struct {
	conn   *resty.Client
	apiURL util.APIURL
	apiKey string
}

// NewNewsAPIRepository ...
func NewNewsAPIRepository(magmar *config.ViperConfig) NewsRepository {
	return &newsAPIRepository{
		conn:   resty.New(),
		apiURL: util.NewsAPIURL,
		apiKey: magmar.GetString(util.NewsAPIKey),
	}
}

// GetNews ...
func (n *newsAPIRepository) GetNews(ctx context.Context, keywords []string) (newses model.Newses, err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	resp, err := n.conn.R().
		SetResult(&newses).
		SetQueryParam("q", strings.Join(keywords, " and ")).
		SetQueryParam("language", "en").
		SetQueryParam("sortBy", "publishedAt").
		SetQueryParam("pageSize", "10").
		SetQueryParam("apiKey", n.apiKey).
		Get(fmt.Sprintf("%s/v2/everything", n.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get news failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get news failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get news failed")
	}

	return newses, nil
}
