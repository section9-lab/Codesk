package slash

import (
	"Codesk/backend/model"
	"Codesk/backend/repository"
	"database/sql"
	"fmt"
	"time"
)

// SlashService 斜杠命令服务
type SlashService struct{}

// NewSlashService 创建斜杠命令服务实例
func NewSlashService() *SlashService {
	return &SlashService{}
}

// ListCommands 列出所有斜杠命令
func (s *SlashService) ListCommands() ([]*model.SlashCommand, error) {
	query := `
		SELECT id, name, description, command, icon, created_at, updated_at
		FROM slash_commands ORDER BY name
	`

	rows, err := repository.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list slash commands: %w", err)
	}
	defer rows.Close()

	var commands []*model.SlashCommand
	for rows.Next() {
		cmd := &model.SlashCommand{}
		err := rows.Scan(
			&cmd.ID,
			&cmd.Name,
			&cmd.Description,
			&cmd.Command,
			&cmd.Icon,
			&cmd.CreatedAt,
			&cmd.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// GetCommand 获取斜杠命令
func (s *SlashService) GetCommand(name string) (*model.SlashCommand, error) {
	query := `
		SELECT id, name, description, command, icon, created_at, updated_at
		FROM slash_commands WHERE name = ?
	`

	cmd := &model.SlashCommand{}
	err := repository.GetDB().QueryRow(query, name).Scan(
		&cmd.ID,
		&cmd.Name,
		&cmd.Description,
		&cmd.Command,
		&cmd.Icon,
		&cmd.CreatedAt,
		&cmd.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("slash command not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get slash command: %w", err)
	}

	return cmd, nil
}

// SaveCommand 保存斜杠命令（创建或更新）
func (s *SlashService) SaveCommand(cmd *model.SlashCommand) error {
	// 检查是否已存在
	existing, err := s.GetCommand(cmd.Name)
	if err == nil && existing != nil {
		// 更新
		return s.updateCommand(cmd)
	}

	// 创建
	return s.createCommand(cmd)
}

// createCommand 创建新命令
func (s *SlashService) createCommand(cmd *model.SlashCommand) error {
	query := `
		INSERT INTO slash_commands (name, description, command, icon)
		VALUES (?, ?, ?, ?)
	`

	result, err := repository.GetDB().Exec(query,
		cmd.Name,
		cmd.Description,
		cmd.Command,
		cmd.Icon,
	)

	if err != nil {
		return fmt.Errorf("failed to create slash command: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	*cmd.ID = id
	cmd.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	cmd.UpdatedAt = cmd.CreatedAt

	return nil
}

// updateCommand 更新命令
func (s *SlashService) updateCommand(cmd *model.SlashCommand) error {
	query := `
		UPDATE slash_commands SET
			description = ?, command = ?, icon = ?, updated_at = CURRENT_TIMESTAMP
		WHERE name = ?
	`

	_, err := repository.GetDB().Exec(query,
		cmd.Description,
		cmd.Command,
		cmd.Icon,
		cmd.Name,
	)

	if err != nil {
		return fmt.Errorf("failed to update slash command: %w", err)
	}

	cmd.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return nil
}

// DeleteCommand 删除斜杠命令
func (s *SlashService) DeleteCommand(name string) error {
	query := `DELETE FROM slash_commands WHERE name = ?`

	result, err := repository.GetDB().Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete slash command: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("slash command not found: %s", name)
	}

	return nil
}
