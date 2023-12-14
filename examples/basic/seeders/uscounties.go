package seeders

import (
	"encoding/json"
	"fmt"

	"github.com/TheAschr/pgseeder"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var usCountyIdNS = uuid.MustParse("8b3fa54b-f47d-4e85-a177-962ce30129cf")

func newUsCountyID(longName string, stateAbbr string) uuid.UUID {
	return uuid.NewSHA1(usCountyIdNS, []byte(fmt.Sprintf("%s,%s", longName, stateAbbr)))
}

func NewUsCounties(fileName string, children []pgseeder.Config) pgseeder.Config {
	return pgseeder.Config{
		FileName:  fileName,
		ChunkSize: 100,
		HandleLine: func(batch *pgx.Batch, line []byte) error {
			type Properties struct {
				StateUsAbbreviation string `json:"stusab"`
				GeoID               string `json:"geoid"`
				NameLSAD            string `json:"namelsad"`
				Name                string `json:"name"`
			}

			type Feature struct {
				Properties Properties  `json:"properties"`
				Geometry   interface{} `json:"geometry"`
			}

			var feature Feature

			if err := json.Unmarshal(line, &feature); err != nil {
				return fmt.Errorf("failed to unmarshall feature from line: %w", err)
			}

			id := newUsCountyID(feature.Properties.NameLSAD, feature.Properties.StateUsAbbreviation)

			batch.Queue(`
	INSERT INTO "UsCounty" (
		"id", 
		"stateId",
		"territoryId",
		"stcoFipsCode",
		"longName",
		"shortName",
		"shapeGeoJSON",
		"deprecated"
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	) ON CONFLICT ("id") DO UPDATE SET
		"stateId" = $2,
		"territoryId" = $3,
		"stcoFipsCode" = $4,
		"longName" = $5,
		"shortName" = $6,
		"shapeGeoJSON" = $7,
		"deprecated" = $8
`,
				id,
				nil,
				nil,
				feature.Properties.GeoID,
				feature.Properties.NameLSAD,
				feature.Properties.Name,
				feature.Geometry,
				false, // deprecated
			)

			return nil
		},
		Children: children,
	}
}
