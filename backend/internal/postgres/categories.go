package postgres

import (
	"context"
	"strings"

	"github.com/flexGURU/flower-haven/backend/internal/postgres/generated"
	"github.com/flexGURU/flower-haven/backend/internal/repository"
	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.CategoryRepository = (*CategoryRepository)(nil)

type CategoryRepository struct {
	queries *generated.Queries
}

func NewCategoryRepository(queries *generated.Queries) *CategoryRepository {
	return &CategoryRepository{queries: queries}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *repository.Category) (*repository.Category, error) {
	generatedCategory, err := cr.queries.CreateCategory(ctx, generated.CreateCategoryParams{
		Name:        category.Name,
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating category: %s", err.Error())
	}

	category.ID = uint32(generatedCategory.ID)
	category.CreatedAt = generatedCategory.CreatedAt

	return category, nil
}

func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*repository.Category, error) {
	generatedCategory, err := cr.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "category with id %d not found", id)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error fetching category by id: %s", err.Error())
	}

	category := &repository.Category{
		ID:          uint32(generatedCategory.ID),
		Name:        generatedCategory.Name,
		Description: generatedCategory.Description,
		ImageUrl:    generatedCategory.ImageUrl,
		DeletedAt:   &generatedCategory.DeletedAt.Time,
		CreatedAt:   generatedCategory.CreatedAt,
	}

	return category, nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, category *repository.UpdateCategory) (*repository.Category, error) {
	params := generated.UpdateCategoryParams{
		ID:          int64(category.ID),
		Name:        pgtype.Text{Valid: false},
		Description: pgtype.Text{Valid: false},
		ImageUrl:    nil,
	}

	if category.Name != nil {
		params.Name = pgtype.Text{
			Valid:  true,
			String: *category.Name,
		}
	}
	if category.Description != nil {
		params.Description = pgtype.Text{
			Valid:  true,
			String: *category.Description,
		}
	}
	if category.ImageUrl != nil {
		params.ImageUrl = *category.ImageUrl
	}

	generatedCategory, err := cr.queries.UpdateCategory(ctx, params)
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "category with id %d not found", category.ID)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating category: %s", err.Error())
	}

	return &repository.Category{
		ID:          uint32(generatedCategory.ID),
		Name:        generatedCategory.Name,
		Description: generatedCategory.Description,
		ImageUrl:    generatedCategory.ImageUrl,
		DeletedAt:   &generatedCategory.DeletedAt.Time,
		CreatedAt:   generatedCategory.CreatedAt,
	}, nil
}

func (cr *CategoryRepository) ListCategories(ctx context.Context, filter *repository.CategoryFilter) ([]*repository.Category, *pkg.Pagination, error) {
	paramListCategories := generated.ListCategoriesParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: int32(filter.Pagination.Page * filter.Pagination.PageSize),
		Search: pgtype.Text{Valid: false},
	}

	paramCountCategories := pgtype.Text{Valid: false}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramListCategories.Search = pgtype.Text{
			Valid:  true,
			String: "%" + search + "%",
		}
		paramCountCategories.Valid = true
		paramCountCategories.String = "%" + search + "%"
	}

	generatedCategories, err := cr.queries.ListCategories(ctx, paramListCategories)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing categories: %s", err.Error())
	}

	totalCount, err := cr.queries.ListCategoriesCount(ctx, paramCountCategories)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting categories: %s", err.Error())
	}

	categoryList := make([]*repository.Category, len(generatedCategories))
	for i, cat := range generatedCategories {
		categoryList[i] = &repository.Category{
			ID:          uint32(cat.ID),
			Name:        cat.Name,
			Description: cat.Description,
			ImageUrl:    cat.ImageUrl,
			DeletedAt:   &cat.DeletedAt.Time,
			CreatedAt:   cat.CreatedAt,
		}
	}

	return categoryList, pkg.CalculatePagination(uint32(totalCount), filter.Pagination.PageSize, filter.Pagination.Page), nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	if err := cr.queries.DeleteCategory(ctx, id); err != nil {
		if pkg.PgxErrorCode(err) == pkg.NOT_FOUND_ERROR {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "category with id %d not found", id)
		}
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting category by id: %s", err.Error())
	}
	return nil
}
