package dto

import (
	"github.com/minhwalker/cqrs-microservices/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type GetProductByIdQuery struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

func NewGetProductByIdQuery(productID uuid.UUID) *GetProductByIdQuery {
	return &GetProductByIdQuery{ProductID: productID}
}

type SearchProductQuery struct {
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

func NewSearchProductQuery(text string, pagination *utils.Pagination) *SearchProductQuery {
	return &SearchProductQuery{Text: text, Pagination: pagination}
}
