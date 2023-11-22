package interfaces

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type IQueryMatch interface {
	GetID() string
	SetID(id string) error
	GetScore() float32
	SetScore(float32) error
	GetMetaData() map[string]string
	SetMetaData(map[string]string) error
	Decode([]byte) error
}

type UniversalRole string

type ChatHistory struct {
	Role    UniversalRole
	Content string
}

type Prompt struct {
	Body            string                 // The actual prompt
	RawBody         string                 // The actual prompt
	History         []*ChatHistory         // The history of the conversation
	HistoryLen      int                    // The length of the history
	Model           string                 // The model being used
	ContextToRender []string               // Rendered context (within token limit)
	ContextTitles   []string               // titles and references to the content
	MetaData        map[string]interface{} // additional template data
	Template        string
	RenderedPrompt  string
	Instructions    string // You are an AI assistant that is happy, helpful and tries to offer insightful answers
	Snippets        map[string]string
	Stop            []string // Human: AI:
	PromptContext   []IQueryMatch
	//Functions         []openai.FunctionDefinition
	IsFuncResponse    bool
	ModelRoles        *ModelRole
	ExpandedResources []ExpandedResource
	tpl               *template.Template
	WithUserID        string
	WithBotID         string
}

type ExpandedResource struct {
	URI    string
	Raw    []byte
	Clean  string
	Status int64
}

// swagger:model ModelRole
type ModelRole struct {
	// The system role start tag (e.g. <system>), can be empty
	// example: "<system>"
	System string

	// The system role end tag (e.g. </system>), can be empty
	// example: "</system>"
	SystemEnd string

	// The instruction role start tag (e.g. 'instructions')
	// example: "<instructions>"
	Instruction string

	// The instruction role end tag (e.g. '</instructions>')
	// example: "</instructions>"
	InstructionEnd string

	// The user role start tag (e.g. '<user>')
	// example: "<user>"
	User string

	// The user role end tag (e.g. '</user>')
	// example: "</user>"
	UserEnd string

	// The AI role start tag (e.g. '<ai>')
	// example: "<ai>"
	AI string

	// The AI role end tag (e.g. '</ai>')
	// example: "</ai>"
	AIEnd string
}

var DefaultModelRoles = &ModelRole{
	System:      "",
	Instruction: "",
	User:        "### Instruction:",
	AI:          "### Assistant:",
}

var DEFAULT_TEMPLATE = `
{{ if .ContextToRender }}Use the following context to help with your response:
{{ range $ctx := .ContextToRender }}
{{$ctx}}
{{ end }}
===={{ end }}

Human: {{.Body}}
AI:
`

func NewPrompt(instructions, promptTemplate string, withContext []IQueryMatch, PromptFormat *ModelRole, withSnippets map[string]string, withUser string, withBot string) *Prompt {
	b := &Prompt{
		ContextToRender: make([]string, 0),
		ContextTitles:   make([]string, 0),
		Snippets:        withSnippets,
		WithUserID:      withUser,
		WithBotID:       withBot,
	}

	b.Template = promptTemplate
	b.Instructions = instructions

	b.ModelRoles = DefaultModelRoles
	if PromptFormat != nil {
		b.ModelRoles = PromptFormat
	}

	if b.Template == "" {
		b.Template = DEFAULT_TEMPLATE
	}

	// load from a file if it's a file
	pf := strings.HasPrefix(promptTemplate, "file://")
	if pf {
		promptTemplate = strings.ReplaceAll(promptTemplate, "file://", "")
		if !filepath.IsAbs(promptTemplate) {
			ex, err := os.Executable()
			if err != nil {
				log.Println(err)
				return nil
			}
			exPath := filepath.Dir(ex)
			promptTemplate = filepath.Join(exPath, promptTemplate)
		}

		c, err := os.ReadFile(promptTemplate)
		if err != nil {
			log.Println(err)
			return nil
		}

		b.Template = string(c)
	}

	if withContext != nil {
		b.PromptContext = withContext
		b.ContextTitles = make([]string, len(withContext))
		b.ContextToRender = make([]string, len(withContext))

		for i, _ := range withContext {
			title := fmt.Sprintf("%s", withContext[i].GetMetaData()["title"])
			b.ContextTitles[i] = title
			b.ContextToRender[i] = withContext[i].GetMetaData()["text"]
		}

	}

	return b
}
