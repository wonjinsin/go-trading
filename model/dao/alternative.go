package dao

// AlternativeGreedIndexType ...
type AlternativeGreedIndexType string

// GreedIndexType ...
const (
	AlternativeGreedIndexTypeNone         AlternativeGreedIndexType = "None"
	AlternativeGreedIndexTypeExtremeFear  AlternativeGreedIndexType = "Extreme Fear"
	AlternativeGreedIndexTypeFear         AlternativeGreedIndexType = "Fear"
	AlternativeGreedIndexTypeNeutral      AlternativeGreedIndexType = "Neutral"
	AlternativeGreedIndexTypeGreed        AlternativeGreedIndexType = "Greed"
	AlternativeGreedIndexTypeExtremeGreed AlternativeGreedIndexType = "Extreme Greed"
)

// AlternativeGreedIndex ...
type AlternativeGreedIndex struct {
	Data []struct {
		Value               string                    `json:"value"`
		ValueClassification AlternativeGreedIndexType `json:"value_classification"`
	} `json:"data"`
}
