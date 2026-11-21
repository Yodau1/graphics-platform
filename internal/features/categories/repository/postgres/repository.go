package categories_postgres_repository

import core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"

// CategoriesRepository выполняет операции с категориями в PostgreSQL
type CategoriesRepository struct {
	pool core_postgres_pool.Pool
}

func NewCategoriesRepository(pool core_postgres_pool.Pool) *CategoriesRepository {
	return &CategoriesRepository{
		pool: pool,
	}
}
