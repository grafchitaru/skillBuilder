package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"time"
)

func (s *Storage) CreateCollection(userID, name, description string) (string, error) {
	const op = "storage.postgresql.CreateCollection"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id := uuid.New()

	now := time.Now()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO collections(id, user_id, name, description, created_at, updated_at)
        VALUES($1, $2, $3, $4, $5, $6);
    `, id, userID, name, description, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))
	if err != nil {
		return "", fmt.Errorf("%s exec: %w", op, err)
	}

	return id.String(), nil
}

func (s *Storage) UpdateCollection(collection models.Collection) error {
	const op = "storage.postgresql.UpdateCollection"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	now := time.Now()

	_, err := s.pool.Exec(ctx, `
        UPDATE collections
        SET name=$1, description=$2, updated_at=$3
        WHERE id=$4 AND user_id=$5;
    `, collection.Name, collection.Description, now.Format("2006-01-02 15:04:05"), collection.Id, collection.UserId)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCollections(userID string) ([]models.Collection, error) {
	const op = "storage.postgresql.GetCollections"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//TODO Need optimize SQL Request + add indexes
	rows, err := s.pool.Query(ctx, "SELECT collections.*,"+
		"(SELECT sum(materials.xp) "+
		"FROM materials WHERE materials.id IN "+
		"(SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		") AS sum_xp, "+
		"( "+
		"SELECT (SELECT sum(materials.xp) FROM materials WHERE materials.id = user_materials.material_id) "+
		"FROM user_materials WHERE user_materials.material_id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		"AND user_materials.completed = true AND user_materials.user_id = $1 "+
		") AS xp "+
		"FROM collections", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var collection models.Collection
		if err = rows.Scan(&collection.Id, &collection.CreatedAt, &collection.UpdatedAt, &collection.UserId, &collection.Name, &collection.Description, &collection.SumXp, &collection.Xp); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return collections, nil
}

func (s *Storage) GetUserCollections(userID string) ([]models.Collection, error) {
	const op = "storage.postgresql.GetUserCollections"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//TODO Need optimize SQL Request + add indexes
	rows, err := s.pool.Query(ctx, "SELECT collections.*, "+
		"(SELECT sum(materials.xp) "+
		"FROM materials WHERE materials.id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		") AS sum_xp, "+
		"( "+
		"SELECT (SELECT sum(materials.xp) FROM materials WHERE materials.id = user_materials.material_id) "+
		"FROM user_materials WHERE user_materials.material_id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		"AND user_materials.completed = true  AND user_materials.user_id = $1 "+
		") AS xp "+
		"FROM collections "+
		"INNER JOIN user_collections ON user_collections.collection_id = collections.id "+
		"AND user_collections.user_id = $1", userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var collection models.Collection
		if err = rows.Scan(&collection.Id, &collection.CreatedAt, &collection.UpdatedAt, &collection.UserId, &collection.Name, &collection.Description, &collection.SumXp, &collection.Xp); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return collections, nil
}

func (s *Storage) GetCollection(id string, userID string) (models.Collection, error) {
	const op = "storage.postgresql.GetCollection"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var collection models.Collection

	//TODO Need optimize SQL Request + add indexes
	err := s.pool.QueryRow(ctx, "SELECT collections.*, "+
		"( "+
		"SELECT sum(materials.xp) "+
		"FROM materials WHERE materials.id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		") AS sum_xp, "+
		"( "+
		"SELECT (SELECT sum(materials.xp) FROM materials WHERE materials.id = user_materials.material_id) "+
		"FROM user_materials WHERE user_materials.material_id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		"AND user_materials.completed = true AND user_materials.user_id = $1 "+
		") AS xp "+
		"FROM collections WHERE id = $2", userID, id).Scan(&collection.Id, &collection.CreatedAt, &collection.UpdatedAt, &collection.UserId, &collection.Name, &collection.Description, &collection.SumXp, &collection.Xp)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Collection{}, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return models.Collection{}, fmt.Errorf("%s: %w", op, err)
	}

	return collection, nil
}

func (s *Storage) AddCollectionToUser(userID, collectionID string) error {
	const op = "storage.postgresql.AddCollectionToUser"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO user_collections(user_id, collection_id)
        VALUES($1, $2)
        ON CONFLICT DO NOTHING;
    `, userID, collectionID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteCollectionFromUser(userID, collectionID string) error {
	const op = "storage.postgresql.DeleteCollectionFromUser"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        DELETE FROM user_collections
        WHERE user_id=$1 AND collection_id=$2;
    `, userID, collectionID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteCollection(userID, collectionID string) error {
	const op = "storage.postgresql.DeleteCollection"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.pool.Exec(ctx, `
        DELETE FROM collections
        WHERE user_id=$1 AND id=$2;
    `, userID, collectionID)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}

	return nil
}

func (s *Storage) SearchCollections(query string, userID string) ([]models.Collection, error) {
	const op = "storage.postgresql.SearchCollections"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//TODO Need optimize SQL Request + add indexes
	rows, err := s.pool.Query(ctx, "SELECT collections.*, "+
		"( "+
		"SELECT sum(materials.xp) "+
		"FROM materials WHERE materials.id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		") AS sum_xp, "+
		"( "+
		"SELECT (SELECT sum(materials.xp) FROM materials WHERE materials.id = user_materials.material_id) "+
		"FROM user_materials WHERE user_materials.material_id IN (SELECT collection_materials.material_id FROM collection_materials WHERE collection_materials.collection_id = collections.id) "+
		"AND user_materials.completed = true AND user_materials.user_id = $1 "+
		") AS xp "+
		"FROM collections "+
		" WHERE name LIKE '%'||$2||'%' OR description LIKE '%'||$2||'%'", userID, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var collection models.Collection
		if err = rows.Scan(&collection.Id, &collection.CreatedAt, &collection.UpdatedAt, &collection.UserId, &collection.Name, &collection.Description, &collection.SumXp, &collection.Xp); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return collections, nil
}
