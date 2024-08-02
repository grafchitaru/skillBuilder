package postgresql

import (
	"context"
	"fmt"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"time"
)

func (s *Storage) GetTypeMaterials() ([]models.TypeMaterial, error) {
	const op = "storage.postgresql.GetTypeMaterials"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT id, created_at, updated_at, name, characteristic, xp FROM type_materials")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var typeMaterials []models.TypeMaterial
	for rows.Next() {
		var typeMaterial models.TypeMaterial
		if err = rows.Scan(&typeMaterial.Id, &typeMaterial.CreatedAt, &typeMaterial.UpdatedAt, &typeMaterial.Name, &typeMaterial.Characteristic, &typeMaterial.Xp); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		typeMaterials = append(typeMaterials, typeMaterial)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return typeMaterials, nil
}
