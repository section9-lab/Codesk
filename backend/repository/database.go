package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
	mu   sync.Mutex
)

// InitDatabase 初始化数据库连接
func InitDatabase(dbPath string) error {
	var err error
	once.Do(func() {
		// 确保数据库目录存在
		dbDir := filepath.Dir(dbPath)
		if err = os.MkdirAll(dbDir, 0755); err != nil {
			return
		}

		// 打开数据库连接
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
		if err != nil {
			return
		}

		// 设置连接池参数
		db.SetMaxOpenConns(1) // SQLite 建议单连接
		db.SetMaxIdleConns(1)

		// 测试连接
		if err = db.Ping(); err != nil {
			return
		}

		// 初始化表结构
		err = createTables()
	})
	return err
}

// GetDB 获取数据库连接
func GetDB() *sql.DB {
	return db
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// createTables 创建所有必要的表
func createTables() error {
	schemas := []string{
		// Agents 表
		`CREATE TABLE IF NOT EXISTS agents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			icon TEXT NOT NULL,
			system_prompt TEXT NOT NULL,
			default_task TEXT,
			model TEXT NOT NULL,
			enable_file_read BOOLEAN NOT NULL DEFAULT 1,
			enable_file_write BOOLEAN NOT NULL DEFAULT 1,
			enable_network BOOLEAN NOT NULL DEFAULT 1,
			hooks TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Agent Runs 表
		`CREATE TABLE IF NOT EXISTS agent_runs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			agent_id INTEGER NOT NULL,
			agent_name TEXT NOT NULL,
			agent_icon TEXT NOT NULL,
			task TEXT NOT NULL,
			model TEXT NOT NULL,
			project_path TEXT NOT NULL,
			session_id TEXT NOT NULL,
			status TEXT NOT NULL,
			pid INTEGER,
			process_started_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
		)`,

		// App Settings 表
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Slash Commands 表
		`CREATE TABLE IF NOT EXISTS slash_commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL,
			command TEXT NOT NULL,
			icon TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Usage Stats 表
		`CREATE TABLE IF NOT EXISTS usage_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			project_path TEXT NOT NULL,
			model TEXT NOT NULL,
			input_tokens INTEGER NOT NULL DEFAULT 0,
			output_tokens INTEGER NOT NULL DEFAULT 0,
			total_tokens INTEGER NOT NULL DEFAULT 0,
			cost_usd REAL NOT NULL DEFAULT 0.0,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			duration_ms INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Checkpoints 表
		`CREATE TABLE IF NOT EXISTS checkpoints (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			project_path TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			message_id TEXT,
			file_states TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 创建索引
		`CREATE INDEX IF NOT EXISTS idx_agent_runs_agent_id ON agent_runs(agent_id)`,
		`CREATE INDEX IF NOT EXISTS idx_agent_runs_session_id ON agent_runs(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_agent_runs_status ON agent_runs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_usage_stats_session_id ON usage_stats(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_usage_stats_start_time ON usage_stats(start_time)`,
		`CREATE INDEX IF NOT EXISTS idx_checkpoints_session_id ON checkpoints(session_id)`,
	}

	mu.Lock()
	defer mu.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, schema := range schemas {
		if _, err := tx.Exec(schema); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExecuteInTransaction 在事务中执行函数
func ExecuteInTransaction(fn func(*sql.Tx) error) error {
	mu.Lock()
	defer mu.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
