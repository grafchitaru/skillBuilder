package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("%s: unable to parse config: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to connect: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.pool.Ping(ctx)
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) GetUser(login string) (string, error) {
	const op = "storage.postgresql.GetUser"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string
	err := s.pool.QueryRow(ctx, "SELECT id FROM users WHERE login = $1", login).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserPassword(login string) (string, error) {
	const op = "storage.postgresql.GetUserPassword"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var password string
	err := s.pool.QueryRow(ctx, "SELECT password FROM users WHERE login = $1", login).Scan(&password)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("%s: operation timed out: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

func (s *Storage) Registration(id string, login string, password string) (string, error) {
	const op = "storage.postgresql.Registration"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("%s begin: %w", op, err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	_, err = tx.Exec(ctx, `
        INSERT INTO users(id, login, password, created_at, updated_at)   
        VALUES($1, $2, $3, $4, $5);
    `, id, login, password, now.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))
	if err != nil {
		return "", fmt.Errorf("%s exec: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("%s commit: %w", op, err)
	}

	return id, nil
}

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

func (s *Storage) GetMaterials(collectionID string) ([]models.Material, error) {
	const op = "storage.postgresql.GetMaterials"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT materials.* "+
		"FROM materials "+
		"INNER JOIN collection_materials ON materials.id = collection_materials.material_id "+
		"WHERE collection_materials.collection_id = $1 ", collectionID)
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

func (s *Storage) GetTypeMaterials() ([]models.TypeMaterial, error) {
	const op = "storage.postgresql.GetTypeMaterials"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT * FROM type_materials")
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
