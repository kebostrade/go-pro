package creational

import (
	"fmt"
	"strings"
)

/*
BUILDER PATTERN

Purpose: Separate the construction of a complex object from its representation.

Use Cases:
- Building complex HTTP requests
- Creating database queries
- Constructing configuration objects
- Building UI components

Go-Specific Implementation:
- Method chaining with pointer receivers
- Fluent interface
- Build() method returns the final object
*/

// HTTPRequest represents a complex HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Timeout int
	Retries int
}

// HTTPRequestBuilder builds HTTPRequest objects
type HTTPRequestBuilder struct {
	request *HTTPRequest
}

// NewHTTPRequestBuilder creates a new builder
func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		request: &HTTPRequest{
			Headers: make(map[string]string),
			Method:  "GET",
			Timeout: 30,
			Retries: 3,
		},
	}
}

// Method sets the HTTP method
func (b *HTTPRequestBuilder) Method(method string) *HTTPRequestBuilder {
	b.request.Method = method
	return b
}

// URL sets the request URL
func (b *HTTPRequestBuilder) URL(url string) *HTTPRequestBuilder {
	b.request.URL = url
	return b
}

// Header adds a header
func (b *HTTPRequestBuilder) Header(key, value string) *HTTPRequestBuilder {
	b.request.Headers[key] = value
	return b
}

// Body sets the request body
func (b *HTTPRequestBuilder) Body(body string) *HTTPRequestBuilder {
	b.request.Body = body
	return b
}

// Timeout sets the timeout in seconds
func (b *HTTPRequestBuilder) Timeout(seconds int) *HTTPRequestBuilder {
	b.request.Timeout = seconds
	return b
}

// Retries sets the number of retries
func (b *HTTPRequestBuilder) Retries(retries int) *HTTPRequestBuilder {
	b.request.Retries = retries
	return b
}

// Build creates the final HTTPRequest
func (b *HTTPRequestBuilder) Build() *HTTPRequest {
	return b.request
}

// String returns a string representation
func (r *HTTPRequest) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n", r.Method, r.URL))
	for k, v := range r.Headers {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	if r.Body != "" {
		sb.WriteString(fmt.Sprintf("\n%s", r.Body))
	}
	return sb.String()
}

// SQLQuery represents a complex SQL query
type SQLQuery struct {
	table      string
	columns    []string
	where      []string
	orderBy    string
	limit      int
	offset     int
	joins      []string
}

// SQLQueryBuilder builds SQL queries
type SQLQueryBuilder struct {
	query *SQLQuery
}

// NewSQLQueryBuilder creates a new SQL query builder
func NewSQLQueryBuilder() *SQLQueryBuilder {
	return &SQLQueryBuilder{
		query: &SQLQuery{
			columns: []string{"*"},
		},
	}
}

// Select sets the columns to select
func (b *SQLQueryBuilder) Select(columns ...string) *SQLQueryBuilder {
	b.query.columns = columns
	return b
}

// From sets the table name
func (b *SQLQueryBuilder) From(table string) *SQLQueryBuilder {
	b.query.table = table
	return b
}

// Where adds a WHERE condition
func (b *SQLQueryBuilder) Where(condition string) *SQLQueryBuilder {
	b.query.where = append(b.query.where, condition)
	return b
}

// OrderBy sets the ORDER BY clause
func (b *SQLQueryBuilder) OrderBy(column string) *SQLQueryBuilder {
	b.query.orderBy = column
	return b
}

// Limit sets the LIMIT
func (b *SQLQueryBuilder) Limit(limit int) *SQLQueryBuilder {
	b.query.limit = limit
	return b
}

// Offset sets the OFFSET
func (b *SQLQueryBuilder) Offset(offset int) *SQLQueryBuilder {
	b.query.offset = offset
	return b
}

// Join adds a JOIN clause
func (b *SQLQueryBuilder) Join(join string) *SQLQueryBuilder {
	b.query.joins = append(b.query.joins, join)
	return b
}

// Build creates the final SQL query string
func (b *SQLQueryBuilder) Build() string {
	var sb strings.Builder
	
	// SELECT
	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(b.query.columns, ", "))
	
	// FROM
	if b.query.table != "" {
		sb.WriteString(" FROM ")
		sb.WriteString(b.query.table)
	}
	
	// JOINs
	for _, join := range b.query.joins {
		sb.WriteString(" ")
		sb.WriteString(join)
	}
	
	// WHERE
	if len(b.query.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.query.where, " AND "))
	}
	
	// ORDER BY
	if b.query.orderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(b.query.orderBy)
	}
	
	// LIMIT
	if b.query.limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", b.query.limit))
	}
	
	// OFFSET
	if b.query.offset > 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", b.query.offset))
	}
	
	return sb.String()
}

// User represents a complex user object
type User struct {
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Address   string
	Phone     string
	IsActive  bool
	Roles     []string
}

// UserBuilder builds User objects
type UserBuilder struct {
	user *User
}

// NewUserBuilder creates a new user builder
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &User{
			IsActive: true,
			Roles:    []string{"user"},
		},
	}
}

// ID sets the user ID
func (b *UserBuilder) ID(id string) *UserBuilder {
	b.user.ID = id
	return b
}

// Username sets the username
func (b *UserBuilder) Username(username string) *UserBuilder {
	b.user.Username = username
	return b
}

// Email sets the email
func (b *UserBuilder) Email(email string) *UserBuilder {
	b.user.Email = email
	return b
}

// Name sets first and last name
func (b *UserBuilder) Name(firstName, lastName string) *UserBuilder {
	b.user.FirstName = firstName
	b.user.LastName = lastName
	return b
}

// Age sets the age
func (b *UserBuilder) Age(age int) *UserBuilder {
	b.user.Age = age
	return b
}

// Address sets the address
func (b *UserBuilder) Address(address string) *UserBuilder {
	b.user.Address = address
	return b
}

// Phone sets the phone number
func (b *UserBuilder) Phone(phone string) *UserBuilder {
	b.user.Phone = phone
	return b
}

// AddRole adds a role
func (b *UserBuilder) AddRole(role string) *UserBuilder {
	b.user.Roles = append(b.user.Roles, role)
	return b
}

// Deactivate sets the user as inactive
func (b *UserBuilder) Deactivate() *UserBuilder {
	b.user.IsActive = false
	return b
}

// Build creates the final User
func (b *UserBuilder) Build() *User {
	return b.user
}

