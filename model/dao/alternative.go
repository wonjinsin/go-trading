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
	Index     string
	IndexType AlternativeGreedIndexType
}
