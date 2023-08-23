package queries

type ProductQueries struct {
	GetProductById GetProductByIdHandler
}

func NewProductQueries(getProductById GetProductByIdHandler) *ProductQueries {
	return &ProductQueries{GetProductById: getProductById}
}
