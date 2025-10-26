package model

// TableInfo 数据库表信息
type TableInfo struct {
	Name       string   `json:"name"`
	RowCount   int      `json:"row_count"`
	Columns    []string `json:"columns"`
}

// TableRow 表行数据
type TableRow map[string]interface{}

// SQLResult SQL 执行结果
type SQLResult struct {
	RowsAffected int64       `json:"rows_affected"`
	LastInsertID int64       `json:"last_insert_id"`
	Rows         []TableRow  `json:"rows,omitempty"`
}

// ProxySettings 代理设置
type ProxySettings struct {
	Enabled    bool    `json:"enabled"`
	HTTPProxy  *string `json:"http_proxy"`
	HTTPSProxy *string `json:"https_proxy"`
	NoProxy    *string `json:"no_proxy"`
	AllProxy   *string `json:"all_proxy"`
}

// SlashCommand 斜杠命令
type SlashCommand struct {
	ID          *int64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Icon        string `json:"icon"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// HooksConfig Hooks 配置
type HooksConfig struct {
	PreExecution  []Hook `json:"pre_execution"`
	PostExecution []Hook `json:"post_execution"`
	OnError       []Hook `json:"on_error"`
}

// Hook 单个 Hook 配置
type Hook struct {
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Enabled     bool     `json:"enabled"`
	Description string   `json:"description"`
}
