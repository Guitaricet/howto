# Howto
Howto is a GPT-3/Codex-powered shell tool that allows you to talk with your shell in natural language.


[![Howto demo](https://user-images.githubusercontent.com/2821124/197363019-b038d31e-fde0-45a5-b347-3e87e0c260a6.gif)](https://youtu.be/VwP9eZdTrGY)

Forgot how to create a conda environment?
```bash
% howto create conda env
conda create -n <env_name> python=3.6
```

Forgot how to add a new env to Jupyter?
```bash
% howto add kernel to jupyter
python -m ipykernel install --user --name=
```

Want to download the biggest Rick Astley's hit?
```bash
% howto download youtube video for never give you up
youtube-dl -f 18 https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

It can also suggest how to be a nicer person
```bash
 % howto be a nicer person
alias please='sudo'
```

It works by sending requests to [OpenAI API](http://openai.com/api/) and requires you to setup your own `OPENAI_API_TOKEN`.

# Installation

## Option 1: Download the binary from Github

Download the binary from the [releases page](https://github.com/Guitaricet/howto/releases).
Then move the binary to your path, e.g., `mv howto /usr/local/bin/`

## Option 2: Build from source

If you have Go installed, you can build the binary from source.

```bash
go build
```
> if you have your `$GOPATH/bin` in your path, just run `go install .` to install the binary

Then move the binary to your path, e.g., `mv howto /usr/local/bin/`

## Environment variables

You need to connect your OpenAI API key to the program by setting the `OPENAI_API_KEY` environment variable. Get your OpenAI API key [here](https://beta.openai.com/docs/quickstart/add-your-api-key).

```bash
OPENAI_API_KEY=<your_api_key>
```

By default we use `text-davinci-002`, you can change it to a different model by setting the `HOWTO_OPENAI_MODEL` environment variable. It's best to use Codex models (e.g., `code-davinci-002`), but they are currently in beta and not available to everyone.

```bash
HOWTO_OPENAI_MODEL=code-davinci-002
```
