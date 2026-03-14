# Podwise CLI Installation and Setup

Use this reference when `podwise` is missing, when the user asks how to install it, or when the CLI is installed but not configured yet.

## Install

### Homebrew (macOS)

```bash
brew tap hardhackerlabs/podwise-tap
brew install podwise
```

### Automatic Install Script

```bash
curl -sL https://raw.githubusercontent.com/hardhackerlabs/podwise-cli/main/install.sh | sh
```

### Manual Binary Install

1. Download the latest binary for the user's OS and architecture from the GitHub Releases page.
2. Unpack the archive, for example:

```bash
tar -xzf podwise_linux_amd64.tar.gz
```

3. Move the binary into a directory on `PATH`, for example:

```bash
mv podwise /usr/local/bin/
```

4. Make sure it is executable:

```bash
chmod +x /usr/local/bin/podwise
```

### Build from Source

```bash
git clone https://github.com/hardhackerlabs/podwise-cli.git
cd podwise-cli
go build -o podwise .
sudo mv podwise /usr/local/bin/
```

## Configure the API Key

Create a Podwise API key in the Podwise dashboard settings page, then run:

```bash
podwise config set api_key your-sk-xxxx
podwise config show
```

The config file is stored at `~/.config/podwise/config.toml`.

## Quick Verification

Run:

```bash
podwise --help
podwise config show
```

If the install is healthy, `podwise --help` should print command usage and `podwise config show` should display the config path and API key status.
