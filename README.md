# podwise-cli

CLI client for [podwise.ai](https://podwise.ai) — turn any podcast episode into AI-powered insights, designed for use in AI agents and skills workflows.

Podwise transforms hours of podcasts into summaries, outlines, transcripts, Q&A, and mind maps. This CLI is purpose-built as a **tool for AI agents** — letting LLMs, skills runtimes, and automation pipelines fetch structured podcast insights without a browser or human in the loop.

## Why a CLI?

Most podcast insight tools assume a human is navigating a UI. `podwise-cli` is different: it's designed as a composable primitive for **AI agent workflows**.

- Drop it into a [Cursor skill](https://docs.cursor.com), [Codex skill](https://github.com/hardhacker/codex), or any agent that can invoke shell tools
- Pipe clean Markdown output directly into an LLM prompt, a RAG pipeline, or a note-taking tool
- Let agents autonomously enrich context with podcast knowledge — no browser, no manual copy-paste

## Installation

```bash
# Build from source
go build -o podwise .

# Or use goreleaser snapshot build
goreleaser release --snapshot --clean
```

## Configuration

```bash
# Set your podwise.ai API key
podwise config set api_key sk-xxxx

# Verify config
podwise config show
```

The config file lives at `~/.config/podwise/config.yaml`.

## Usage

```
podwise get <episode-url> [flags]
```

### Flags

```
--type    string   output type: summary | outline | transcript | qa | mindmap  (default: summary)
--lang    string   output language: en | zh | ja | ko | fr | de | es | pt      (default: episode language)
--output  string   write to a file or directory instead of stdout
--export  string   export to: notion | obsidian | readwise | logseq
--format  string   file format: md | pdf | srt | xmind                         (default: md)
```

### Examples

```bash
# Print AI summary to stdout (pipe into an LLM prompt)
podwise get https://podwise.ai/dashboard/episodes/7360326

# Fetch transcript in Chinese for downstream processing
podwise get <url> --type transcript --lang zh

# Get Q&A pairs — useful for RAG ingestion
podwise get <url> --type qa

# Save mind map as xmind file
podwise get <url> --type mindmap --format xmind --output ./notes/

# Export summary directly to Obsidian
podwise get <url> --export obsidian
```

## Using with AI Agents & Skills

`podwise-cli` is designed to be called as a shell tool inside agent runtimes. Example skill usage:

```bash
# Inside a Cursor or Codex skill — fetch episode insights and feed to the agent
podwise get <episode-url> --type summary
podwise get <episode-url> --type transcript | your-rag-ingest-script
```

Because all output goes to stdout as plain Markdown by default, agents can consume it directly without parsing or post-processing.
