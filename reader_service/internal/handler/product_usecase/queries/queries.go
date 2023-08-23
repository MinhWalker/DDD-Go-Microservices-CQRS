package queries

type ProductQueries struct {
	GetProductById GetProductByIdHandler
	SearchProduct  SearchProductHandler
}

func NewProductQueries(getProductById GetProductByIdHandler, searchProduct SearchProductHandler) *ProductQueries {
	return &ProductQueries{GetProductById: getProductById, SearchProduct: searchProduct}
}
