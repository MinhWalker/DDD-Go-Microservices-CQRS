package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type productRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewProductRepository(log logger.Logger, cfg *config.Config, db *pgxpool.Pool) RepositoryReader {
	return &productRepository{log: log, cfg: cfg, db: db}
}

func (p *productRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.CreateProduct")
	defer span.Finish()

	var created models.Product
	if err := p.db.QueryRow(ctx, createProductQuery, &product.ProductID, &product.Name, &product.Description, &product.Price).Scan(
		&created.ProductID,
		&created.Name,
		&created.Description,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}

	return &created, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.UpdateProduct")
	defer span.Finish()

	var prod models.Product
	if err := p.db.QueryRow(
		ctx,
		updateProductQuery,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ProductID,
	).Scan(&prod.ProductID, &prod.Name, &prod.Description, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &prod, nil
}

func (p *productRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.GetProductById")
	defer span.Finish()

	var product models.Product
	if err := p.db.QueryRow(ctx, getProductByIdQuery, uuid).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &product, nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.DeleteProductByID")
	defer span.Finish()

	_, err := p.db.Exec(ctx, deleteProductByIdQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}

	return nil
}

func (p *productRepository) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresRepository.Search")
	defer span.Finish()

	query := `
		SELECT product_id, name, description, price 
		FROM products
		WHERE name ILIKE $1 OR description ILIKE $1
		LIMIT $2 OFFSET $3
	`

	// Construct the search query to perform case-insensitive search using ILIKE
	searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(search))

	rows, err := p.db.Query(ctx, query, searchTerm, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		p.traceErr(span, err)
		return nil, errors.Wrap(err, "Query")
	}
	defer rows.Close()

	products := make([]*models.Product, 0)

	for rows.Next() {
		var prod models.Product
		if err := rows.Scan(&prod.ProductID, &prod.Name, &prod.Description, &prod.Price); err != nil {
			p.traceErr(span, err)
			return nil, errors.Wrap(err, "Scan")
		}
		products = append(products, &prod)
	}

	if err := rows.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return nil, errors.Wrap(err, "rows.Err")
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM products
		WHERE name ILIKE $1 OR description ILIKE $1
	`

	var count int64
	err = p.db.QueryRow(ctx, countQuery, searchTerm).Scan(&count)
	if err != nil {
		p.traceErr(span, err)
		return nil, errors.Wrap(err, "QueryRow")
	}

	return models.NewProductListWithPagination(products, count, pagination), nil
}

func (p *productRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
