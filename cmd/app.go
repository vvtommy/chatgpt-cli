package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/viper"
	"github.com/vvtommy/chatgpt-cli/internal/chat"
	"github.com/vvtommy/chatgpt-cli/openai"
	"io"
	"os"
	"regexp"
	"strings"
)

func ifStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

const DefaultPrefix = ">>> "
const MultiLinePrefix = "... "

type appContext struct {
	name       string
	version    string
	apiKey     string
	organizeID string
	chat       *chat.Context

	jsonOutput bool

	prompt bool
	pipe   bool

	p *prompt.Prompt

	multiline        bool
	heredocDelimiter string
	buffer           string
}

var _shared *appContext

func newApp() *appContext {
	if _shared != nil {
		return _shared
	}
	apiKey := viper.GetString(ConfigApiKey)
	organizeID := viper.GetString(ConfigOrganizeID)

	oa := openai.NewOpenAIContext(apiKey, organizeID)
	chatOption := chat.NewChatContextOption{
		OpenAI: oa,
	}

	pipe := ifStdin()
	return &appContext{
		name:       AppName,
		version:    Version,
		apiKey:     apiKey,
		organizeID: organizeID,
		chat:       chat.NewChatContext(chatOption),
		prompt:     !pipe && argInput == "",
		pipe:       pipe,
		jsonOutput: argJSONOutput,
		buffer:     argInput,
	}
}

func (a *appContext) readStdin() string {
	var (
		err   error
		input []byte
	)

	input, err = io.ReadAll(os.Stdin)
	if err != nil {
		return ""
	}
	return string(input)
}

func (a *appContext) run() {
	if a.prompt {
		a.startPrompt()
		return
	}
	a.runOnce()
}

type Exit int

func handleExit() {
	v := recover()
	switch v := v.(type) {
	case nil:
		return
	case Exit:
		os.Exit(int(v))
	default:
		fmt.Printf("%+v", v)
	}
}

func bye() {
	fmt.Println("Bye!")
	panic(Exit(0))
}

var (
	heredocBeginPattern = regexp.MustCompile(`^<<(\w+)$`)
	heredocEndPattern   = regexp.MustCompile(`^(\w+)$`)
)

func (a *appContext) raise(err error) {
	if !a.prompt {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Printf("error: %s\n", err.Error())
}

func (a *appContext) promptExecutor(input string) {
	if a.chat == nil {
		return
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return
	} else if input == "quit" || input == "exit" {
		bye()
		return
	}

	if a.multiline {
		match := heredocEndPattern.FindStringSubmatch(input)
		if len(match) > 1 && match[1] == a.heredocDelimiter {
			a.flush()
			return
		}
		a.buffer = fmt.Sprintf("%s\n%s", a.buffer, input)
		return
	} else {
		match := heredocBeginPattern.FindStringSubmatch(input)
		if len(match) > 1 {
			a.multiline = true
			a.heredocDelimiter = match[1]
			return
		}
	}
	a.buffer = input
	a.flush()
}

func (a *appContext) getPromptPrefix() (string, bool) {
	if a.multiline {
		return MultiLinePrefix, true
	}
	return DefaultPrefix, false
}

func (a *appContext) flush() {
	var (
		output   string
		response openai.ChatCompletionResponse
		err      error
	)
	a.heredocDelimiter = ""
	a.multiline = false
	if a.jsonOutput {
		response, err = a.chat.ChatRequest(a.buffer)
		if err != nil {
			a.raise(err)
			return
		}
		output, err = response.ToJSON()
		if err != nil {
			a.raise(err)
			return
		}
	} else {
		output, err = a.chat.Chat(a.buffer)
	}
	if err != nil {
		a.raise(err)
		return
	}
	fmt.Println(strings.TrimSpace(output))
	a.buffer = ""
}

func (a *appContext) promptCompleter(prompt.Document) []prompt.Suggest {
	return nil
}

var ccExitKeyBind = prompt.KeyBind{
	Key: prompt.ControlC,
	Fn: func(b *prompt.Buffer) {
		bye()
	},
}

func (a *appContext) startPrompt() {
	defer handleExit()
	a.p = prompt.New(
		a.promptExecutor,
		a.promptCompleter,
		prompt.OptionPrefix(DefaultPrefix),
		prompt.OptionLivePrefix(a.getPromptPrefix),
		prompt.OptionTitle(a.name),
		prompt.OptionAddKeyBind(ccExitKeyBind),
	)
	a.p.Run()
}

func (a *appContext) runOnce() {
	if a.pipe {
		input := a.readStdin()
		a.buffer = fmt.Sprintf("%s\n%s", a.buffer, input)
	}
	a.flush()
}
