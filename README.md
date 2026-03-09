# podwise-cli

CLI client for [podwise.ai](https://podwise.ai) — turn any podcast episode into AI-powered insights, designed for use in AI agents and skills workflows.

Podwise transforms hours of podcasts into summaries, outlines, transcripts, Q&A, and mind maps. This CLI is purpose-built as a **tool for AI agents** — letting LLMs, skills runtimes, and automation pipelines fetch structured podcast insights without a browser or human in the loop.

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

### Examples

```bash
# Fetch transcript in Chinese for downstream processing
podwise get https://podwise.ai/dashboard/episodes/7360326 --type transcript

# Get Q&A pairs — useful for RAG ingestion
podwise get <url> --type qa

# Save mind map as xmind file
podwise get <url> --type mindmap --format xmind --output ./notes/
```