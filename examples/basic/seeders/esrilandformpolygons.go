package seeders

import (
	"encoding/json"
	"fmt"

	"github.com/TheAschr/pgseeder"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var esriLandformPolygonIdNS = uuid.MustParse("e0fefc56-9345-439a-a85e-164e447dfa2a")

func newEsriLandformPolygonID(permId string) uuid.UUID {
	return uuid.NewSHA1(esriLandformPolygonIdNS, []byte(permId))
}

func NewEsriLandformPolygons(fileName string, children []pgseeder.Config) pgseeder.Config {
	return pgseeder.Config{
		FileName:  fileName,
		ChunkSize: 100,
		HandleLine: func(batch *pgx.Batch, line []byte) error {
			type Properties struct {
				PermanentIdentifier string `json:"PERMANENT_IDENTIFIER"`
				Name                string `json:"NAME"`
				FCode               int    `json:"FCODE"`
				GazID               int    `json:"GAZ_ID"`
			}

			type Feature struct {
				Properties Properties  `json:"properties"`
				Geometry   interface{} `json:"geometry"`
			}

			var feature Feature

			if err := json.Unmarshal(line, &feature); err != nil {
				return fmt.Errorf("failed to unmarshall feature from line: %w", err)
			}

			id := newEsriLandformPolygonID(feature.Properties.PermanentIdentifier)

			batch.Queue(`
	INSERT INTO "EsriLandformPolygon" (
		"id", 
		"name",
		"featureCodeId",
		"gazId",
		"geoJSON"
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	) ON CONFLICT ("id") DO UPDATE SET
		"name" = $2,
		"featureCodeId" = $3,
		"gazId" = $4,
		"geoJSON" = $5
`,
				id,
				feature.Properties.Name,
				feature.Properties.FCode,
				feature.Properties.GazID,
				feature.Geometry,
			)

			return nil
		},
		Children: children,
	}
}
