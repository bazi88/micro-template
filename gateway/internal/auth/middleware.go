package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtConfig     JWTConfig
	apiConfig     APIKeyConfig
	permissionMgr *PermissionManager
	excludedPaths []string
}

func NewAuthMiddleware(jwtCfg JWTConfig, apiCfg APIKeyConfig, permMgr *PermissionManager, excluded []string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtConfig:     jwtCfg,
		apiConfig:     apiCfg,
		permissionMgr: permMgr,
		excludedPaths: excluded,
	}
}

func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
	path := c.Path()

	// Kiểm tra excluded paths
	if m.isExcludedPath(path) {
		return c.Next()
	}

	// Xác thực JWT
	claims, err := m.validateAuth(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Parse resource và action từ request
	resource, action, targetID := m.parseRequest(c)

	// Kiểm tra permission
	if !m.permissionMgr.HasPermission(claims, resource, action, targetID) {
		return fiber.NewError(fiber.StatusForbidden, "Permission denied")
	}

	c.Locals("user", claims)
	return c.Next()
}

func (m *AuthMiddleware) isExcludedPath(path string) bool {
	for _, excluded := range m.excludedPaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

func (m *AuthMiddleware) validateAuth(c *fiber.Ctx) (*JWTClaims, error) {
	// Thử xác thực bằng API key trước
	apiKey := c.Get(m.apiConfig.Header)
	if apiKey != "" {
		if err := ValidateAPIKey(c, m.apiConfig); err == nil {
			return &JWTClaims{
				UserID: "api",
				Roles:  []string{"api"},
			}, nil
		}
	}

	// Xác thực JWT
	auth := c.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid authorization header")
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	return ValidateToken(token, m.jwtConfig)
}

func (m *AuthMiddleware) parseRequest(c *fiber.Ctx) (string, string, string) {
	path := c.Path()
	method := c.Method()

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "", "", ""
	}

	resource := parts[1]
	targetID := ""
	if len(parts) > 2 {
		targetID = parts[2]
	}

	action := m.mapMethodToAction(method, parts)

	return resource, action, targetID
}

func (m *AuthMiddleware) mapMethodToAction(method string, pathParts []string) string {
	defaultActions := map[string]string{
		"GET":    "read",
		"POST":   "create",
		"PUT":    "update",
		"DELETE": "delete",
	}

	if len(pathParts) > 3 {
		specialAction := pathParts[3]
		return specialAction
	}

	return defaultActions[method]
}
