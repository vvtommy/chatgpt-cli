package main

import (
	"bufio"
	"chatgit-cli/internal/chat"
	"chatgit-cli/openai"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/viper"
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
	openai     *openai.Context

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
		err    error
		reader *bufio.Reader
		input  string
	)

	reader = bufio.NewReader(os.Stdin)
	input, err = reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return input
}

func (a *appContext) run() {
	if a.prompt {
		a.startPrompt()
		return
	}
	a.runOnce()
}

func bye() {
	fmt.Println("Bye!")
	os.Exit(0)
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
		b.InsertText("Ctrl+C", false, false)
		bye()
	},
}

func (a *appContext) startPrompt() {
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
		output, err := a.chat.Chat(input)
		if err != nil {
			a.raise(err)
			return
		}
		a.buffer = output
	}
	a.flush()
}
