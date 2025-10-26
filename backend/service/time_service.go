package service

import (
	"time"
)

// TimeService 处理时间相关的业务逻辑
type TimeService struct {
}

// NewTimeService 创建新的时间服务实例
func NewTimeService() *TimeService {
	return &TimeService{}
}

// GetCurrentTime 获取当前时间
func (s *TimeService) GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
