package models

import (
	"time"

	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	readerService "github.com/minhwalker/cqrs-microservices/reader_service/internal/dto/proto/product_reader"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Product struct {
	models.Product
}

// Command
type CreateProductCommand struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type UpdateProductCommand struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

// Queries
type SearchProductQuery struct {
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

type GetProductByIdQuery struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

// ProductsList products list response with pagination
type ProductsList struct {
	TotalCount int64      `json:"totalCount" bson:"totalCount"`
	TotalPages int64      `json:"totalPages" bson:"totalPages"`
	Page       int64      `json:"page" bson:"page"`
	Size       int64      `json:"size" bson:"size"`
	HasMore    bool       `json:"hasMore" bson:"hasMore"`
	Products   []*Product `json:"products" bson:"products"`
}

func NewProductListWithPagination(products []*Product, count int64, pagination *utils.Pagination) *ProductsList {
	return &ProductsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Products:   products,
	}
}

func ProductToGrpcMessage(product *Product) *readerService.Product {
	return &readerService.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ProductListToGrpc(products *ProductsList) *readerService.SearchRes {
	list := make([]*readerService.Product, 0, len(products.Products))
	for _, product := range products.Products {
		list = append(list, ProductToGrpcMessage(product))
	}

	return &readerService.SearchRes{
		TotalCount: products.TotalCount,
		TotalPages: products.TotalPages,
		Page:       products.Page,
		Size:       products.Size,
		HasMore:    products.HasMore,
		Products:   list,
	}
}
