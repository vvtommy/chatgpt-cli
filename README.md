# chatgpt-cli

This is a command-line tool for chat completions, currently only supporting GPT-3 Turbo.

## Example

Simply type chatgpt and hit enter to start the conversation

![截屏2023-03-07 00 12 29](https://user-images.githubusercontent.com/743350/223187939-9a0a3744-acda-47e9-aa28-f40d05ba5742.png)

If you wanna use the ChatGPT command-line tool like a pro, just pipe your input into it and you're good to go.

```bash
cat ./awsome-prompt.txt | chatgpt --json | jq -r '.usage.total_tokens'
```

even instead of `jq`

```bash
$ cat ./awsome-prompt.txt | chatgpt --json | \
  chatgpt 'Calculate the amount spent on tokens, given that you acted as a JSON parser and \
    the cost is $0.0002 per 1k tokens. The total number of tokens used can be determined from \
    the attribute .usage.total_tokens'
    
The total tokens used are 44.
Cost per 1k tokens is $0.0002.
So, the cost for 44 tokens can be calculated as:
(44/1000) * $0.0002 = $0.0000088
Therefore, the amount spent on tokens is $0.0000088.
```

Use heredoc to support multiple-line input.

![截屏2023-03-07 01 35 56](https://user-images.githubusercontent.com/743350/223187997-7ca35234-32b5-4332-a43d-96bca7ed1124.png)

proxy support 🪜

```
$ HTTPS_PROXY=127.0.0.1:7890 chatgpt
```

## Usage

```bash
$ chatgpt --help
ChatGPT command line tool that supports pipe and repl.

Usage:
  chatgpt [flags]

Flags:
  -h, --help   help for chatgpt
  -j, --json   output as json
```
