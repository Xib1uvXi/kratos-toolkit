package xmysql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMakeDBUtil(t *testing.T) {
	t.Skip("Skip TestMakeDBUtil")
	dbConf := &DBConfig{
		Username:     "root",
		Password:     "root",
		Host:         "127.0.0.1",
		Port:         "3306",
		DbName:       "",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		DbCharset:    "utf8mb4",
	}

	dbConf.DbName = "hahaha_test"

	var utilDB DBUtil
	require.NotPanics(t, func() {
		utilDB = MakeDBUtil(dbConf)
	})

	utilDB.CreateDB()
	utilDB.DropDB()
	utilDB.Close()
}

func TestMakeDB(t *testing.T) {
	t.Skip("Skip TestMakeDBUtil")
	dbConf := &DBConfig{
		Username:     "root",
		Password:     "root",
		Host:         "127.0.0.1",
		Port:         "3306",
		DbName:       "",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		DbCharset:    "utf8mb4",
	}

	dbConf.DbName = "hahaha_test"

	var utilDB DBUtil
	require.NotPanics(t, func() {
		utilDB = MakeDBUtil(dbConf)
	})

	utilDB.CreateDB()
	defer utilDB.DropDB()

	require.NotPanics(t, func() {
		db := MakeDB(dbConf)
		db.ClearAllData()
	})
}
