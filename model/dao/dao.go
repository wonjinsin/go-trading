package dao

// Queryable ...
type Queryable interface {
	// GetQuery should guarantee order of parameters, can't use url.Values(it is ordering inside logic)
	GetQuery() string
}
