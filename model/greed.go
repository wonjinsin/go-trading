package model

import (
	"errors"
	"magmar/model/dao"
	"strconv"
)

type GreedIndexType string

// GreedIndexType ...
const (
	GreedIndexTypeNone         GreedIndexType = "None"
	GreedIndexTypeExtremeFear  GreedIndexType = "Extreme Fear"
	GreedIndexTypeFear         GreedIndexType = "Fear"
	GreedIndexTypeNeutral      GreedIndexType = "Neutral"
	GreedIndexTypeGreed        GreedIndexType = "Greed"
	GreedIndexTypeExtremeGreed GreedIndexType = "Extreme Greed"
)

// GreedIndex ...
type GreedIndex struct {
	Index     uint
	IndexType GreedIndexType
}

// NewGreedIndexByAlternative ...
func NewGreedIndexByAlternative(index string, alternativeIndexType dao.AlternativeGreedIndexType) (*GreedIndex, error) {
	indexUint, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		return nil, err
	}

	var indexType GreedIndexType = GreedIndexTypeNone
	switch alternativeIndexType {
	case dao.AlternativeGreedIndexTypeExtremeFear:
		indexType = GreedIndexTypeExtremeFear
	case dao.AlternativeGreedIndexTypeFear:
		indexType = GreedIndexTypeFear
	case dao.AlternativeGreedIndexTypeNeutral:
		indexType = GreedIndexTypeNeutral
	case dao.AlternativeGreedIndexTypeGreed:
		indexType = GreedIndexTypeGreed
	case dao.AlternativeGreedIndexTypeExtremeGreed:
		indexType = GreedIndexTypeExtremeGreed
	default:
		return nil, errors.New("invalid index type")
	}

	return &GreedIndex{
		Index:     uint(indexUint),
		IndexType: indexType,
	}, nil
}
