package main

import (
	"context"
	"fmt"

	core_pgx_pool "github.com/Yodau1/graphics-platform/internal/core/repository/postgres/pool/pgx"
)

func main() {
	ctx := context.Background()

	postgresPool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		panic(fmt.Errorf("init PostgreSQL pool: %w", err))
	}
	defer postgresPool.Close()

	fmt.Println("PostgreSQL connection established")
}
