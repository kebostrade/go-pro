// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package postgres provides functionality for the GO-PRO Learning Platform.
package postgres

import (
	"fmt"
	"strings"
)

// QueryBuilder helps build complex SQL queries.
type QueryBuilder struct {
	table      string
	columns    []string
	joins      []string
	where      []string
	whereArgs  []interface{}
	orderBy    []string
	groupBy    []string
	having     []string
	havingArgs []interface{}
	limit      int
	offset     int
	argCounter int
}

// NewQueryBuilder creates a new query builder.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:      table,
		columns:    make([]string, 0),
		joins:      make([]string, 0),
		where:      make([]string, 0),
		whereArgs:  make([]interface{}, 0),
		orderBy:    make([]string, 0),
		groupBy:    make([]string, 0),
		having:     make([]string, 0),
		havingArgs: make([]interface{}, 0),
		argCounter: 1,
	}
}

// Select specifies the columns to select.
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = append(qb.columns, columns...)
	return qb
}

// Join adds a JOIN clause.
func (qb *QueryBuilder) Join(joinType, table, condition string) *QueryBuilder {
	join := fmt.Sprintf("%s JOIN %s ON %s", joinType, table, condition)
	qb.joins = append(qb.joins, join)

	return qb
}

// InnerJoin adds an INNER JOIN clause.
func (qb *QueryBuilder) InnerJoin(table, condition string) *QueryBuilder {
	return qb.Join("INNER", table, condition)
}

// LeftJoin adds a LEFT JOIN clause.
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	return qb.Join("LEFT", table, condition)
}

// RightJoin adds a RIGHT JOIN clause.
func (qb *QueryBuilder) RightJoin(table, condition string) *QueryBuilder {
	return qb.Join("RIGHT", table, condition)
}

// Where adds a WHERE condition.
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	// Replace ? placeholders with $1, $2, etc.
	condition = qb.replacePlaceholders(condition, len(args))
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, args...)

	return qb
}

// WhereIn adds a WHERE IN condition.
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", qb.argCounter)
		qb.argCounter++
	}

	condition := fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, values...)

	return qb
}

// WhereNotIn adds a WHERE NOT IN condition.
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", qb.argCounter)
		qb.argCounter++
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, values...)

	return qb
}

// WhereBetween adds a WHERE BETWEEN condition.
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	condition := fmt.Sprintf("%s BETWEEN $%d AND $%d", column, qb.argCounter, qb.argCounter+1)
	qb.argCounter += 2
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, start, end)

	return qb
}

// WhereLike adds a WHERE LIKE condition.
func (qb *QueryBuilder) WhereLike(column, pattern string) *QueryBuilder {
	condition := fmt.Sprintf("%s LIKE $%d", column, qb.argCounter)
	qb.argCounter++
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, pattern)

	return qb
}

// WhereNull adds a WHERE IS NULL condition.
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	condition := fmt.Sprintf("%s IS NULL", column)
	qb.where = append(qb.where, condition)

	return qb
}

// WhereNotNull adds a WHERE IS NOT NULL condition.
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	condition := fmt.Sprintf("%s IS NOT NULL", column)
	qb.where = append(qb.where, condition)

	return qb
}

// OrderBy adds an ORDER BY clause.
func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
	order := fmt.Sprintf("%s %s", column, strings.ToUpper(direction))
	qb.orderBy = append(qb.orderBy, order)

	return qb
}

// GroupBy adds a GROUP BY clause.
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having adds a HAVING condition.
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	// Replace ? placeholders with $1, $2, etc.
	condition = qb.replacePlaceholders(condition, len(args))
	qb.having = append(qb.having, condition)
	qb.havingArgs = append(qb.havingArgs, args...)

	return qb
}

// Limit sets the LIMIT clause.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets the OFFSET clause.
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Build builds the final SQL query and returns it with arguments.
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var query strings.Builder

	// SELECT clause.
	query.WriteString("SELECT ")
	if len(qb.columns) == 0 {
		query.WriteString("*")
	} else {
		query.WriteString(strings.Join(qb.columns, ", "))
	}

	// FROM clause.
	query.WriteString(fmt.Sprintf(" FROM %s", qb.table))

	// JOIN clauses.
	if len(qb.joins) > 0 {
		query.WriteString(" ")
		query.WriteString(strings.Join(qb.joins, " "))
	}

	// WHERE clause.
	if len(qb.where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.where, " AND "))
	}

	// GROUP BY clause.
	if len(qb.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupBy, ", "))
	}

	// HAVING clause.
	if len(qb.having) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(qb.having, " AND "))
	}

	// ORDER BY clause.
	if len(qb.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(qb.orderBy, ", "))
	}

	// LIMIT clause.
	if qb.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	// OFFSET clause.
	if qb.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	// Combine all arguments.
	args := make([]interface{}, 0, len(qb.whereArgs)+len(qb.havingArgs))
	args = append(args, qb.whereArgs...)
	args = append(args, qb.havingArgs...)

	return query.String(), args
}

// BuildCount builds a COUNT query.
func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var query strings.Builder

	// SELECT COUNT clause.
	query.WriteString("SELECT COUNT(*)")

	// FROM clause.
	query.WriteString(fmt.Sprintf(" FROM %s", qb.table))

	// JOIN clauses.
	if len(qb.joins) > 0 {
		query.WriteString(" ")
		query.WriteString(strings.Join(qb.joins, " "))
	}

	// WHERE clause.
	if len(qb.where) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.where, " AND "))
	}

	// GROUP BY clause (for COUNT with GROUP BY)
	if len(qb.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(qb.groupBy, ", "))
	}

	// HAVING clause.
	if len(qb.having) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(qb.having, " AND "))
	}

	// Combine all arguments.
	args := make([]interface{}, 0, len(qb.whereArgs)+len(qb.havingArgs))
	args = append(args, qb.whereArgs...)
	args = append(args, qb.havingArgs...)

	return query.String(), args
}

// replacePlaceholders replaces ? with $1, $2, etc.
func (qb *QueryBuilder) replacePlaceholders(condition string, argCount int) string {
	result := condition
	for i := 0; i < argCount; i++ {
		result = strings.Replace(result, "?", fmt.Sprintf("$%d", qb.argCounter), 1)
		qb.argCounter++
	}

	return result
}

// Reset resets the query builder to its initial state.
func (qb *QueryBuilder) Reset() *QueryBuilder {
	qb.columns = make([]string, 0)
	qb.joins = make([]string, 0)
	qb.where = make([]string, 0)
	qb.whereArgs = make([]interface{}, 0)
	qb.orderBy = make([]string, 0)
	qb.groupBy = make([]string, 0)
	qb.having = make([]string, 0)
	qb.havingArgs = make([]interface{}, 0)
	qb.limit = 0
	qb.offset = 0
	qb.argCounter = 1

	return qb
}
