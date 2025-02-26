package types

import (
	"fmt"
	"strings"
)

const (
	CompletionToolTypeFunction CompletionToolType = "function"
)

type CompletionToolType string

type CompletionRequest struct {
	Model        string
	Tools        []CompletionTool
	Messages     []CompletionMessage
	MaxToken     int
	Temperature  *float32
	JSONResponse bool
	Cache        *bool
}

type CompletionTool struct {
	Type     CompletionToolType           `json:"type"`
	Function CompletionFunctionDefinition `json:"function,omitempty"`
}

type CompletionFunctionDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Domain      string      `json:"domain,omitempty"`
	Parameters  *JSONSchema `json:"parameters"`
}

// Chat message role defined by the OpenAI API.
const (
	CompletionMessageRoleTypeUser      = CompletionMessageRoleType("user")
	CompletionMessageRoleTypeSystem    = CompletionMessageRoleType("system")
	CompletionMessageRoleTypeAssistant = CompletionMessageRoleType("assistant")
	CompletionMessageRoleTypeTool      = CompletionMessageRoleType("tool")
)

type CompletionMessageRoleType string

type CompletionMessage struct {
	Role     CompletionMessageRoleType `json:"role,omitempty"`
	Content  []ContentPart             `json:"content,omitempty" column:"name=Message,jsonpath=.spec.content"`
	ToolCall *CompletionToolCall       `json:"toolCall,omitempty"`
}

func (in CompletionMessage) IsToolCall() bool {
	for _, content := range in.Content {
		if content.ToolCall != nil {
			return true
		}
	}
	return false
}

func Text(text string) []ContentPart {
	return []ContentPart{
		{
			Text: text,
		},
	}
}

func (in CompletionMessage) String() string {
	buf := strings.Builder{}
	for i, content := range in.Content {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(content.Text)
		if content.ToolCall != nil {
			buf.WriteString(fmt.Sprintf("tool call %s -> %s", content.ToolCall.Function.Name, content.ToolCall.Function.Arguments))
		}
	}
	return buf.String()
}

type ContentPart struct {
	Text     string              `json:"text,omitempty"`
	ToolCall *CompletionToolCall `json:"toolCall,omitempty"`
}

type CompletionToolCall struct {
	Index    *int                   `json:"index,omitempty"`
	ID       string                 `json:"id,omitempty"`
	Type     CompletionToolType     `json:"type,omitempty"`
	Function CompletionFunctionCall `json:"function,omitempty"`
}

type CompletionFunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}
