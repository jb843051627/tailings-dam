package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"tailings-dam/internal/model"

	_ "modernc.org/sqlite"
)

// Store 数据存储层，包含数据库连接和缓存
type Store struct {
	db *sql.DB
	mu sync.RWMutex

	// 缓存映射
	damCache             map[int64]*model.Dam
	monitoringPointCache map[int64]*model.MonitoringPoint
	alertCache           map[int64]*model.Alert
	alertByDamCache      []*model.Alert
	alertByStatusCache   []*model.Alert
	inspectionCache      map[int64]*model.Inspection
	drainageCache        map[int64]*model.DrainageSystem
}

// New 创建新的 Store 实例
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set journal mode: %v", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %v", err)
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy timeout: %v", err)
	}

	s := &Store{
		db:                   db,
		damCache:             make(map[int64]*model.Dam),
		monitoringPointCache: make(map[int64]*model.MonitoringPoint),
		alertCache:           make(map[int64]*model.Alert),
		alertByDamCache:      nil,
		alertByStatusCache:   nil,
		inspectionCache:     make(map[int64]*model.Inspection),
		drainageCache:        make(map[int64]*model.DrainageSystem),
	}

	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init schema: %v", err)
	}

	return s, nil
}

// Close 关闭数据库连接
func (s *Store) Close() error {
	return s.db.Close()
}

// Ping 检查数据库连接
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// initSchema 初始化数据库表结构
func (s *Store) initSchema() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS dams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			location TEXT NOT NULL,
			province TEXT,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			capacity REAL NOT NULL,
			current_level REAL DEFAULT 0,
			hazard_level TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			description TEXT,
			operator TEXT,
			constructed_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS monitoring_points (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dam_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			code TEXT NOT NULL,
			type TEXT NOT NULL,
			latitude REAL,
			longitude REAL,
			elevation REAL,
			status TEXT NOT NULL DEFAULT 'active',
			description TEXT,
			last_reading DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (dam_id) REFERENCES dams(id)
		)`,
		`CREATE TABLE IF NOT EXISTS seepage_readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			point_id INTEGER NOT NULL,
			flow_rate REAL NOT NULL,
			turbidity REAL,
			ph REAL,
			temperature REAL,
			recorded_at DATETIME,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (point_id) REFERENCES monitoring_points(id)
		)`,
		`CREATE TABLE IF NOT EXISTS displacement_readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			point_id INTEGER NOT NULL,
			horizontal_displacement REAL NOT NULL,
			vertical_displacement REAL NOT NULL,
			cumulative_horizontal REAL,
			cumulative_vertical REAL,
			recorded_at DATETIME,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (point_id) REFERENCES monitoring_points(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pore_pressure_readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			point_id INTEGER NOT NULL,
			pressure REAL NOT NULL,
			depth REAL,
			water_level REAL,
			recorded_at DATETIME,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (point_id) REFERENCES monitoring_points(id)
		)`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dam_id INTEGER NOT NULL,
			point_id INTEGER,
			level TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			title TEXT,
			message TEXT NOT NULL,
			threshold REAL,
			current_value REAL,
			reading_type TEXT,
			acknowledged_by TEXT,
			acknowledged_at DATETIME,
			resolved_by TEXT,
			resolved_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (dam_id) REFERENCES dams(id)
		)`,
		`CREATE TABLE IF NOT EXISTS inspections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dam_id INTEGER NOT NULL,
			inspector TEXT NOT NULL,
			title TEXT NOT NULL,
			scheduled_date DATETIME NOT NULL,
			completed_date DATETIME,
			findings TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			priority TEXT NOT NULL DEFAULT 'normal',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (dam_id) REFERENCES dams(id)
		)`,
		`CREATE TABLE IF NOT EXISTS drainage_systems (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dam_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'normal',
			design_flow REAL DEFAULT 0,
			actual_flow REAL DEFAULT 0,
			diameter REAL,
			length REAL,
			material TEXT,
			last_inspection DATETIME,
			next_inspection DATETIME,
			notes TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (dam_id) REFERENCES dams(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mp_dam_id ON monitoring_points(dam_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sr_point_id ON seepage_readings(point_id)`,
		`CREATE INDEX IF NOT EXISTS idx_dr_point_id ON displacement_readings(point_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ppr_point_id ON pore_pressure_readings(point_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_dam_id ON alerts(dam_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,
		`CREATE INDEX IF NOT EXISTS idx_insp_dam_id ON inspections(dam_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ds_dam_id ON drainage_systems(dam_id)`,
	}

	for _, stmt := range statements {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema: %v", err)
		}
	}
	return nil
}

// nullTime 将 sql.NullTime 转换为 time.Time
func nullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

// nullString 将 sql.NullString 转换为 string
func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// nullableTime 将 time.Time 转换为可用于 SQL 的接口值（零值返回 nil）
func nullableTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}

// deepCopyAlerts returns a deep copy of the given alert slice
func deepCopyAlerts(src []*model.Alert) []*model.Alert {
	if src == nil {
		return nil
	}
	dst := make([]*model.Alert, len(src))
	for i, a := range src {
		cp := *a
		dst[i] = &cp
	}
	return dst
}

// clearCache 清除所有缓存
func (s *Store) clearCache() {
	s.mu.Lock()
	s.damCache = make(map[int64]*model.Dam)
	s.monitoringPointCache = make(map[int64]*model.MonitoringPoint)
	s.alertCache = make(map[int64]*model.Alert)
	s.alertByDamCache = nil
	s.alertByStatusCache = nil
	s.inspectionCache = make(map[int64]*model.Inspection)
	s.drainageCache = make(map[int64]*model.DrainageSystem)
	s.mu.Unlock()
}
