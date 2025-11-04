package aicmd

import (
	"context"
	"fmt"

	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/strutil"
	"github.com/inhere/kite-go/internal/apputil"
	"github.com/sashabaranov/go-openai"
)

var langMap = map[string]string{
	"en": "English",
	"zh": "Chinese",
	"ja": "Japanese",
	"ko": "Korean",
	"fr": "French",
	"de": "German",
	"es": "Spanish",
	"it": "Italian",
	"pt": "Portuguese",
	"ru": "Russian",
	"ar": "Arabic",
	"hi": "Hindi",
	"bn": "Bengali",
}

// NewTranslateCmd instance
func NewTranslateCmd() *gcli.Command {
	var transOpts = struct {
		source string
		target string
		model  string
	}{}

	return &gcli.Command{
		Name:    "translate",
		Aliases: []string{"tr", "trans"},
		Desc:    "translate input text by AI model",
		Config: func(c *gcli.Command) {
			c.StrOpt2(&transOpts.source, "source,s", "source language, default is auto check. eg: en, zh")
			c.StrOpt2(&transOpts.target, "target,t", "target language, default is Chinese. eg: en, zh")
			c.StrOpt2(&transOpts.model, "model", "the model name, default is gpt-3.5-turbo\n allow set by TRANSLATE_MODEL_NAME")

			c.AddArg("text", "input text to translate. use '@c' refer clipboard", true)
		},
		Func: func(c *gcli.Command, args []string) error {
			text, err := apputil.ReadSource(c.Arg("text").String())
			if err != nil {
				return err
			}

			srcLang := strutil.OrElse(transOpts.source, "input")
			tgtLang := strutil.OrElse(transOpts.target, "Chinese")
			if realLang := langMap[srcLang]; realLang != "" {
				srcLang = realLang
			}
			if realLang := langMap[tgtLang]; realLang != "" {
				tgtLang = realLang
			}
			colorp.Infoln("- Translate from:", strutil.OrCond(srcLang == "input", "auto", srcLang), "to:", tgtLang)

			modelName := envutil.GetOne([]string{"TRANSLATE_MODEL_NAME", "DEFAULT_LLM_MODEL"}, "gpt-3.5-turbo")
			colorp.Infoln("- Translate by model:", modelName)

			fmt.Println("INPUT:", text)
			fmt.Println("---------------------------------------------------------------------")

			client := apputil.NewOpenaiClient()
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					// Model: openai.GPT3Dot5Turbo,
					Model: modelName,
					Messages: []openai.ChatCompletionMessage{
						// system message
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: fmt.Sprintf("You are a helpful assistant that translates %s to %s.", srcLang, tgtLang),
						},
						// user message
						{
							Role:    openai.ChatMessageRoleUser,
							Content: text,
						},
					},
				},
			)

			if err == nil {
				fmt.Println("RESULT:", resp.Choices[0].Message.Content)
			}
			return err
		},
	}
}
