package resources_postgres_repository

import core_postgres_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool"

type ResourcesRepository struct {
	pool core_postgres_pool.Pool
}

func NewResourcesRepository(pool core_postgres_pool.Pool) *ResourcesRepository {
	return &ResourcesRepository{
		pool: pool,
	}
}
