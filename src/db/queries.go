package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

type QueryBuilder struct {
	builder squirrel.StatementBuilderType
	db      *sql.DB // Add this field
}


func NewQueryBuilder(db *sql.DB) *QueryBuilder {
	return &QueryBuilder{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db:      db,
	}
}


// Insert (single row)
func (qb *QueryBuilder) Insert(table string, data map[string]interface{}) error {
	query := qb.builder.Insert(table).SetMap(data)

	sqlStr, args, err := query.ToSql() 
	if err != nil {
		return err
	}

	_, err = qb.db.Exec(sqlStr, args...)
	return err
}

// InsertManyOrUpdate (upsert)
func (qb *QueryBuilder) InsertManyOrUpdate(
	table string,
	columns []string,
	values [][]interface{},
	conflictColumns []string,
	conflictAction string,
	returns []string,
) (*sql.Rows, error) {

	baseQuery := qb.builder.Insert(table).Columns(columns...)

	// Squirrel makes upsert construction more readable
	if len(conflictColumns) > 0 { 
		baseQuery = baseQuery.Suffix(fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", 
			strings.Join(conflictColumns, ","), conflictAction)) 
	}

	query := baseQuery.Options("VALUES (?)", values) // Use VALUES for bulk rows

	if len(returns) > 0 {
		query = query.Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returns, ",")))
	}

	sqlStr, args, err := query.ToSql() 
	if err != nil { 
		return nil, err 
	} 

	return qb.db.Query(sqlStr, args...) 
}

// InsertMany (bulk insert)
func (qb *QueryBuilder) InsertMany(
	table string,
	columns []string,
	values [][]interface{},
) error {
	query := qb.builder.Insert(table).Columns(columns...).Options("VALUES (?)", values) 

	sqlStr, args, err := query.ToSql() 
	if err != nil {
		return err 
	}

	_, err = qb.db.Exec(sqlStr, args...)
	return err
}

// Query (for general queries - Squirrel is less useful here if simply executing)
func (qb *QueryBuilder) Query(query string, data ...interface{}) (*sql.Rows, error) {
	return qb.db.Query(query, data...)
}
