package database

import (
	"context"
	"log"

	"gorm.io/gorm"
)

// Query param context keys - using string type to match middleware
const (
	// PageContextKey is the key for storing page number in context
	PageContextKey = "page"
	// LimitContextKey is the key for storing limit in context
	LimitContextKey = "limit"
	// SearchContextKey is the key for storing search parameters in context
	SearchContextKey = "search"
	// SortContextKey is the key for storing sort direction in context
	SortContextKey = "sort"
	// SortColumnContextKey is the key for storing sort column in context
	SortColumnContextKey = "sort_column"
)

// SelectAll retrieves records of the model type with pagination, sorting, and search.
// It automatically extracts page, limit, search, sort, and sortColumn parameters from context.
// Default limit is 10, maximum allowed is 30.
// Optional conditions can be passed to filter the query.
func (r *Repository[T]) SelectAll(
	ctx context.Context,
	preloads []string,
	conditions ...func(*gorm.DB) *gorm.DB,
) ([]T, error) {
	// Extract page from context (default to 1)
	defaultPage := 1
	if ctxPage, ok := ctx.Value(PageContextKey).(int); ok && ctxPage > 0 {
		defaultPage = ctxPage
	}

	// Extract limit from context (default to 10, max 30)
	defaultLimit := 10
	if ctxLimit, ok := ctx.Value(LimitContextKey).(int); ok && ctxLimit > 0 {
		defaultLimit = ctxLimit
		// Enforce maximum limit of 30
		if defaultLimit > 30 {
			defaultLimit = 30
		}
	}

	// Extract search from context
	var search map[string]string
	if ctxSearch, ok := ctx.Value(SearchContextKey).(map[string]string); ok {
		search = ctxSearch
	}
	if search == nil {
		search = make(map[string]string)
	}

	// Initialize models slice
	var models []T
	// Initialize query
	query := r.db.DB.WithContext(ctx)

	// Apply preloads
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	// Extract sort from context (default to desc)
	sortBy := "desc"
	if ctxSort, ok := ctx.Value(SortContextKey).(string); ok && ctxSort != "" {
		sortBy = ctxSort
	}

	// Extract sortColumn from context (default to created_at)
	sortColumn := "created_at"
	if ctxSortColumn, ok := ctx.Value(SortColumnContextKey).(string); ok && ctxSortColumn != "" {
		sortColumn = ctxSortColumn
	}

	query = query.Order(sortColumn + " " + sortBy)

	// Apply search
	for key, value := range search {
		if value != "" && len(value) >= 3 {
			columnName := key
			// If key starts with "search_", strip it to get the actual column name
			if len(key) > 7 && key[:7] == "search_" {
				columnName = key[7:]
			} else if key == "search" {
				// Handle generic search if needed, otherwise skip to avoid "column 'search' does not exist"
				continue
			}
			query = query.Where(columnName+" LIKE ?", "%"+value+"%")
		}
	}

	// Apply optional static conditions
	for _, condition := range conditions {
		query = condition(query)
	}

	// Calculate offset based on page number (page is 1-indexed, so subtract 1)
	offset := (defaultPage - 1) * defaultLimit
	if offset < 0 {
		offset = 0
	}
	query = query.Limit(defaultLimit).Offset(offset)

	// Find records
	err := query.Find(&models).Error
	if err != nil {
		log.Printf("Database select all error: %v", err)
	}
	return models, err
}

func (r *Repository[T]) SelectByID(ctx context.Context, id uint, preloads []string) (*T, error) {
	var model T
	query := r.db.DB.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	err := query.First(&model, id).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Database select by ID error: %v, id: %d", err, id)
		}
		return nil, err
	}
	return &model, nil
}

// Create inserts a new record into the database.
func (r *Repository[T]) Create(ctx context.Context, model *T) error {
	err := r.db.DB.WithContext(ctx).Create(model).Error
	if err != nil {
		log.Printf("Database create error: %v", err)
	}
	return err
}
