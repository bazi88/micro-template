package auth

import (
	"fmt"
	"sync"
	"time"
)

type permissionCacheKey struct {
	userID       string
	resourceName string
	actionName   string
	targetID     string
}

type permissionCacheEntry struct {
	result    bool
	timestamp time.Time
}

type PermissionManager struct {
	resources map[string]Resource
	roles     map[string]Role
	tenants   map[string]Tenant
	cacheTTL  time.Duration
	cache     map[permissionCacheKey]permissionCacheEntry
	mu        sync.RWMutex
}

func NewPermissionManager(cacheTTL time.Duration) *PermissionManager {
	return &PermissionManager{
		resources: make(map[string]Resource),
		roles:     make(map[string]Role),
		tenants:   make(map[string]Tenant),
		cacheTTL:  cacheTTL,
		cache:     make(map[permissionCacheKey]permissionCacheEntry),
	}
}

func (pm *PermissionManager) HasPermission(claims *JWTClaims, resourceName string, actionName string, targetID string) bool {
	pm.mu.RLock()

	// Kiểm tra cache
	cacheKey := permissionCacheKey{
		userID:       claims.UserID,
		resourceName: resourceName,
		actionName:   actionName,
		targetID:     targetID,
	}

	if entry, exists := pm.cache[cacheKey]; exists {
		if time.Since(entry.timestamp) < pm.cacheTTL {
			pm.mu.RUnlock()
			return entry.result
		}
	}
	pm.mu.RUnlock()

	// Nếu không có trong cache hoặc cache đã hết hạn, tính toán lại
	pm.mu.Lock()
	defer pm.mu.Unlock()

	result := pm.checkPermission(claims, resourceName, actionName, targetID)

	// Lưu kết quả vào cache
	pm.cache[cacheKey] = permissionCacheEntry{
		result:    result,
		timestamp: time.Now(),
	}

	return result
}

func (pm *PermissionManager) checkPermission(claims *JWTClaims, resourceName string, actionName string, targetID string) bool {
	// Kiểm tra resource tồn tại
	resource, exists := pm.resources[resourceName]
	if !exists {
		return false
	}

	// Kiểm tra action tồn tại
	_, exists = resource.Actions[actionName]
	if !exists {
		return false
	}

	// Kiểm tra tenant active
	if tenant, exists := pm.tenants[claims.TenantID]; !exists || tenant.Status != "active" {
		return false
	}

	// Kiểm tra từng role của user
	for _, roleID := range claims.Roles {
		role, exists := pm.roles[roleID]
		if !exists || role.TenantID != claims.TenantID {
			continue
		}

		// Kiểm tra permissions của role
		for _, perm := range role.Permissions {
			if perm.ResourceName == resourceName && perm.ActionName == actionName {
				// Kiểm tra scope
				for _, permScope := range perm.Scopes {
					switch permScope {
					case "all":
						return true
					case "tenant":
						return pm.isInSameTenant(targetID, claims.TenantID)
					case "own":
						return targetID == claims.UserID
					}
				}
			}
		}
	}
	return false
}

func (pm *PermissionManager) isInSameTenant(targetID string, tenantID string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	// Kiểm tra trong danh sách users
	for _, role := range pm.roles {
		if role.ID == targetID {
			return role.TenantID == tenantID
		}
	}

	return false
}

func (pm *PermissionManager) GetResourceActions(resourceName string) (map[string]ResourceAction, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	resource, exists := pm.resources[resourceName]
	if !exists {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}
	return resource.Actions, nil
}

// Các method quản lý role
func (pm *PermissionManager) GetRole(roleID string) (Role, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	role, exists := pm.roles[roleID]
	return role, exists
}

func (pm *PermissionManager) AddRole(role Role) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.roles[role.ID] = role
}

func (pm *PermissionManager) UpdateRole(role Role) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.roles[role.ID]; !exists {
		return false
	}
	pm.roles[role.ID] = role
	return true
}

func (pm *PermissionManager) DeleteRole(roleID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.roles[roleID]; !exists {
		return false
	}
	delete(pm.roles, roleID)
	return true
}

// Các method quản lý tenant
func (pm *PermissionManager) GetTenant(tenantID string) (Tenant, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	tenant, exists := pm.tenants[tenantID]
	return tenant, exists
}

func (pm *PermissionManager) AddTenant(tenant Tenant) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.tenants[tenant.ID] = tenant
}

func (pm *PermissionManager) UpdateTenant(tenant Tenant) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.tenants[tenant.ID]; !exists {
		return false
	}
	pm.tenants[tenant.ID] = tenant
	return true
}

func (pm *PermissionManager) DeleteTenant(tenantID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.tenants[tenantID]; !exists {
		return false
	}
	delete(pm.tenants, tenantID)
	return true
}

// Thêm method để đăng ký service mới
func (pm *PermissionManager) RegisterService(serviceName string, prefixes []string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	resource := NewResourceFromService(serviceName, prefixes)
	pm.resources[serviceName] = resource
}

// Thêm method để lấy danh sách services đã đăng ký
func (pm *PermissionManager) GetRegisteredServices() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	services := make([]string, 0, len(pm.resources))
	for name := range pm.resources {
		services = append(services, name)
	}
	return services
}

// Thêm method để xóa service
func (pm *PermissionManager) UnregisterService(serviceName string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.resources, serviceName)
}
