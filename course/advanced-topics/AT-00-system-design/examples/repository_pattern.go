package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type UserRepository interface {
	FindByID(id string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type User struct {
	ID    string
	Name  string
	Email string
}

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) FindByID(id string) (*User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"

	var user User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepo) Create(user *User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"

	_, err := r.db.Exec(query, user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (r *PostgresUserRepo) Update(user *User) error {
	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"

	result, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresUserRepo) Delete(id string) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

type MemoryUserRepo struct {
	users map[string]*User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users: make(map[string]*User),
	}
}

func (r *MemoryUserRepo) FindByID(id string) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *MemoryUserRepo) Create(user *User) error {
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) Update(user *User) error {
	if _, ok := r.users[user.ID]; !ok {
		return ErrNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) Delete(id string) error {
	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}

var ErrNotFound = errors.New("not found")

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email required")
	}

	user := &User{
		ID:    generateID(),
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func generateID() string {
	return "user-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}

func main() {
	repo := NewMemoryUserRepo()
	service := NewUserService(repo)

	user, err := service.CreateUser("John Doe", "john@example.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Created user: %+v\n", user)

	found, err := service.GetUser(user.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Found user: %+v\n", found)
}
