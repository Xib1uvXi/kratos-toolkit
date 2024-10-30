package xmysql

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

func newGormMysql(dbConfig *DBConfig, forUtil bool) *gormMysql {
	gm := &gormMysql{dbConfig: dbConfig}

	if forUtil {
		gm.initCdDb()
		return gm
	}

	// init db
	gm.initGormDB()

	return gm
}

type gormMysql struct {
	dbConfig *DBConfig
	db       *gorm.DB
	utilDB   *gorm.DB
	sqlDB    *sql.DB
}

// close
func (gm *gormMysql) Close() error {
	if gm.sqlDB != nil {
		err := gm.sqlDB.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (gm *gormMysql) CreateDB() {
	createDbSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8 COLLATE utf8_general_ci;", gm.dbConfig.DbName)

	err := gm.utilDB.Exec(createDbSQL).Error
	if err != nil {
		log.Errorf("create db failed: %v", err)
		return
	}

	log.Infof("%s database create success", gm.dbConfig.DbName)
}

func (gm *gormMysql) DropDB() {
	dropDbSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", gm.dbConfig.DbName)

	err := gm.utilDB.Exec(dropDbSQL).Error
	if err != nil {
		log.Errorf("drop db failed: %v", err)
		return
	}

	log.Infof("%s database drop success", gm.dbConfig.DbName)
}

func (gm *gormMysql) GetUtilDB() *gorm.DB {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	log.Infof("init db connection: %s, name: %s", gm.dbConfig.Host, gm.dbConfig.DbName)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName)

	sqlDB, err := sql.Open("mysql", connStr)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})

	if err != nil {
		panic("failed to open gorm: " + err.Error())
	}

	return db
}

func (gm *gormMysql) GetDB() *gorm.DB {
	return gm.db
}

func (gm *gormMysql) ClearAllData() {
	if flag.Lookup("test.v") != nil && (strings.Contains(gm.dbConfig.DbName, "test") || strings.Contains(gm.dbConfig.DbName, "dev")) {
		tmpDb := gm.db
		if tmpDb == nil {
			panic("db is nil, please init db first")
		}

		if rs, err := tmpDb.Raw("show tables;").Rows(); err == nil {
			var tName string
			for rs.Next() {
				if err := rs.Scan(&tName); err != nil || tName == "" {
					log.Errorf("get table name %s failed: %v", tName, err)
					panic("get table name failed")
				}
				if err := tmpDb.Exec(fmt.Sprintf("delete from %s", tName)).Error; err != nil {
					panic("clear data failed: " + err.Error())
				}
			}
		} else {
			panic("get table list failed：" + err.Error())
		}
	} else {
		panic("only test environment can clear data")
	}
}

func (gm *gormMysql) initGormDB() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	log.Infof("init db connection: %s, name: %s", gm.dbConfig.Host, gm.dbConfig.DbName)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName)

	sqlDB, err := sql.Open("mysql", connStr)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	// silence gorm log
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to open gorm: " + err.Error())
	}

	gm.sqlDB = sqlDB

	gm.db = db
}

func (gm *gormMysql) initCdDb() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, "information_schema", gm.dbConfig.DbCharset)

	sqlDB, err := sql.Open("mysql", cStr)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})

	if err != nil {
		panic("failed to open gorm: " + err.Error())
	}

	gm.utilDB = db
	gm.sqlDB = sqlDB
}
