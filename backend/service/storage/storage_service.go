package storage

import (
	"Codesk/backend/model"
	"Codesk/backend/repository"
	"fmt"
)

// StorageService 存储管理服务
type StorageService struct {
	repo *repository.StorageRepository
}

// NewStorageService 创建存储服务实例
func NewStorageService() *StorageService {
	return &StorageService{
		repo: repository.NewStorageRepository(),
	}
}

// ListTables 列出所有表
func (s *StorageService) ListTables() ([]model.TableInfo, error) {
	return s.repo.ListTables()
}

// ReadTable 读取表数据
func (s *StorageService) ReadTable(tableName string, limit, offset int) ([]model.TableRow, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ReadTable(tableName, limit, offset)
}

// UpdateRow 更新表行（简化实现）
func (s *StorageService) UpdateRow(tableName string, id int64, data model.TableRow) error {
	// 构建 UPDATE 语句
	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	var args []interface{}
	first := true

	for key, value := range data {
		if key == "id" {
			continue
		}
		if !first {
			query += ", "
		}
		query += fmt.Sprintf("%s = ?", key)
		args = append(args, value)
		first = false
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.repo.ExecuteSQL(query, args...)
	return err
}

// DeleteRow 删除表行
func (s *StorageService) DeleteRow(tableName string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := s.repo.ExecuteSQL(query, id)
	return err
}

// InsertRow 插入表行
func (s *StorageService) InsertRow(tableName string, data model.TableRow) error {
	// 构建 INSERT 语句
	query := fmt.Sprintf("INSERT INTO %s (", tableName)
	values := "VALUES ("
	var args []interface{}
	first := true

	for key, value := range data {
		if key == "id" {
			continue
		}
		if !first {
			query += ", "
			values += ", "
		}
		query += key
		values += "?"
		args = append(args, value)
		first = false
	}

	query += ") " + values + ")"

	_, err := s.repo.ExecuteSQL(query, args...)
	return err
}

// ExecuteSQL 执行 SQL 语句
func (s *StorageService) ExecuteSQL(query string, args ...interface{}) (*model.SQLResult, error) {
	return s.repo.ExecuteSQL(query, args...)
}

// ResetDatabase 重置数据库
func (s *StorageService) ResetDatabase() error {
	return s.repo.ResetDatabase()
}

// GetSetting 获取设置
func (s *StorageService) GetSetting(key string) (string, error) {
	return s.repo.GetSetting(key)
}

// SetSetting 设置值
func (s *StorageService) SetSetting(key, value string) error {
	return s.repo.SetSetting(key, value)
}
