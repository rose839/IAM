package mysql

import "gorm.io/gorm"

type users struct {
	db *gorm.DB
}

func newUsers(ds *dataStore) *users {
	return &users{db: ds.db}
}
