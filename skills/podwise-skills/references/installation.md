---
name: podwise-installation
description: "Instructions for installing and configuring the Podwise CLI. Load this when the user needs to set up Podwise for the first time or troubleshoot their installation."
---

# Podwise Installation

Follow these steps to install and configure the `podwise` CLI.

## Step 1: Install Podwise

### macOS

```bash
brew install podwiseai/tap/podwise
```

### Linux

```bash
curl -fsSL https://install.podwise.ai | sh
```

### Via npm

```bash
npm install -g @podwise/cli
```

## Step 2: Verify Installation

```bash
podwise --help
```

## Step 3: Configure API Key

1. Sign up at [podwise.ai](https://podwise.ai) and obtain your API key
2. Configure the CLI:

```bash
podwise config set api-key YOUR_API_KEY
```

## Step 4: Verify Configuration

```bash
podwise config show
```

You should see your configured API key and default settings.

## Troubleshooting

- If `podwise --help` fails, ensure the installation path is in your `$PATH`
- If API key is invalid, re-run `podwise config set api-key` with a fresh key from podwise.ai
- For full documentation, visit [docs.podwise.ai](https://docs.podwise.ai)