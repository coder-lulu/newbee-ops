package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// TemplateEngine 模板引擎接口
type TemplateEngine interface {
	// Render 渲染模板
	Render(templateContent string, parameters map[string]interface{}) (string, error)
	// Validate 验证模板语法
	Validate(templateContent string) error
	// RenderWithCache 使用缓存渲染模板
	RenderWithCache(templateName string, templateContent string, parameters map[string]interface{}) (string, error)
}

// templateEngineImpl 模板引擎实现
type templateEngineImpl struct {
	// 模板缓存
	cache     map[string]*template.Template
	cacheLock sync.RWMutex
	// 自定义函数映射
	funcMap template.FuncMap
	// 日志
	logger logx.Logger
}

// NewTemplateEngine 创建模板引擎实例
func NewTemplateEngine() TemplateEngine {
	engine := &templateEngineImpl{
		cache:  make(map[string]*template.Template),
		logger: logx.WithContext(nil),
	}

	// 初始化自定义函数
	engine.funcMap = template.FuncMap{
		// JSON序列化
		"json": engine.toJSON,
		// 默认值
		"default": engine.defaultValue,
		// 字符串包含
		"contains": strings.Contains,
		// 字符串分割
		"split": strings.Split,
		// 字符串连接
		"join": strings.Join,
		// 字符串替换
		"replace": strings.ReplaceAll,
		// 转大写
		"upper": strings.ToUpper,
		// 转小写
		"lower": strings.ToLower,
		// 去空格
		"trim": strings.TrimSpace,
		// 格式化时间
		"formatTime": engine.formatTime,
		// 时间戳
		"now": time.Now,
		// 字符串前缀
		"hasPrefix": strings.HasPrefix,
		// 字符串后缀
		"hasSuffix": strings.HasSuffix,
	}

	return engine
}

// Render 渲染模板
func (e *templateEngineImpl) Render(templateContent string, parameters map[string]interface{}) (string, error) {
	// 创建临时模板
	tmpl, err := template.New("script").Funcs(e.funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, parameters); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buf.String(), nil
}

// Validate 验证模板语法
func (e *templateEngineImpl) Validate(templateContent string) error {
	// 尝试编译模板
	_, err := template.New("validate").Funcs(e.funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}
	return nil
}

// RenderWithCache 使用缓存渲染模板
func (e *templateEngineImpl) RenderWithCache(templateName string, templateContent string, parameters map[string]interface{}) (string, error) {
	// 先尝试从缓存获取
	e.cacheLock.RLock()
	tmpl, exists := e.cache[templateName]
	e.cacheLock.RUnlock()

	if !exists {
		// 不存在则编译并缓存
		var err error
		tmpl, err = template.New(templateName).Funcs(e.funcMap).Parse(templateContent)
		if err != nil {
			return "", fmt.Errorf("template parse error: %w", err)
		}

		// 缓存模板
		e.cacheLock.Lock()
		e.cache[templateName] = tmpl
		e.cacheLock.Unlock()

		e.logger.Infof("Template cached: %s", templateName)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, parameters); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buf.String(), nil
}

// ClearCache 清除缓存
func (e *templateEngineImpl) ClearCache() {
	e.cacheLock.Lock()
	defer e.cacheLock.Unlock()
	e.cache = make(map[string]*template.Template)
	e.logger.Info("Template cache cleared")
}

// ClearCacheByName 清除指定模板缓存
func (e *templateEngineImpl) ClearCacheByName(templateName string) {
	e.cacheLock.Lock()
	defer e.cacheLock.Unlock()
	delete(e.cache, templateName)
	e.logger.Infof("Template cache cleared: %s", templateName)
}

// 自定义函数实现

// toJSON 将对象转换为JSON字符串
func (e *templateEngineImpl) toJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// defaultValue 提供默认值
func (e *templateEngineImpl) defaultValue(defaultVal interface{}, val interface{}) interface{} {
	if val == nil || val == "" {
		return defaultVal
	}
	return val
}

// formatTime 格式化时间
func (e *templateEngineImpl) formatTime(format string, t time.Time) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return t.Format(format)
}

// GetCacheSize 获取缓存大小
func (e *templateEngineImpl) GetCacheSize() int {
	e.cacheLock.RLock()
	defer e.cacheLock.RUnlock()
	return len(e.cache)
}

// TemplateRenderResult 模板渲染结果
type TemplateRenderResult struct {
	Content   string        // 渲染后内容
	Duration  time.Duration // 渲染耗时
	FromCache bool          // 是否来自缓存
}

// RenderWithMetrics 带性能指标的渲染
func (e *templateEngineImpl) RenderWithMetrics(templateName string, templateContent string, parameters map[string]interface{}) (*TemplateRenderResult, error) {
	startTime := time.Now()

	// 检查缓存
	e.cacheLock.RLock()
	_, fromCache := e.cache[templateName]
	e.cacheLock.RUnlock()

	// 渲染
	content, err := e.RenderWithCache(templateName, templateContent, parameters)
	if err != nil {
		return nil, err
	}

	return &TemplateRenderResult{
		Content:   content,
		Duration:  time.Since(startTime),
		FromCache: fromCache,
	}, nil
}
