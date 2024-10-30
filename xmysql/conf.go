package xmysql

type DBConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	DbCharset    string
}
