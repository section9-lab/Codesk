package repository

import (
	"Codesk/backend/model"
	"database/sql"
	"fmt"
)

// StorageRepository 通用存储数据访问层
type StorageRepository struct{}

// NewStorageRepository 创建存储仓库实例
func NewStorageRepository() *StorageRepository {
	return &StorageRepository{}
}

// ListTables 列出所有表
func (r *StorageRepository) ListTables() ([]model.TableInfo, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`
	
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()
	
	var tables []model.TableInfo
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		
		// 获取行数
		var rowCount int
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
		GetDB().QueryRow(countQuery).Scan(&rowCount)
		
		// 获取列信息
		columns, _ := r.GetTableColumns(tableName)
		
		tables = append(tables, model.TableInfo{
			Name:     tableName,
			RowCount: rowCount,
			Columns:  columns,
		})
	}
	
	return tables, nil
}

// GetTableColumns 获取表的列名
func (r *StorageRepository) GetTableColumns(tableName string) ([]string, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}
	defer rows.Close()
	
	var columns []string
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue sql.NullString
		
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		
		columns = append(columns, name)
	}
	
	return columns, nil
}

// ReadTable 读取表数据
func (r *StorageRepository) ReadTable(tableName string, limit, offset int) ([]model.TableRow, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", tableName)
	
	rows, err := GetDB().Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}
	defer rows.Close()
	
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	
	var results []model.TableRow
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		
		row := make(model.TableRow)
		for i, col := range columns {
			row[col] = values[i]
		}
		
		results = append(results, row)
	}
	
	return results, nil
}

// ExecuteSQL 执行 SQL 语句
func (r *StorageRepository) ExecuteSQL(query string, args ...interface{}) (*model.SQLResult, error) {
	result, err := GetDB().Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	lastInsertID, _ := result.LastInsertId()
	
	return &model.SQLResult{
		RowsAffected: rowsAffected,
		LastInsertID: lastInsertID,
	}, nil
}

// GetSetting 获取应用设置
func (r *StorageRepository) GetSetting(key string) (string, error) {
	query := `SELECT value FROM app_settings WHERE key = ?`
	
	var value string
	err := GetDB().QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get setting: %w", err)
	}
	
	return value, nil
}

// SetSetting 设置应用设置
func (r *StorageRepository) SetSetting(key, value string) error {
	query := `
		INSERT INTO app_settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := GetDB().Exec(query, key, value, value)
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}
	
	return nil
}

// ResetDatabase 重置数据库（清空所有表）
func (r *StorageRepository) ResetDatabase() error {
	tables := []string{"agents", "agent_runs", "usage_stats", "checkpoints", "slash_commands"}
	
	return ExecuteInTransaction(func(tx *sql.Tx) error {
		for _, table := range tables {
			query := fmt.Sprintf("DELETE FROM %s", table)
			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("failed to clear table %s: %w", table, err)
			}
		}
		return nil
	})
}
