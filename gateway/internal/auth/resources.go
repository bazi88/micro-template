package auth

// Định nghĩa các action cố định
var DefaultActions = map[string]ResourceAction{
	"create": {
		Name:        "create",
		Scopes:      []string{"tenant"},
		Description: "Create new resource",
	},
	"read": {
		Name:        "read",
		Scopes:      []string{"all", "tenant", "own"},
		Description: "Read resource information",
	},
	"update": {
		Name:        "update",
		Scopes:      []string{"tenant", "own"},
		Description: "Update resource information",
	},
	"delete": {
		Name:        "delete",
		Scopes:      []string{"tenant"},
		Description: "Delete resource",
	},
	"list": {
		Name:        "list",
		Scopes:      []string{"all", "tenant"},
		Description: "List resources",
	},
	"publish": {
		Name:        "publish",
		Scopes:      []string{"tenant"},
		Description: "Publish resource",
	},
	"unpublish": {
		Name:        "unpublish",
		Scopes:      []string{"tenant"},
		Description: "Unpublish resource",
	},
}

// Hàm tạo resource từ service
func NewResourceFromService(serviceName string, prefixes []string) Resource {
	actions := make(map[string]ResourceAction)

	// Copy các action mặc định
	for name, action := range DefaultActions {
		actions[name] = action
	}

	// Thêm các action đặc biệt dựa trên service
	switch serviceName {
	case "user-service":
		actions["change_role"] = ResourceAction{
			Name:        "change_role",
			Scopes:      []string{"tenant"},
			Description: "Change user role",
		}
		actions["change_password"] = ResourceAction{
			Name:        "change_password",
			Scopes:      []string{"tenant", "own"},
			Description: "Change user password",
		}
	case "notification-service":
		actions["send"] = ResourceAction{
			Name:        "send",
			Scopes:      []string{"tenant"},
			Description: "Send notification",
		}
		actions["broadcast"] = ResourceAction{
			Name:        "broadcast",
			Scopes:      []string{"tenant"},
			Description: "Broadcast notification to all users",
		}
	case "logging-service":
		actions["export"] = ResourceAction{
			Name:        "export",
			Scopes:      []string{"tenant"},
			Description: "Export logs",
		}
		actions["clear"] = ResourceAction{
			Name:        "clear",
			Scopes:      []string{"tenant"},
			Description: "Clear logs",
		}
	}

	return Resource{
		Name:        serviceName,
		Actions:     actions,
		Prefixes:    prefixes,
		Description: "Resource for " + serviceName,
	}
}
