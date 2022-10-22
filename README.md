# howto
How to do bash commands you always forget. OpenAI-powered.

# Installation

## Build from source

```bash
go build ./howto.go
```

Then you can move the binary to your path, e.g., `mv howto /usr/local/bin/`

## Environment variables

You need to connect your OpenAI API key to the program. You can do this by setting the `OPENAI_API_KEY` environment variable.

```bash
OPENAI_API_KEY=<your_api_key>
```

By default we use `text-davinci-002`, you can change it to a different model by setting the `HOWTO_OPENAI_MODEL` environment variable. It's best to use Codex models (e.g., `code-davinci-002`), but they are currently in beta and not available to everyone.

```bash
HOWTO_OPENAI_MODEL=code-davinci-002
```
