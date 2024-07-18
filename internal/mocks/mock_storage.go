package mocks

import (
	"errors"
)

type GetUserFunc func(login string) (string, error)
type GetUserPasswordFunc func(login string) (string, error)
type RegistrationFunc func(id string, login string, password string) (string, error)
type CreateCollectionFunc func(userID string, name string, description string) (string, error)
type UpdateCollectionFunc func(collectionID string, name string, description string) error
type AddMaterialToCollectionFunc func(collectionID, materialID string) error
type UpdateMaterialFunc func(materialID string, name string, description string, materialType string, link string, xp int) error
type DeleteMaterialFunc func(materialID string) error
type GetCollectionsByServiceFunc func(service string) ([]string, error)
type GetUserCollectionsFunc func(userID string) ([]string, error)
type GetCollectionFunc func(collectionID string) (string, string, error)
type GetMaterialFunc func(materialID string) (string, string, string, string, int, error)
type AddCollectionToUserFunc func(userID, collectionID string) error
type DeleteCollectionFromUserFunc func(userID, collectionID string) error
type MarkMaterialAsCompletedFunc func(userID, materialID string) error
type MarkMaterialAsNotCompletedFunc func(userID, materialID string) error

type MockStorage struct {
	PingError                      error
	GetUserFunc                    GetUserFunc
	RegistrationFunc               RegistrationFunc
	GetUserPasswordFunc            GetUserPasswordFunc
	Users                          map[string]string
	IDs                            map[string]string
	Passwords                      map[string]string
	CreateCollectionFunc           CreateCollectionFunc
	UpdateCollectionFunc           UpdateCollectionFunc
	AddMaterialToCollectionFunc    AddMaterialToCollectionFunc
	UpdateMaterialFunc             UpdateMaterialFunc
	DeleteMaterialFunc             DeleteMaterialFunc
	GetCollectionsByServiceFunc    GetCollectionsByServiceFunc
	GetUserCollectionsFunc         GetUserCollectionsFunc
	GetCollectionFunc              GetCollectionFunc
	GetMaterialFunc                GetMaterialFunc
	AddCollectionToUserFunc        AddCollectionToUserFunc
	DeleteCollectionFromUserFunc   DeleteCollectionFromUserFunc
	MarkMaterialAsCompletedFunc    MarkMaterialAsCompletedFunc
	MarkMaterialAsNotCompletedFunc MarkMaterialAsNotCompletedFunc
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Users:     make(map[string]string),
		IDs:       make(map[string]string),
		Passwords: make(map[string]string),
	}
}

func (ms *MockStorage) Ping() error {
	return ms.PingError
}

func (ms *MockStorage) Close() {
	// Implementation for Close method
}

func (ms *MockStorage) GetUser(login string) (string, error) {
	if ms.GetUserFunc != nil {
		return ms.GetUserFunc(login)
	}
	user, exists := ms.Users[login]
	if !exists {
		return "", errors.New("user not found")
	}
	return user, nil
}

func (ms *MockStorage) GetUserPassword(login string) (string, error) {
	if ms.GetUserPasswordFunc != nil {
		return ms.GetUserPasswordFunc(login)
	}
	password, exists := ms.Passwords[login]
	if !exists {
		return "", errors.New("user not found")
	}
	return password, nil
}

func (ms *MockStorage) SetUserPassword(login string, password string) {
	ms.Passwords[login] = password
}

func (ms *MockStorage) Registration(id string, login string, password string) (string, error) {
	if ms.RegistrationFunc != nil {
		return ms.RegistrationFunc(id, login, password)
	}
	if _, exists := ms.Users[login]; exists {
		return "", errors.New("user already exists")
	}
	ms.Users[login] = password
	ms.IDs[login] = id
	return id, nil
}

func (ms *MockStorage) CreateCollection(userID string, name string, description string) (string, error) {
	if ms.CreateCollectionFunc != nil {
		return ms.CreateCollectionFunc(userID, name, description)
	}
	return "", errors.New("not implemented")
}

func (ms *MockStorage) UpdateCollection(collectionID string, name string, description string) error {
	if ms.UpdateCollectionFunc != nil {
		return ms.UpdateCollectionFunc(collectionID, name, description)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) AddMaterialToCollection(collectionID, materialID string) error {
	if ms.AddMaterialToCollectionFunc != nil {
		return ms.AddMaterialToCollectionFunc(collectionID, materialID)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) UpdateMaterial(materialID string, name string, description string, materialType string, link string, xp int) error {
	if ms.UpdateMaterialFunc != nil {
		return ms.UpdateMaterialFunc(materialID, name, description, materialType, link, xp)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) DeleteMaterial(materialID string) error {
	if ms.DeleteMaterialFunc != nil {
		return ms.DeleteMaterialFunc(materialID)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) GetCollectionsByService(service string) ([]string, error) {
	if ms.GetCollectionsByServiceFunc != nil {
		return ms.GetCollectionsByServiceFunc(service)
	}
	return nil, errors.New("not implemented")
}

func (ms *MockStorage) GetUserCollections(userID string) ([]string, error) {
	if ms.GetUserCollectionsFunc != nil {
		return ms.GetUserCollectionsFunc(userID)
	}
	return nil, errors.New("not implemented")
}

func (ms *MockStorage) GetCollection(collectionID string) (string, string, error) {
	if ms.GetCollectionFunc != nil {
		return ms.GetCollectionFunc(collectionID)
	}
	return "", "", errors.New("not implemented")
}

func (ms *MockStorage) GetMaterial(materialID string) (string, string, string, string, int, error) {
	if ms.GetMaterialFunc != nil {
		return ms.GetMaterialFunc(materialID)
	}
	return "", "", "", "", 0, errors.New("not implemented")
}

func (ms *MockStorage) AddCollectionToUser(userID, collectionID string) error {
	if ms.AddCollectionToUserFunc != nil {
		return ms.AddCollectionToUserFunc(userID, collectionID)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) DeleteCollectionFromUser(userID, collectionID string) error {
	if ms.DeleteCollectionFromUserFunc != nil {
		return ms.DeleteCollectionFromUserFunc(userID, collectionID)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) MarkMaterialAsCompleted(userID, materialID string) error {
	if ms.MarkMaterialAsCompletedFunc != nil {
		return ms.MarkMaterialAsCompletedFunc(userID, materialID)
	}
	return errors.New("not implemented")
}

func (ms *MockStorage) MarkMaterialAsNotCompleted(userID, materialID string) error {
	if ms.MarkMaterialAsNotCompletedFunc != nil {
		return ms.MarkMaterialAsNotCompletedFunc(userID, materialID)
	}
	return errors.New("not implemented")
}
