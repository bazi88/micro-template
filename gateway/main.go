package main

import (
	"log"
	"time"

	"gateway/internal/auth"
	"gateway/internal/gateway"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Khởi tạo config từ file và environment
	config := gateway.NewConfig()

	// Khởi tạo permission manager với cache TTL từ config
	permissionMgr := auth.NewPermissionManager(
		time.Duration(config.Auth.Permissions.CacheTTL) * time.Second,
	)

	// Đăng ký các services với permission manager
	for serviceName, serviceConfig := range config.ServiceRegistry {
		permissionMgr.RegisterService(serviceName, serviceConfig.Prefixes)
	}

	// Khởi tạo auth middleware
	authMiddleware := auth.NewAuthMiddleware(
		auth.JWTConfig{
			Secret:     config.Auth.JWTSecret,
			Expiration: time.Duration(config.Auth.JWTExpiration) * time.Hour,
		},
		auth.APIKeyConfig{
			Header: config.Auth.APIKeyHeader,
			Secret: config.Auth.APIKeySecret,
		},
		permissionMgr,
		config.Auth.ExcludedPaths,
	)

	// Khởi tạo Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Thêm CORS middleware
	gateway.ConfigureCORS(app, config.CORS)

	// Thêm auth middleware nếu được enable
	if config.Auth.Enabled {
		app.Use(authMiddleware.Handle)
	}

	// Thêm các route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Thêm route để lấy danh sách services và permissions
	app.Get("/api/services", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"services": permissionMgr.GetRegisteredServices(),
		})
	})

	app.Get("/api/services/:service/permissions", func(c *fiber.Ctx) error {
		serviceName := c.Params("service")
		actions, err := permissionMgr.GetResourceActions(serviceName)
		if err != nil {
			return err
		}
		return c.JSON(actions)
	})

	// API cho role management
	app.Post("/api/roles", func(c *fiber.Ctx) error {
		var role auth.Role
		if err := c.BodyParser(&role); err != nil {
			return err
		}
		role.CreatedAt = time.Now()
		role.UpdatedAt = time.Now()
		permissionMgr.AddRole(role)
		return c.JSON(role)
	})

	app.Get("/api/roles/:id", func(c *fiber.Ctx) error {
		roleID := c.Params("id")
		role, exists := permissionMgr.GetRole(roleID)
		if !exists {
			return fiber.NewError(fiber.StatusNotFound, "Role not found")
		}
		return c.JSON(role)
	})

	// API cho tenant management
	app.Post("/api/tenants", func(c *fiber.Ctx) error {
		var tenant auth.Tenant
		if err := c.BodyParser(&tenant); err != nil {
			return err
		}
		tenant.CreatedAt = time.Now()
		tenant.UpdatedAt = time.Now()
		permissionMgr.AddTenant(tenant)
		return c.JSON(tenant)
	})

	// Khởi động server
	log.Fatal(app.Listen(":" + config.Gateway.Port))
}
