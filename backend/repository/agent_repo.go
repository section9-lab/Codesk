package repository

import (
	"Codesk/backend/model"
	"database/sql"
	"fmt"
	"time"
)

// AgentRepository Agent 数据访问层
type AgentRepository struct{}

// NewAgentRepository 创建 Agent 仓库实例
func NewAgentRepository() *AgentRepository {
	return &AgentRepository{}
}

// Create 创建新 Agent
func (r *AgentRepository) Create(agent *model.Agent) error {
	query := `
		INSERT INTO agents (name, icon, system_prompt, default_task, model, 
			enable_file_read, enable_file_write, enable_network, hooks)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := GetDB().Exec(query,
		agent.Name,
		agent.Icon,
		agent.SystemPrompt,
		agent.DefaultTask,
		agent.Model,
		agent.EnableFileRead,
		agent.EnableFileWrite,
		agent.EnableNetwork,
		agent.Hooks,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	*agent.ID = id
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	
	return nil
}

// GetByID 根据 ID 获取 Agent
func (r *AgentRepository) GetByID(id int64) (*model.Agent, error) {
	query := `
		SELECT id, name, icon, system_prompt, default_task, model,
			enable_file_read, enable_file_write, enable_network, hooks,
			created_at, updated_at
		FROM agents WHERE id = ?
	`
	
	agent := &model.Agent{}
	var createdAt, updatedAt string
	
	err := GetDB().QueryRow(query, id).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Icon,
		&agent.SystemPrompt,
		&agent.DefaultTask,
		&agent.Model,
		&agent.EnableFileRead,
		&agent.EnableFileWrite,
		&agent.EnableNetwork,
		&agent.Hooks,
		&createdAt,
		&updatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	
	agent.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	agent.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	
	return agent, nil
}

// List 获取所有 Agents
func (r *AgentRepository) List() ([]*model.Agent, error) {
	query := `
		SELECT id, name, icon, system_prompt, default_task, model,
			enable_file_read, enable_file_write, enable_network, hooks,
			created_at, updated_at
		FROM agents ORDER BY created_at DESC
	`
	
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}
	defer rows.Close()
	
	var agents []*model.Agent
	for rows.Next() {
		agent := &model.Agent{}
		var createdAt, updatedAt string
		
		err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Icon,
			&agent.SystemPrompt,
			&agent.DefaultTask,
			&agent.Model,
			&agent.EnableFileRead,
			&agent.EnableFileWrite,
			&agent.EnableNetwork,
			&agent.Hooks,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		agent.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		agent.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		
		agents = append(agents, agent)
	}
	
	return agents, nil
}

// Update 更新 Agent
func (r *AgentRepository) Update(agent *model.Agent) error {
	query := `
		UPDATE agents SET
			name = ?, icon = ?, system_prompt = ?, default_task = ?,
			model = ?, enable_file_read = ?, enable_file_write = ?,
			enable_network = ?, hooks = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	
	_, err := GetDB().Exec(query,
		agent.Name,
		agent.Icon,
		agent.SystemPrompt,
		agent.DefaultTask,
		agent.Model,
		agent.EnableFileRead,
		agent.EnableFileWrite,
		agent.EnableNetwork,
		agent.Hooks,
		agent.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}
	
	agent.UpdatedAt = time.Now()
	return nil
}

// Delete 删除 Agent
func (r *AgentRepository) Delete(id int64) error {
	query := `DELETE FROM agents WHERE id = ?`
	
	_, err := GetDB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}
	
	return nil
}

// CreateRun 创建 Agent 运行记录
func (r *AgentRepository) CreateRun(run *model.AgentRun) error {
	query := `
		INSERT INTO agent_runs (agent_id, agent_name, agent_icon, task, model,
			project_path, session_id, status, pid, process_started_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := GetDB().Exec(query,
		run.AgentID,
		run.AgentName,
		run.AgentIcon,
		run.Task,
		run.Model,
		run.ProjectPath,
		run.SessionID,
		run.Status,
		run.PID,
		run.ProcessStartedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create agent run: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	*run.ID = id
	run.CreatedAt = time.Now()
	
	return nil
}

// GetRunByID 根据 ID 获取运行记录
func (r *AgentRepository) GetRunByID(id int64) (*model.AgentRun, error) {
	query := `
		SELECT id, agent_id, agent_name, agent_icon, task, model,
			project_path, session_id, status, pid, process_started_at,
			created_at, completed_at
		FROM agent_runs WHERE id = ?
	`
	
	run := &model.AgentRun{}
	var createdAt string
	var completedAt sql.NullString
	var processStartedAt sql.NullString
	
	err := GetDB().QueryRow(query, id).Scan(
		&run.ID,
		&run.AgentID,
		&run.AgentName,
		&run.AgentIcon,
		&run.Task,
		&run.Model,
		&run.ProjectPath,
		&run.SessionID,
		&run.Status,
		&run.PID,
		&processStartedAt,
		&createdAt,
		&completedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent run not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get agent run: %w", err)
	}
	
	run.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if completedAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
		run.CompletedAt = &t
	}
	if processStartedAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", processStartedAt.String)
		run.ProcessStartedAt = &t
	}
	
	return run, nil
}

// ListRuns 获取 Agent 的所有运行记录
func (r *AgentRepository) ListRuns(agentID int64) ([]*model.AgentRun, error) {
	query := `
		SELECT id, agent_id, agent_name, agent_icon, task, model,
			project_path, session_id, status, pid, process_started_at,
			created_at, completed_at
		FROM agent_runs WHERE agent_id = ? ORDER BY created_at DESC
	`
	
	rows, err := GetDB().Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list agent runs: %w", err)
	}
	defer rows.Close()
	
	var runs []*model.AgentRun
	for rows.Next() {
		run := &model.AgentRun{}
		var createdAt string
		var completedAt sql.NullString
		var processStartedAt sql.NullString
		
		err := rows.Scan(
			&run.ID,
			&run.AgentID,
			&run.AgentName,
			&run.AgentIcon,
			&run.Task,
			&run.Model,
			&run.ProjectPath,
			&run.SessionID,
			&run.Status,
			&run.PID,
			&processStartedAt,
			&createdAt,
			&completedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		run.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if completedAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
			run.CompletedAt = &t
		}
		if processStartedAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", processStartedAt.String)
			run.ProcessStartedAt = &t
		}
		
		runs = append(runs, run)
	}
	
	return runs, nil
}

// UpdateRunStatus 更新运行状态
func (r *AgentRepository) UpdateRunStatus(id int64, status string, completedAt *time.Time) error {
	query := `UPDATE agent_runs SET status = ?, completed_at = ? WHERE id = ?`
	
	_, err := GetDB().Exec(query, status, completedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update run status: %w", err)
	}
	
	return nil
}
