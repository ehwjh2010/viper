package postgresql

import (
	"github.com/ehwjh2010/cobra/client"
	"github.com/ehwjh2010/cobra/db/rdb"
)

//InitPostgresql 初始化Mysql
func InitPostgresql(dbConfig *client.DB) (client *rdb.DBClient, err error) {

	db, err := rdb.InitDBWithGorm(dbConfig, rdb.Postgresql)

	if err != nil {
		return nil, err
	}

	client = rdb.NewDBClient(db, rdb.Postgresql)

	return client, nil
}