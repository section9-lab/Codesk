package service

import (
	"fmt"
)

// GreetService 处理问候相关的业务逻辑
type GreetService struct {
}

// NewGreetService 创建新的问候服务实例
func NewGreetService() *GreetService {
	return &GreetService{}
}

// Greet 生成问候消息
func (s *GreetService) Greet(name string) string {
	// 这里可以扩展复杂的业务逻辑
	// 比如：用户验证、个性化问候、多语言支持等
	if name == "" {
		return "Hello there, It's show time!"
	}
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 未来可以添加更多方法：
// func (s *GreetService) GetPersonalizedGreeting(userID string) string { ... }
// func (s *GreetService) GetGreetingHistory(userID string) []string { ... }