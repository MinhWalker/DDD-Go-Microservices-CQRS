package dto

import (
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/models"
	uuid "github.com/satori/go.uuid"
	"time"
)

func NewCreateProductCommand(productID string, name string, description string, price float64, createdAt time.Time, updatedAt time.Time) *models.CreateProductCommand {
	return &models.CreateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

func NewUpdateProductCommand(productID string, name string, description string, price float64, updatedAt time.Time) *models.UpdateProductCommand {
	return &models.UpdateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, UpdatedAt: updatedAt}
}

func NewDeleteProductCommand(productID uuid.UUID) *models.DeleteProductCommand {
	return &models.DeleteProductCommand{ProductID: productID}
}

func NewGetProductByIdQuery(productID uuid.UUID) *models.GetProductByIdQuery {
	return &models.GetProductByIdQuery{ProductID: productID}
}

func NewSearchProductQuery(text string, pagination *utils.Pagination) *models.SearchProductQuery {
	return &models.SearchProductQuery{Text: text, Pagination: pagination}
}
