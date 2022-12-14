// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TasksDao is the data access object for table tasks.
type TasksDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns TasksColumns // columns contains all the column names of Table for convenient usage.
}

// TasksColumns defines and stores column names for table tasks.
type TasksColumns struct {
	Id         string //
	UploadId   string //
	BucketName string //
	FileKey    string //
	Src        string //
	CreatedAt  string //
}

//  tasksColumns holds the columns for table tasks.
var tasksColumns = TasksColumns{
	Id:         "id",
	UploadId:   "upload_id",
	BucketName: "bucket_name",
	FileKey:    "file_key",
	Src:        "src",
	CreatedAt:  "created_at",
}

// NewTasksDao creates and returns a new DAO object for table data access.
func NewTasksDao() *TasksDao {
	return &TasksDao{
		group:   "default",
		table:   "tasks",
		columns: tasksColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TasksDao) Columns() TasksColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TasksDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
