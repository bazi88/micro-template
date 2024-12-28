package database

import (
	"fmt"
	"time"
)

// Migration định nghĩa interface cho database migrations
type Migration interface {
	Up() error   // Thực hiện migration
	Down() error // Rollback migration
	Version() string
	Description() string
}

// MigrationHistory lưu trữ lịch sử các migrations đã chạy
type MigrationHistory struct {
	Version     string    `json:"version" db:"version"`
	Description string    `json:"description" db:"description"`
	AppliedAt   time.Time `json:"applied_at" db:"applied_at"`
	Status      string    `json:"status" db:"status"`
}

// Migrator quản lý việc thực thi migrations
type Migrator struct {
	db         DataStore
	migrations []Migration
}

// NewMigrator tạo instance mới của Migrator
func NewMigrator(db DataStore) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// AddMigration thêm migration mới vào danh sách
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// RunMigrations thực hiện tất cả migrations chưa được áp dụng
func (m *Migrator) RunMigrations() error {
	// Tạo bảng migration history nếu chưa tồn tại
	err := m.createMigrationTable()
	if err != nil {
		return fmt.Errorf("failed to create migration table: %v", err)
	}

	// Lấy danh sách migrations đã chạy
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Thực hiện các migrations chưa chạy
	for _, migration := range m.migrations {
		if !m.isApplied(migration.Version(), applied) {
			err := m.runMigration(migration)
			if err != nil {
				return fmt.Errorf("failed to run migration %s: %v", migration.Version(), err)
			}
		}
	}

	return nil
}

// RollbackMigration rollback migration cuối cùng
func (m *Migrator) RollbackMigration() error {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	if len(applied) == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Lấy migration cuối cùng
	lastApplied := applied[len(applied)-1]

	// Tìm migration tương ứng
	var migration Migration
	for _, m := range m.migrations {
		if m.Version() == lastApplied.Version {
			migration = m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", lastApplied.Version)
	}

	// Thực hiện rollback
	err = migration.Down()
	if err != nil {
		return fmt.Errorf("failed to rollback migration %s: %v", migration.Version(), err)
	}

	// Cập nhật trạng thái
	err = m.updateMigrationStatus(migration.Version(), "rolled_back")
	if err != nil {
		return fmt.Errorf("failed to update migration status: %v", err)
	}

	return nil
}

// Helper functions
func (m *Migrator) createMigrationTable() error {
	// Implementation depends on database type
	return nil
}

func (m *Migrator) getAppliedMigrations() ([]MigrationHistory, error) {
	// Implementation depends on database type
	return nil, nil
}

func (m *Migrator) isApplied(version string, applied []MigrationHistory) bool {
	for _, m := range applied {
		if m.Version == version {
			return true
		}
	}
	return false
}

func (m *Migrator) runMigration(migration Migration) error {
	err := migration.Up()
	if err != nil {
		return err
	}

	// Lưu thông tin migration đã chạy
	history := MigrationHistory{
		Version:     migration.Version(),
		Description: migration.Description(),
		AppliedAt:   time.Now(),
		Status:      "applied",
	}

	// Lưu vào database
	// Implementation depends on database type

	return nil
}

func (m *Migrator) updateMigrationStatus(version, status string) error {
	// Implementation depends on database type
	return nil
}
