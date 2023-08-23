package commands

type ProductCommands struct {
	CreateProduct CreateProductCmdHandler
	UpdateProduct UpdateProductCmdHandler
	DeleteProduct DeleteProductCmdHandler
}

func NewProductCommands(createProduct CreateProductCmdHandler, updateProduct UpdateProductCmdHandler, deleteProduct DeleteProductCmdHandler) *ProductCommands {
	return &ProductCommands{CreateProduct: createProduct, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}
