package xmysql

import "gorm.io/gorm"

type DBUtil interface {
	CreateDB()
	DropDB()
	GetUtilDB() *gorm.DB
	Close() error
}

type DB interface {
	GetDB() *gorm.DB
	ClearAllData()
	Close() error
}
