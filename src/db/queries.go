package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database connection pool
var pool *pgxpool.Pool 

func init() {
    var err error
    pool, err = pgxpool.Connect(context.Background(), createSQLPool())
    if err != nil {
        fmt.Println("Error connecting to database:", err)
    }
}

// 1. insert 
func Insert(table string, columns []string, values []any) error {
    builder := sq.Insert(table).Columns(columns...)

    for _, val := range values {
        builder = builder.Values(val)
    }

    query, args, err := builder.ToSql()
    if err != nil {
        return err
    }

    _, err = pool.Exec(context.Background(), query, args...)
    return err 
}

// // 2. insertManyOrUpdate
// func insertManyOrUpdate(
//     table string,
//     columns []string,
//     values [][]any,  
//     conflictColumns []string,
//     conflictAction string,
//     returns []string,
// ) (*pgx.Rows, error) {

//     builder := sq.Insert(table).Columns(columns...)

//     // Construct VALUES with unnest, assuming types align with values
//     for _, row := range values {
//         builder = builder.Values(row...)
//     }

//     builder = builder.Suffix(fmt.Sprintf(
//           "ON CONFLICT (%s) DO UPDATE SET %s",
//           strings.Join(conflictColumns, ","),
//           conflictAction,
//     ))

// 	if len(returns) != 0 {
// 		builder = builder.Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returns, ",")))
// 	}

// 	query, args, err := builder.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	x, err := pool.Query(context.Background(), query, args...) 

// 	return x, err
// }

// 3. insertMany
func InsertMany(
    table string,
    columns []string,
    values [][]any, 
) error {
    builder := sq.Insert(table).Columns(columns...)

    // Construct VALUES with unnest
    for _, row := range values {
        builder = builder.Values(row...)
    }

    query, args, err := builder.ToSql()
    if err != nil {
        return err
    }

    _, err = pool.Exec(context.Background(), query, args...)
    return err
}

// func query(query string, data ...any) (*pgx.Rows, error) {
//     return pool.Query(context.Background(), query, data...)
// }