package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/TheAschr/pgseeder"
	"github.com/TheAschr/pgseeder/examples/basic/seeders"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbauerster/mpb/v8"
)

const dataDir = "./data"

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgresql://user:password@localhost:5435/example")
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := initDb(pool); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	p := mpb.NewWithContext(ctx)

	s := pgseeder.New(pool, pgseeder.WithProgress(p))

	startTime := time.Now()

	if err := s.Run(ctx, []pgseeder.Config{
		seeders.NewUsers(
			path.Join(dataDir, "users.gz"),
			nil,
		),
		seeders.NewUsStates(
			path.Join(dataDir, "us-states.gz"),
			[]pgseeder.Config{
				seeders.NewUsCounties(path.Join(dataDir, "us-counties.gz"), nil),
			},
		),
		seeders.NewEsriLandformPolygons(path.Join(dataDir, "esri-landform-polygons.gz"), nil),
	}); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Printf("Completed in %v\n", time.Since(startTime).Round(time.Millisecond))
}

func initDb(pool *pgxpool.Pool) error {
	if _, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS "User" (
		"id" INT NOT NULL,
		"name" TEXT NOT NULL,
	
		CONSTRAINT "User_pkey" PRIMARY KEY ("id")
	)`); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS "EsriLandformPolygon" (
		"id" TEXT NOT NULL,
		"featureCodeId" INTEGER NOT NULL,
		"gazId" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"geoJSON" JSONB NOT NULL,
	
		CONSTRAINT "EsriLandformPolygon_pkey" PRIMARY KEY ("id")
	)`); err != nil {
		return fmt.Errorf("failed to create esri landform polygon table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS "UsState" (
		"id" TEXT NOT NULL,
		"fipsCode" TEXT NOT NULL,
		"alpha" TEXT NOT NULL,
		"name" TEXT NOT NULL,
		"shapeGeoJSON" JSONB NOT NULL,
		"districtOfColumbiaId" TEXT,
	
		CONSTRAINT "UsState_pkey" PRIMARY KEY ("id")
	)`); err != nil {
		return fmt.Errorf("failed to create us state table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS "UsCounty" (
		"id" TEXT NOT NULL,
		"stcoFipsCode" TEXT NOT NULL,
		"shortName" TEXT NOT NULL,
		"longName" TEXT NOT NULL,
		"deprecated" BOOLEAN NOT NULL DEFAULT false,
		"shapeGeoJSON" JSONB NOT NULL,
		"stateId" TEXT,
		"districtOfColumbiaId" TEXT,
		"territoryId" TEXT,
	
		CONSTRAINT "UsCounty_pkey" PRIMARY KEY ("id")
	)`); err != nil {
		return fmt.Errorf("failed to create us county table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'UsCounty_stateId_fkey') THEN
			ALTER TABLE "UsCounty" ADD CONSTRAINT "UsCounty_stateId_fkey" FOREIGN KEY ("stateId") REFERENCES "UsState"("id") ON DELETE SET NULL ON UPDATE CASCADE;
		END IF;
	END;
	$$;
	`); err != nil {
		return fmt.Errorf("failed to create us state reference key: %w", err)
	}

	return nil
}
