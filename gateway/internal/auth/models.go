package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResourceAction struct {
	Name        string   `json:"name"`
	Scopes      []string `json:"scopes"`
	Description string   `json:"description"`
}

type Resource struct {
	Name        string                    `json:"name"`
	Actions     map[string]ResourceAction `json:"actions"`
	Prefixes    []string                  `json:"prefixes"`
	Description string                    `json:"description"`
}

type Permission struct {
	ResourceName string   `json:"resource_name"`
	ActionName   string   `json:"action_name"`
	Scopes       []string `json:"scopes"`
}

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	TenantID    string       `json:"tenant_id"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	TenantID  string    `json:"tenant_id"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JWTClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}
