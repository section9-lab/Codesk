package proxy

import (
	"Codesk/backend/model"
	"Codesk/backend/repository"
	// "encoding/json"
	"fmt"
	"os"
)

// ProxyService 代理设置服务
type ProxyService struct {
	repo *repository.StorageRepository
}

// NewProxyService 创建代理服务实例
func NewProxyService() *ProxyService {
	return &ProxyService{
		repo: repository.NewStorageRepository(),
	}
}

// GetProxySettings 获取代理设置
func (s *ProxyService) GetProxySettings() (*model.ProxySettings, error) {
	settings := &model.ProxySettings{}

	// 从数据库读取设置
	if enabled, err := s.repo.GetSetting("proxy_enabled"); err == nil && enabled == "true" {
		settings.Enabled = true
	}

	if httpProxy, err := s.repo.GetSetting("proxy_http"); err == nil && httpProxy != "" {
		settings.HTTPProxy = &httpProxy
	}

	if httpsProxy, err := s.repo.GetSetting("proxy_https"); err == nil && httpsProxy != "" {
		settings.HTTPSProxy = &httpsProxy
	}

	if noProxy, err := s.repo.GetSetting("proxy_no"); err == nil && noProxy != "" {
		settings.NoProxy = &noProxy
	}

	if allProxy, err := s.repo.GetSetting("proxy_all"); err == nil && allProxy != "" {
		settings.AllProxy = &allProxy
	}

	return settings, nil
}

// SaveProxySettings 保存代理设置
func (s *ProxyService) SaveProxySettings(settings *model.ProxySettings) error {
	// 保存到数据库
	if err := s.repo.SetSetting("proxy_enabled", fmt.Sprintf("%t", settings.Enabled)); err != nil {
		return err
	}

	if settings.HTTPProxy != nil {
		if err := s.repo.SetSetting("proxy_http", *settings.HTTPProxy); err != nil {
			return err
		}
	}

	if settings.HTTPSProxy != nil {
		if err := s.repo.SetSetting("proxy_https", *settings.HTTPSProxy); err != nil {
			return err
		}
	}

	if settings.NoProxy != nil {
		if err := s.repo.SetSetting("proxy_no", *settings.NoProxy); err != nil {
			return err
		}
	}

	if settings.AllProxy != nil {
		if err := s.repo.SetSetting("proxy_all", *settings.AllProxy); err != nil {
			return err
		}
	}

	// 应用代理设置到环境变量
	s.ApplyProxySettings(settings)

	return nil
}

// ApplyProxySettings 应用代理设置到环境变量
func (s *ProxyService) ApplyProxySettings(settings *model.ProxySettings) {
	if !settings.Enabled {
		// 清除代理环境变量
		os.Unsetenv("HTTP_PROXY")
		os.Unsetenv("HTTPS_PROXY")
		os.Unsetenv("NO_PROXY")
		os.Unsetenv("ALL_PROXY")
		os.Unsetenv("http_proxy")
		os.Unsetenv("https_proxy")
		os.Unsetenv("no_proxy")
		os.Unsetenv("all_proxy")
		return
	}

	// 设置代理环境变量
	if settings.HTTPProxy != nil && *settings.HTTPProxy != "" {
		os.Setenv("HTTP_PROXY", *settings.HTTPProxy)
		os.Setenv("http_proxy", *settings.HTTPProxy)
	}

	if settings.HTTPSProxy != nil && *settings.HTTPSProxy != "" {
		os.Setenv("HTTPS_PROXY", *settings.HTTPSProxy)
		os.Setenv("https_proxy", *settings.HTTPSProxy)
	}

	if settings.NoProxy != nil && *settings.NoProxy != "" {
		os.Setenv("NO_PROXY", *settings.NoProxy)
		os.Setenv("no_proxy", *settings.NoProxy)
	}

	if settings.AllProxy != nil && *settings.AllProxy != "" {
		os.Setenv("ALL_PROXY", *settings.AllProxy)
		os.Setenv("all_proxy", *settings.AllProxy)
	}
}
