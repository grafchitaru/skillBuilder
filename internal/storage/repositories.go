package storage

import "github.com/grafchitaru/skillBuilder/internal/models"

type Repositories interface {
	Ping() error
	Close()
	GetUser(login string) (string, error)
	GetUserPassword(login string) (string, error)
	Registration(id string, login string, password string) (string, error)
	CreateCollection(userID string, name string, description string) (string, error)
	CreateMaterial(userID string, name string, description string, typed string, xp int, link string) (string, error)
	DeleteCollection(userID, collectionID string) error
	UpdateCollection(collection models.Collection) error
	AddMaterialToCollection(collectionID, materialID string) error
	UpdateMaterial(material models.Material) error
	DeleteMaterial(userID, materialID string) error
	GetCollections() ([]models.Collection, error)
	GetUserCollections(userID string) ([]models.Collection, error)
	GetCollection(collectionID string) (models.Collection, error)
	GetMaterial(materialID string) (models.Material, error)
	GetMaterials(collectionID string) ([]models.Material, error)
	AddCollectionToUser(userID, collectionID string) error
	DeleteCollectionFromUser(userID, collectionID string) error
	MarkMaterialAsCompleted(userID, materialID string) error
	MarkMaterialAsNotCompleted(userID, materialID string) error
	SearchMaterials(query string) ([]models.Material, error)
	SearchCollections(query string) ([]models.Collection, error)
}
