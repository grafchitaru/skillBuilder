package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"time"
)

func (s *Storage) CreateMaterial(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
	const op = "storage.postgresql.CreateMaterial"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := uuid.New()

	now := time.Now()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO materials(id, user_id, name, description, created_at, updated_at, type_id, xp, link)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
    `, id, userID, name, description, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), typeId, xp, link)
	if err != nil {
		return "", fmt.Errorf("%s exec: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) UpdateMaterial(material models.Material) error {
	const op = "storage.postgresql.UpdateMaterial"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	now := time.Now()

	_, err := s.pool.Exec(ctx, `
        UPDATE materials
        SET name=$1, description=$2, type_id=$3, link=$4, xp=$5, updated_at=$6
        WHERE id=$7 AND user_id=$8;
    `, material.Name, material.Description, material.TypeId, material.Link, material.Xp, now.Format("2006-01-02 15:04:05"), material.Id, material.UserId)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMaterial(userID, materialID string) error {
	const op = "storage.postgresql.DeleteMaterial"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        DELETE FROM materials
        WHERE id=$1 AND user_id=$2;
    `, materialID, userID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) GetMaterials(collectionID, userID string) ([]models.Material, error) {
	const op = "storage.postgresql.GetMaterials"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, `
		SELECT materials.*,
		       COALESCE(user_materials.completed, false) AS completed
		FROM materials
		INNER JOIN collection_materials ON materials.id = collection_materials.material_id
		LEFT JOIN user_materials ON materials.id = user_materials.material_id AND user_materials.user_id = $2
		WHERE collection_materials.collection_id = $1
	`, collectionID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var materials []models.Material
	for rows.Next() {
		var material models.Material
		var completed bool
		if err := rows.Scan(&material.Id, &material.CreatedAt, &material.UpdatedAt, &material.UserId, &material.Name, &material.Description, &material.TypeId, &material.Xp, &material.Link, &completed); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		material.Completed = completed
		materials = append(materials, material)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return materials, nil
}

func (s *Storage) GetMaterial(materialID string) (models.Material, error) {
	const op = "storage.postgresql.GetMaterial"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var material models.Material

	err := s.pool.QueryRow(ctx, "SELECT * FROM materials WHERE id = $1", materialID).Scan(&material.Id, &material.CreatedAt, &material.UpdatedAt, &material.UserId, &material.Name, &material.Description, &material.TypeId, &material.Xp, &material.Link)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Material{}, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return models.Material{}, fmt.Errorf("%s: %w", op, err)
	}

	return material, nil
}

func (s *Storage) AddMaterialToCollection(collectionID, materialID string) error {
	const op = "storage.postgresql.AddMaterialToCollection"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO collection_materials(collection_id, material_id)
        VALUES($1, $2);
    `, collectionID, materialID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) MarkMaterialAsCompleted(userID, materialID string) error {
	const op = "storage.postgresql.MarkMaterialAsCompleted"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        WITH upsert AS (
    UPDATE user_materials
    SET completed = true
    WHERE user_id = $1 AND material_id = $2
    RETURNING *
	)
	INSERT INTO user_materials (user_id, material_id, completed)
	SELECT $1, $2, true
	WHERE NOT EXISTS (SELECT * FROM upsert);
    `, userID, materialID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) MarkMaterialAsNotCompleted(userID, materialID string) error {
	const op = "storage.postgresql.MarkMaterialAsNotCompleted"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
WITH upsert AS (
    UPDATE user_materials
    SET completed = false
    WHERE user_id = $1 AND material_id = $2
    RETURNING *
)
INSERT INTO user_materials (user_id, material_id, completed)
SELECT $1, $2, false
WHERE NOT EXISTS (SELECT * FROM upsert);
    `, userID, materialID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) SearchMaterials(query string) ([]models.Material, error) {
	const op = "storage.postgresql.SearchMaterials"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT * FROM materials WHERE name LIKE '%'||$1||'%' OR description LIKE '%'||$1||'%'", query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var materials []models.Material
	for rows.Next() {
		var material models.Material
		if err := rows.Scan(&material.Id, &material.CreatedAt, &material.UpdatedAt, &material.UserId, &material.Name, &material.Description, &material.TypeId, &material.Xp, &material.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		materials = append(materials, material)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return materials, nil
}
