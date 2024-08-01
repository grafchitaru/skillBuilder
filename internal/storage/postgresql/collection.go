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

	query := `
	SELECT
		collections.id, collections.user_id, collections.name, collections.description, collections.created_at, collections.updated_at,
		COALESCE(
			(SELECT SUM(materials.xp)
			 FROM materials
			 JOIN collection_materials ON materials.id = collection_materials.material_id
			 WHERE collection_materials.collection_id = collections.id), 0
		) AS sum_xp,
		COALESCE(
			(SELECT SUM(materials.xp)
			 FROM materials
			 JOIN collection_materials ON materials.id = collection_materials.material_id
			 JOIN user_materials ON materials.id = user_materials.material_id
			 WHERE collection_materials.collection_id = collections.id
			   AND user_materials.completed = true
			   AND user_materials.user_id = $1), 0
		) AS xp
	FROM
		collections;
	`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var collection models.Collection
		if err = rows.Scan(&collection.Id, &collection.UserId, &collection.Name, &collection.Description, &collection.CreatedAt, &collection.UpdatedAt, &collection.SumXp, &collection.Xp); err != nil {
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
	rows, err := s.pool.Query(ctx, `SELECT collections.id, collections.user_id, collections.name, collections.description, collections.created_at, collections.updated_at,
       COALESCE(sum_xp.total_xp, 0) AS sum_xp,
       COALESCE(user_xp.total_xp, 0) AS xp
FROM collections
INNER JOIN user_collections ON user_collections.collection_id = collections.id
LEFT JOIN (
    SELECT collection_id, SUM(materials.xp) AS total_xp
    FROM collection_materials
    INNER JOIN materials ON collection_materials.material_id = materials.id
    GROUP BY collection_id
) AS sum_xp ON sum_xp.collection_id = collections.id
LEFT JOIN (
    SELECT collection_materials.collection_id, SUM(materials.xp) AS total_xp
    FROM collection_materials
    INNER JOIN materials ON collection_materials.material_id = materials.id
    INNER JOIN user_materials ON materials.id = user_materials.material_id
    WHERE user_materials.completed = true
      AND user_materials.user_id = $1
    GROUP BY collection_materials.collection_id
) AS user_xp ON user_xp.collection_id = collections.id
WHERE user_collections.user_id = $1`, userID)
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
		if err = rows.Scan(&collection.Id, &collection.UserId, &collection.Name, &collection.Description, &collection.CreatedAt, &collection.UpdatedAt, &collection.SumXp, &collection.Xp); err != nil {
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

	// Validate input parameters
	if _, err := uuid.Parse(id); err != nil {
		return models.Collection{}, fmt.Errorf("%s: invalid collection ID: %w", op, err)
	}
	if _, err := uuid.Parse(userID); err != nil {
		return models.Collection{}, fmt.Errorf("%s: invalid user ID: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var collection models.Collection

	// Log the input parameters for debugging
	fmt.Printf("GetCollection: id=%s, userID=%s\n", id, userID)

	// TODO Need optimize SQL Request + add indexes
	err := s.pool.QueryRow(ctx, `SELECT collections.id, collections.user_id, collections.name, collections.description, collections.created_at, collections.updated_at,
       COALESCE(sum_xp.total_xp, 0) AS sum_xp,
       COALESCE(user_xp.total_xp, 0) AS xp
FROM collections
LEFT JOIN (
    SELECT collection_id, SUM(materials.xp) AS total_xp
    FROM collection_materials
    INNER JOIN materials ON collection_materials.material_id = materials.id
    GROUP BY collection_id
) AS sum_xp ON sum_xp.collection_id = collections.id
LEFT JOIN (
    SELECT collection_materials.collection_id, SUM(materials.xp) AS total_xp
    FROM collection_materials
    INNER JOIN materials ON collection_materials.material_id = materials.id
    INNER JOIN user_materials ON materials.id = user_materials.material_id
    WHERE user_materials.completed = true
      AND user_materials.user_id = $1
    GROUP BY collection_materials.collection_id
) AS user_xp ON user_xp.collection_id = collections.id
WHERE collections.id = $2
`, userID, id).Scan(&collection.Id, &collection.UserId, &collection.Name, &collection.Description, &collection.CreatedAt, &collection.UpdatedAt, &collection.SumXp, &collection.Xp)
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
	rows, err := s.pool.Query(ctx, "SELECT collections.id, collections.user_id, collections.name, collections.description, collections.created_at, collections.updated_at, "+
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
		if err = rows.Scan(&collection.Id, &collection.UserId, &collection.Name, &collection.Description, &collection.CreatedAt, &collection.UpdatedAt, &collection.SumXp, &collection.Xp); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return collections, nil
}
