package service

import "github.com/inhere/kite-go/internal/service/aiservice"

// AI 获取 AI 服务
func AI() *aiservice.AIService {
	return aiservice.New()
}
