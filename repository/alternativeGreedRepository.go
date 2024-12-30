package repository

import (
	"context"
	"fmt"
	"magmar/model"
	"magmar/model/dao"
	"magmar/util"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
)

type alternativeGreedRepository struct {
	conn   *resty.Client
	apiURL util.APIURL
}

// NewAlternativeGreedRepository ...
func NewAlternativeGreedRepository() GreedRepository {
	return &alternativeGreedRepository{
		conn:   resty.New(),
		apiURL: util.AlternativeURL,
	}
}

// GetGreedIndex ...
func (r *alternativeGreedRepository) GetGreedIndex(ctx context.Context) (*model.GreedIndex, error) {
	zlog.With(ctx).Infow(util.LogRepo)
	var daoIndex dao.AlternativeGreedIndex
	resp, err := r.conn.R().
		SetResult(&daoIndex).
		Get(fmt.Sprintf("%s/fng", r.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get greed index failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get greed index failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get greed index failed")
	}

	greedIndex, err := model.NewGreedIndexByAlternative(daoIndex.Index, daoIndex.IndexType)
	if err != nil {
		zlog.With(ctx).Errorw("Parsing greed index failed", "err", err)
		return nil, err
	}
	return greedIndex, nil
}
