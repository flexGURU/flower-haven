package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	queries *generated.Queries
}

func NewProductRepository(queries *generated.Queries) *ProductRepository {
	return &ProductRepository{queries: queries}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (*repository.Product, error) {
	generatedProduct, err := pr.queries.CreateProduct(ctx, generated.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         pkg.Float64ToPgTypeNumeric(product.Price),
		CategoryID:    int64(product.CategoryID),
		ImageUrl:      product.ImageUrl,
		StockQuantity: product.StockQuantity,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating product: %s", err.Error())
	}

	product.ID = uint32(generatedProduct.ID)
	product.CreatedAt = generatedProduct.CreatedAt
	product.CategoryData = nil

	return product, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, id int64) (*repository.Product, error) {
	generatedProduct, err := pr.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching product by id: %s", err.Error())
	}

	product := &repository.Product{
		ID:            uint32(generatedProduct.ID),
		Name:          generatedProduct.Name,
		Description:   generatedProduct.Description,
		Price:         pkg.PgTypeNumericToFloat64(generatedProduct.Price),
		CategoryID:    uint32(generatedProduct.CategoryID),
		ImageUrl:      generatedProduct.ImageUrl,
		StockQuantity: generatedProduct.StockQuantity,
		DeletedAt:     nil,
		CreatedAt:     generatedProduct.CreatedAt,
		CategoryData:  nil,
	}

	if generatedProduct.DeletedAt.Valid {
		product.DeletedAt = &generatedProduct.DeletedAt.Time
	}

	if generatedProduct.CategoryID_2.Valid {
		product.CategoryData = &repository.Category{
			ID:          uint32(generatedProduct.CategoryID_2.Int64),
			Name:        generatedProduct.CategoryName.String,
			Description: generatedProduct.CategoryDescription.String,
		}
	}

	return product, nil
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *repository.UpdateProduct) (*repository.Product, error) {
	params := generated.UpdateProductParams{
		ID:            int64(product.ID),
		Name:          pgtype.Text{Valid: false},
		Description:   pgtype.Text{Valid: false},
		Price:         pgtype.Numeric{Valid: false},
		CategoryID:    pgtype.Int8{Valid: false},
		ImageUrl:      nil,
		StockQuantity: pgtype.Int8{Valid: false},
	}

	if product.Name != nil {
		params.Name = pgtype.Text{
			Valid:  true,
			String: *product.Name,
		}
	}
	if product.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *product.Description,
		}
	}
	if product.Price != nil {
		params.Price = pkg.Float64ToPgTypeNumeric(*product.Price)
	}
	if product.CategoryID != nil {
		params.CategoryID = pgtype.Int8{
			Valid: true,
			Int64: int64(*product.CategoryID),
		}
	}
	if product.ImageURL != nil {
		params.ImageUrl = *product.ImageURL
	}
	if product.StockQuantity != nil {
		params.StockQuantity = pgtype.Int8{Int64: int64(*product.StockQuantity), Valid: true}
	}

	generatedProduct, err := pr.queries.UpdateProduct(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", product.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating product: %s", err.Error())
	}

	productData := &repository.Product{
		ID:            uint32(generatedProduct.ID),
		Name:          generatedProduct.Name,
		Description:   generatedProduct.Description,
		CategoryID:    uint32(generatedProduct.CategoryID),
		Price:         pkg.PgTypeNumericToFloat64(generatedProduct.Price),
		ImageUrl:      generatedProduct.ImageUrl,
		StockQuantity: generatedProduct.StockQuantity,
		CreatedAt:     generatedProduct.CreatedAt,
		DeletedAt:     nil,
	}

	return productData, nil
}

func (pr *ProductRepository) ListProducts(ctx context.Context, filter *repository.ProductFilter) ([]*repository.Product, *pkg.Pagination, error) {
	paramsListProducts := generated.ListProductsParams{
		Limit:       int32(filter.Pagination.PageSize),
		Offset:      pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
		Search:      pgtype.Text{Valid: false},
		PriceFrom:   pgtype.Float8{Valid: false},
		PriceTo:     pgtype.Float8{Valid: false},
		CategoryIds: nil,
	}

	paramsCountProducts := generated.ListCountProductsParams{
		Search:      pgtype.Text{Valid: false},
		PriceFrom:   pgtype.Float8{Valid: false},
		PriceTo:     pgtype.Float8{Valid: false},
		CategoryIds: nil,
	}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramsListProducts.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramsCountProducts.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
	}

	if filter.PriceFrom != nil && filter.PriceTo != nil {
		paramsListProducts.PriceFrom = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceFrom,
		}
		paramsListProducts.PriceTo = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceTo,
		}

		paramsCountProducts.PriceFrom = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceFrom,
		}
		paramsCountProducts.PriceTo = pgtype.Float8{
			Valid:   true,
			Float64: *filter.PriceTo,
		}
	}

	if filter.CategoryIDs != nil {
		paramsListProducts.CategoryIds = make([]int32, len(*filter.CategoryIDs))
		paramsCountProducts.CategoryIds = make([]int32, len(*filter.CategoryIDs))
		for _, id := range *filter.CategoryIDs {
			paramsListProducts.CategoryIds = append(paramsListProducts.CategoryIds, int32(id))
			paramsCountProducts.CategoryIds = append(paramsCountProducts.CategoryIds, int32(id))
		}
	}

	generatedProducts, err := pr.queries.ListProducts(ctx, paramsListProducts)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing products: %s", err.Error())
	}

	totalCount, err := pr.queries.ListCountProducts(ctx, paramsCountProducts)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting products: %s", err.Error())
	}

	products := make([]*repository.Product, len(generatedProducts))
	for i, p := range generatedProducts {
		products[i] = &repository.Product{
			ID:            uint32(p.ID),
			Name:          p.Name,
			Description:   p.Description,
			Price:         pkg.PgTypeNumericToFloat64(p.Price),
			CategoryID:    uint32(p.CategoryID),
			ImageUrl:      p.ImageUrl,
			StockQuantity: p.StockQuantity,
			DeletedAt:     nil,
			CreatedAt:     p.CreatedAt,
			CategoryData:  nil,
		}

		if p.CategoryID_2.Valid {
			products[i].CategoryData = &repository.Category{
				ID:          uint32(p.CategoryID_2.Int64),
				Name:        p.CategoryName.String,
				Description: p.CategoryDescription.String,
			}
		}
	}

	return products, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	if err := pr.queries.DeleteProduct(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting product: %s", err.Error())
	}
	return nil
}
