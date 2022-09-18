package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	_ "github.com/mattn/go-sqlite3"
)

func GenerateNewDb(dbPath, dbName string) (err error) {
	pwd := dbPath + "/" + dbName
	os.Remove(pwd)

	db, err := sql.Open("sqlite3", pwd)
	if err != nil {
		return
	}
	defer db.Close()

	sqlStmt := `
	create table tasks (id integer not null primary key, upload_id text not null, bucket_name text not null, file_key text not null, src text not null, created_at text);
	create table parts (id integer not null primary key, task_id integer not null, part_num integer not null, content_md5 text not null, etag text not null, status int2);
	`
	_, err = db.Exec(sqlStmt)
	return
}

func GetDbName(bucketName, fileName string) (dbName string) {
	return bucketName + "_" + fileName + ".db"
}

func InitUploadTaskDbConfig(dbPath string) (err error) {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Type: "sqlite",
				Link: dbPath,
			},
		},
	})
	return
}

func GetDbInstance(configName string) (db gdb.DB, err error) {
	if configName == "" {
		configName = "default"
	}
	return gdb.NewByGroup(configName)
}
