package database

import (
	"fmt"
	"time"
)

// Seeder định nghĩa interface cho database seeding
type Seeder interface {
	Seed() error  // Thêm dữ liệu mẫu
	Clean() error // Xóa dữ liệu mẫu
	Name() string
}

// SeederHistory lưu trữ lịch sử các seeder đã chạy
type SeederHistory struct {
	Name      string    `json:"name" db:"name"`
	RunAt     time.Time `json:"run_at" db:"run_at"`
	Status    string    `json:"status" db:"status"`
	BatchNo   int       `json:"batch_no" db:"batch_no"`
	CleanedAt time.Time `json:"cleaned_at,omitempty" db:"cleaned_at"`
}

// SeederManager quản lý việc thực thi seeders
type SeederManager struct {
	db      DataStore
	seeders []Seeder
}

// NewSeederManager tạo instance mới của SeederManager
func NewSeederManager(db DataStore) *SeederManager {
	return &SeederManager{
		db:      db,
		seeders: make([]Seeder, 0),
	}
}

// AddSeeder thêm seeder mới vào danh sách
func (sm *SeederManager) AddSeeder(seeder Seeder) {
	sm.seeders = append(sm.seeders, seeder)
}

// RunSeeders thực hiện tất cả seeders chưa được chạy
func (sm *SeederManager) RunSeeders() error {
	// Tạo bảng seeder history nếu chưa tồn tại
	err := sm.createSeederTable()
	if err != nil {
		return fmt.Errorf("failed to create seeder table: %v", err)
	}

	// Lấy batch number mới
	batchNo, err := sm.getNextBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get next batch number: %v", err)
	}

	// Thực hiện các seeders
	for _, seeder := range sm.seeders {
		if !sm.isSeeded(seeder.Name()) {
			err := sm.runSeeder(seeder, batchNo)
			if err != nil {
				return fmt.Errorf("failed to run seeder %s: %v", seeder.Name(), err)
			}
		}
	}

	return nil
}

// CleanSeeders xóa dữ liệu từ batch cuối cùng
func (sm *SeederManager) CleanSeeders() error {
	// Lấy batch number cuối cùng
	lastBatch, err := sm.getLastBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get last batch number: %v", err)
	}

	// Lấy danh sách seeders của batch cuối
	seeders, err := sm.getSeedersInBatch(lastBatch)
	if err != nil {
		return fmt.Errorf("failed to get seeders in batch %d: %v", lastBatch, err)
	}

	// Thực hiện clean theo thứ tự ngược lại
	for i := len(seeders) - 1; i >= 0; i-- {
		seederName := seeders[i].Name
		for _, s := range sm.seeders {
			if s.Name() == seederName {
				err := s.Clean()
				if err != nil {
					return fmt.Errorf("failed to clean seeder %s: %v", seederName, err)
				}
				err = sm.updateSeederStatus(seederName, "cleaned", lastBatch)
				if err != nil {
					return fmt.Errorf("failed to update seeder status: %v", err)
				}
				break
			}
		}
	}

	return nil
}

// Helper functions
func (sm *SeederManager) createSeederTable() error {
	// Implementation depends on database type
	return nil
}

func (sm *SeederManager) getNextBatchNumber() (int, error) {
	// Implementation depends on database type
	return 1, nil
}

func (sm *SeederManager) getLastBatchNumber() (int, error) {
	// Implementation depends on database type
	return 1, nil
}

func (sm *SeederManager) isSeeded(name string) bool {
	// Implementation depends on database type
	return false
}

func (sm *SeederManager) runSeeder(seeder Seeder, batchNo int) error {
	err := seeder.Seed()
	if err != nil {
		return err
	}

	// Lưu thông tin seeder đã chạy
	history := SeederHistory{
		Name:    seeder.Name(),
		RunAt:   time.Now(),
		Status:  "seeded",
		BatchNo: batchNo,
	}

	// Lưu vào database
	// Implementation depends on database type

	return nil
}

func (sm *SeederManager) getSeedersInBatch(batchNo int) ([]SeederHistory, error) {
	// Implementation depends on database type
	return nil, nil
}

func (sm *SeederManager) updateSeederStatus(name, status string, batchNo int) error {
	// Implementation depends on database type
	return nil
}
