package simpleai

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gookit/goutil/maputil"
	"github.com/sashabaranov/go-openai"
)

// FunctionCall 函数调用请求
type FunctionCall = openai.FunctionCall

// FunctionSpec 函数调用元数据
type FunctionSpec = openai.FunctionDefinition

type FunctionHandler func(ctx context.Context, params maputil.Map) (any, error)

type ToolCallValidator interface {
	Validate(call FunctionCall) error
}

type DefaultValidator struct {
	allowedFunctions map[string]bool
}

func (v *DefaultValidator) Validate(call FunctionCall) error {
	if !v.allowedFunctions[call.Name] {
		return fmt.Errorf("function %s is not allowed", call.Name)
	}
	// 参数类型校验...
	return nil
}

// ToolRegistry 工具函数注册表
type ToolRegistry struct {
	tools map[string]FunctionHandler
	specs []FunctionSpec
}

// NewToolRegistry 初始化函数注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]FunctionHandler),
		specs: []FunctionSpec{},
	}
}

// Register 注册函数
func (tr *ToolRegistry) Register(spec FunctionSpec, handler FunctionHandler) error {
	if _, exists := tr.tools[spec.Name]; exists {
		return fmt.Errorf("function %s already registered", spec.Name)
	}

	tr.tools[spec.Name] = handler
	tr.specs = append(tr.specs, spec)
	return nil
}

func (tr *ToolRegistry) Specs() []FunctionSpec {
	return tr.specs
}

// Handler get by name
func (tr *ToolRegistry) Handler(name string) (FunctionHandler, bool) {
	handler, ok := tr.tools[name]
	return handler, ok
}

// CallHandler 处理函数调用
func (tr *ToolRegistry) CallHandler(ctx context.Context, call FunctionCall) (any, error) {
	handler, ok := tr.tools[call.Name]
	if !ok {
		return nil, fmt.Errorf("tool function %q not found", call.Name)
	}

	params := make(maputil.Map)
	err := json.Unmarshal([]byte(call.Arguments), &params)
	if err != nil {
		return nil, fmt.Errorf("invalid function arguments: %w", err)
	}
	return handler(ctx, params)
}

var toolRegistry = NewToolRegistry()
var openaiClient = openai.NewClient("")

func callAIWithFunctions(query string, tools []FunctionSpec) (string, []FunctionCall) {
	oaiTools := make([]openai.Tool, len(tools))
	for i, tool := range tools {
		oaiTools[i] = openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &tool,
		}
	}

	// 构造带函数定义的API请求
	req := openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: query},
		},
		Tools: oaiTools,
	}

	// 发送API请求
	resp, err := openaiClient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", nil
	}

	// 解析响应
	var functionCalls []FunctionCall
	finalResponse := ""

	for _, choice := range resp.Choices {
		if choice.Message.FunctionCall != nil {
			functionCalls = append(functionCalls, *choice.Message.FunctionCall)
		} else {
			finalResponse += choice.Message.Content
		}
	}
	return finalResponse, functionCalls
}

// 处理函数调用循环
func processFunctionCalls(input string) string {
	// 匹配工具
	response, calls := callAIWithFunctions(input, toolRegistry.Specs())

	for _, call := range calls {
		result, err := toolRegistry.CallHandler(context.Background(), call)
		if err != nil {
			continue
		}

		jsonResult, _ := json.Marshal(result)

		// 将tool执行结果反馈给AI
		followupResp, _ := callAIWithFunctions(
			fmt.Sprintf("Function result: %s", jsonResult),
			nil,
		)
		response += "\n" + followupResp
	}

	return response
}
