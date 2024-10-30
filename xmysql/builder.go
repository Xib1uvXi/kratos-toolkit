package xmysql

func MakeDBUtil(dbConfig *DBConfig) DBUtil {
	return newGormMysql(dbConfig, true)
}

func MakeDB(dbConfig *DBConfig) DB {
	return newGormMysql(dbConfig, false)
}
