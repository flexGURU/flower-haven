package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
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
	db      *Store
}

func NewProductRepository(db *Store) *ProductRepository {
	return &ProductRepository{
		db:      db,
		queries: generated.New(db.pool),
	}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (*repository.Product, error) {
	err := pr.db.ExecTx(ctx, func(q *generated.Queries) error {
		generatedProduct, err := q.CreateProduct(ctx, generated.CreateProductParams{
			Name:          product.Name,
			Description:   product.Description,
			Price:         pkg.Float64ToPgTypeNumeric(product.Price),
			CategoryID:    int64(product.CategoryID),
			HasStems:      product.HasStems,
			IsMessageCard: product.IsMessageCard,
			IsFlowers:     product.IsFlowers,
			IsAddOn:       product.IsAddOn,
			ImageUrl:      product.ImageUrl,
			StockQuantity: product.StockQuantity,
		})
		if err != nil {
			if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
				return pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
			}
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating product: %s", err.Error())
		}

		product.ID = uint32(generatedProduct.ID)
		product.CreatedAt = generatedProduct.CreatedAt
		product.CategoryData = nil

		if product.HasStems {
			for _, stem := range product.Stems {
				_, err := q.CreateProductStem(ctx, generated.CreateProductStemParams{
					ProductID: int64(product.ID),
					StemCount: int64(stem.StemCount),
					Price:     pkg.Float64ToPgTypeNumeric(stem.Price),
				})
				if err != nil {
					if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
						return pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
					}
					return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating product stem: %s", err.Error())
				}
			}
		}

		return nil
	})

	return product, err
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
		HasStems:      generatedProduct.HasStems,
		IsMessageCard: generatedProduct.IsMessageCard,
		IsFlowers:     generatedProduct.IsFlowers,
		IsAddOn:       generatedProduct.IsAddOn,
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

	if product.HasStems {
		var stems []repository.ProductStem
		productStems, err := pr.queries.GetProductStemsByProductID(ctx, int64(product.ID))
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting product stems: %s", err.Error())
			}
		}

		for _, stem := range productStems {
			stems = append(stems, repository.ProductStem{
				ID:        uint32(stem.ID),
				ProductID: uint32(stem.ProductID),
				StemCount: uint32(stem.StemCount),
				Price:     pkg.PgTypeNumericToFloat64(stem.Price),
			})
		}
		product.Stems = stems
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
		HasStems:      pgtype.Bool{Valid: false},
		IsMessageCard: pgtype.Bool{Valid: false},
		IsFlowers:     pgtype.Bool{Valid: false},
		IsAddOn:       pgtype.Bool{Valid: false},
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
	if product.HasStems != nil {
		params.HasStems = pgtype.Bool{Bool: *product.HasStems, Valid: true}
	}
	if product.IsMessageCard != nil {
		params.IsMessageCard = pgtype.Bool{Bool: *product.IsMessageCard, Valid: true}
	}
	if product.IsFlowers != nil {
		params.IsFlowers = pgtype.Bool{Bool: *product.IsFlowers, Valid: true}
	}
	if product.IsAddOn != nil {
		params.IsAddOn = pgtype.Bool{Bool: *product.IsAddOn, Valid: true}
	}

	err := pr.db.ExecTx(ctx, func(q *generated.Queries) error {
		if product.Stems != nil {
			// delete all existing stems and recreate them
			if err := q.DeleteProductStemsByProductID(ctx, int64(product.ID)); err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting existing product stems: %s", err.Error())
			}

			for _, stem := range product.Stems {
				stemProductId := product.ID
				_, err := q.CreateProductStem(ctx, generated.CreateProductStemParams{
					ProductID: int64(stemProductId),
					StemCount: int64(*stem.StemCount),
					Price:     pkg.Float64ToPgTypeNumeric(*stem.Price),
				})
				if err != nil {
					if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
						return pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
					}
					return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating product stem: %s", err.Error())
				}
			}
		}

		_, err := pr.queries.UpdateProduct(ctx, params)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with ID %d not found", product.ID)
			}
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error updating product: %s", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pr.GetProductByID(ctx, int64(product.ID))
}

// func (pr *ProductRepository) ListProducts(ctx context.Context, filter *repository.ProductFilter) ([]*repository.Product, *pkg.Pagination, error) {
// 	paramsListProducts := generated.ListProductsParams{
// 		Limit:       int32(filter.Pagination.PageSize),
// 		Offset:      pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
// 		Search:      pgtype.Text{Valid: false},
// 		PriceFrom:   pgtype.Float8{Valid: false},
// 		PriceTo:     pgtype.Float8{Valid: false},
// 		CategoryIds: nil,
// 	}

// 	paramsCountProducts := generated.ListCountProductsParams{
// 		Search:      pgtype.Text{Valid: false},
// 		PriceFrom:   pgtype.Float8{Valid: false},
// 		PriceTo:     pgtype.Float8{Valid: false},
// 		CategoryIds: nil,
// 	}

// 	if filter.Search != nil {
// 		search := strings.ToLower(*filter.Search)
// 		paramsListProducts.Search = pgtype.Text{
// 			Valid:  true,
// 			String: "%" + search + "%",
// 		}
// 		paramsCountProducts.Search = pgtype.Text{
// 			Valid:  true,
// 			String: "%" + search + "%",
// 		}
// 	}

// 	if filter.PriceFrom != nil && filter.PriceTo != nil {
// 		paramsListProducts.PriceFrom = pgtype.Float8{
// 			Valid:   true,
// 			Float64: *filter.PriceFrom,
// 		}
// 		paramsListProducts.PriceTo = pgtype.Float8{
// 			Valid:   true,
// 			Float64: *filter.PriceTo,
// 		}

// 		paramsCountProducts.PriceFrom = pgtype.Float8{
// 			Valid:   true,
// 			Float64: *filter.PriceFrom,
// 		}
// 		paramsCountProducts.PriceTo = pgtype.Float8{
// 			Valid:   true,
// 			Float64: *filter.PriceTo,
// 		}
// 	}

// 	if filter.CategoryIDs != nil {
// 		paramsListProducts.CategoryIds = make([]int32, len(*filter.CategoryIDs))
// 		paramsCountProducts.CategoryIds = make([]int32, len(*filter.CategoryIDs))
// 		for _, id := range *filter.CategoryIDs {
// 			paramsListProducts.CategoryIds = append(paramsListProducts.CategoryIds, int32(id))
// 			paramsCountProducts.CategoryIds = append(paramsCountProducts.CategoryIds, int32(id))
// 		}
// 	}

// 	generatedProducts, err := pr.queries.ListProducts(ctx, paramsListProducts)
// 	if err != nil {
// 		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing products: %s", err.Error())
// 	}

// 	totalCount, err := pr.queries.ListCountProducts(ctx, paramsCountProducts)
// 	if err != nil {
// 		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting products: %s", err.Error())
// 	}

// 	products := make([]*repository.Product, len(generatedProducts))
// 	for i, p := range generatedProducts {
// 		products[i] = &repository.Product{
// 			ID:            uint32(p.ID),
// 			Name:          p.Name,
// 			Description:   p.Description,
// 			Price:         pkg.PgTypeNumericToFloat64(p.Price),
// 			CategoryID:    uint32(p.CategoryID),
// 			ImageUrl:      p.ImageUrl,
// 			StockQuantity: p.StockQuantity,
// 			DeletedAt:     nil,
// 			CreatedAt:     p.CreatedAt,
// 			CategoryData:  nil,
// 		}

// 		if p.CategoryID_2.Valid {
// 			products[i].CategoryData = &repository.Category{
// 				ID:          uint32(p.CategoryID_2.Int64),
// 				Name:        p.CategoryName.String,
// 				Description: p.CategoryDescription.String,
// 			}
// 		}
// 	}

// 	return products, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
// }

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
		paramsListProducts.PriceFrom = pgtype.Float8{Valid: true, Float64: *filter.PriceFrom}
		paramsListProducts.PriceTo = pgtype.Float8{Valid: true, Float64: *filter.PriceTo}
		paramsCountProducts.PriceFrom = pgtype.Float8{Valid: true, Float64: *filter.PriceFrom}
		paramsCountProducts.PriceTo = pgtype.Float8{Valid: true, Float64: *filter.PriceTo}
	}

	if filter.CategoryIDs != nil {
		paramsListProducts.CategoryIds = make([]int32, 0, len(*filter.CategoryIDs))
		paramsCountProducts.CategoryIds = make([]int32, 0, len(*filter.CategoryIDs))
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
		product := &repository.Product{
			ID:            uint32(p.ID),
			Name:          p.Name,
			Description:   p.Description,
			Price:         pkg.PgTypeNumericToFloat64(p.Price),
			CategoryID:    uint32(p.CategoryID),
			ImageUrl:      p.ImageUrl,
			HasStems:      p.HasStems,
			IsMessageCard: p.IsMessageCard,
			IsFlowers:     p.IsFlowers,
			IsAddOn:       p.IsAddOn,
			StockQuantity: p.StockQuantity,
			DeletedAt:     nil,
			CreatedAt:     p.CreatedAt,
			CategoryData:  nil,
		}

		if p.CategoryID_2.Valid {
			product.CategoryData = &repository.Category{
				ID:          uint32(p.CategoryID_2.Int64),
				Name:        p.CategoryName.String,
				Description: p.CategoryDescription.String,
			}
		}

		// Unmarshal stems JSON
		if p.Stems != nil {
			var stems []repository.ProductStem

			switch v := p.Stems.(type) {
			case []byte:
				if err := json.Unmarshal(v, &stems); err != nil {
					return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshaling stems: %s", err.Error())
				}
			case []interface{}:
				raw, err := json.Marshal(v)
				if err != nil {
					return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error marshaling stems interface{}: %s", err.Error())
				}
				if err := json.Unmarshal(raw, &stems); err != nil {
					return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshaling stems: %s", err.Error())
				}
			default:
				return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "unexpected stems type: %T", p.Stems)
			}

			product.Stems = stems
		}

		products[i] = product
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

func (pr *ProductRepository) GetDashboardData(ctx context.Context) (interface{}, error) {
	totalRevenue, err := pr.queries.TotalRevenue(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching total revenue: %s", err.Error())
	}
	totalProducts, err := pr.queries.TotalProducts(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching total products: %s", err.Error())
	}
	totalOrders, err := pr.queries.TotalOrders(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching total orders: %s", err.Error())
	}
	activeSubscriptions, err := pr.queries.ActiveSubscriptions(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching active subscriptions: %s", err.Error())
	}
	recentOrders, err := pr.queries.GetRecentOrders(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching recent orders: %s", err.Error())
	}
	categoriesData, err := pr.queries.GetCategoriesWithProductCount(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching categories with product count: %s", err.Error())
	}

	return map[string]interface{}{
		"total_revenue":        totalRevenue,
		"total_products":       totalProducts,
		"total_orders":         totalOrders,
		"active_subscriptions": activeSubscriptions,
		"recent_orders":        recentOrders,
		"categories":           categoriesData,
	}, nil
}
