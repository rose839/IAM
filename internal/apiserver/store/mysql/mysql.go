package mysql

import (
	"github.com/rose839/IAM/internal/apiserver/store"
	"gorm.io/gorm"
)

type dataStore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *dataStore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *dataStore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *dataStore) Policies() store.PolicyStore {
	return newPolicies(ds)
}
