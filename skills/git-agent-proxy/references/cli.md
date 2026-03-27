# Podwise CLI Reference

Use the `podwise_proxy` command to search episodes, process media, retrieve AI artifacts, manage subscriptions, browse history, and ask questions across podcast transcripts.

Check that the CLI is installed and configured before running any command:

```bash
podwise_proxy --help
podwise_proxy config show
```

Every command supports `--help` for full flag details. Add `--json` for machine-readable output when results will be parsed by another step.

---

## Commands

### search

Find episodes or podcasts by keyword or title.

```bash
# Search episodes
podwise_proxy search episode "AI agents" --limit 10 --json
podwise_proxy search episode "Hard Fork" --limit 5 --json

# Search podcasts by name
podwise_proxy search podcast "Lex Fridman" --limit 5 --json
podwise_proxy search podcast "All-In" --limit 5 --json
```

Use `search episode` when the user wants a list of matching episodes to browse.
Use `search podcast` when the user is looking for a show by name or wants candidates to subscribe to.
Use `ask` — not `search` — when the user wants a synthesized answer from transcript content.

---

### popular

Discover what is currently trending across Podwise.

```bash
podwise_proxy popular --json
```

---

### list

Browse new episodes and followed podcasts.

```bash
# New episodes from followed podcasts — last N days
podwise_proxy list episodes --json --latest 7
podwise_proxy list episodes --json --latest 30

# New episodes by specific date
podwise_proxy list episodes --json --date today
podwise_proxy list episodes --json --date yesterday
podwise_proxy list episodes --json --date 2026-03-01    

# Followed podcasts with recent updates
podwise_proxy list podcasts --json
podwise_proxy list podcasts --json --date today
podwise_proxy list podcasts --json --latest 14
```

---

### drill

List recent episodes for a specific podcast by its Podwise URL.

```bash
podwise_proxy drill https://podwise.ai/dashboard/podcasts/{id} --json
podwise_proxy drill https://podwise.ai/dashboard/podcasts/{id} --latest 7 --json
```

Use `drill` when the user wants to explore a specific show's back catalogue rather than their full followed feed.

---

### follow / unfollow

Manage podcast subscriptions. Both commands are idempotent.

```bash
podwise_proxy follow https://podwise.ai/dashboard/podcasts/{id}
podwise_proxy unfollow https://podwise.ai/dashboard/podcasts/{id}
```

---

### history

Browse the user's listening and reading activity. Results are sorted by most recent first.

```bash
# Episodes the user has read (viewed in Podwise)
podwise_proxy history read --json
podwise_proxy history read --json --limit 20

# Episodes the user has listened to
podwise_proxy history listened --json
podwise_proxy history listened --json --limit 20
```

---

### ask

Ask a question and get a transcript-grounded answer synthesized across Podwise's corpus.

```bash
podwise_proxy ask "the future of AI agents"
podwise_proxy ask "how do founders think about pricing?" --sources
podwise_proxy ask "what do podcasts say about longevity and sleep?"
```

`--sources` appends source excerpts and episode links. Use it when the user wants to trace claims back to specific episodes.

Allow up to 60 seconds — `ask` searches transcripts and generates a synthesis. Do not cancel early.

Do not use `ask` to locate episodes by keyword — use `search episode` for that.

---

### process

Submit an episode, video, or local file for AI processing.

```bash
# Podwise episode URL
podwise_proxy process https://podwise.ai/dashboard/episodes/7360326

# YouTube
podwise_proxy process https://www.youtube.com/watch?v=d0-Gn_Bxf8s
podwise_proxy process https://youtu.be/d0-Gn_Bxf8s

# Xiaoyuzhou
podwise_proxy process https://www.xiaoyuzhoufm.com/episode/abc123

# Local audio or video file
podwise_proxy process ./interview.mp3
podwise_proxy process ./meeting.wav --title "Product Review Meeting"
podwise_proxy process ./demo.mp4 --title "Launch Demo" --hotwords "Podwise,LLM,ASR" 
```

Supported local file extensions: `.mp3 .wav .m4a .mp4 .m4v .mov .webm`

`process` waits for completion before exiting. The output includes the resolved Podwise episode URL — use that URL for all subsequent `get` commands.

**`process` consumes quota/credits. Always confirm with the user before running it**, regardless of input type.

---

### get

Retrieve AI-generated artifacts for a processed episode.

```bash
podwise_proxy get summary <episode-url>
podwise_proxy get transcript <episode-url>
podwise_proxy get chapters <episode-url>
podwise_proxy get highlights <episode-url>
podwise_proxy get qa <episode-url>
podwise_proxy get mindmap <episode-url>
podwise_proxy get keywords <episode-url>
```

Transcript format options:

```bash
podwise_proxy get transcript <episode-url>               # plain text (default)
podwise_proxy get transcript <episode-url> --format srt  # SubRip
podwise_proxy get transcript <episode-url> --format vtt  # WebVTT
```

`get` accepts only a Podwise episode URL — never a YouTube, Xiaoyuzhou, or local file path. Run `process` on external sources first, then use the resulting Podwise URL with `get`.

---

## Artifact Reference

| Artifact | Command | Description |
|---|---|---|
| Summary | `get summary` | AI-generated summary and key takeaways |
| Transcript | `get transcript` | Full transcript in text, SRT, or VTT |
| Chapters | `get chapters` | Chapter breakdown with timestamps |
| Highlights | `get highlights` | Notable highlights with timestamps |
| Q&A | `get qa` | AI-extracted question and answer pairs |
| Mind map | `get mindmap` | Topic structure as a nested tree |
| Keywords | `get keywords` | Topic keywords with descriptions |

---

## Intent → Command Mapping

| User wants to… | Command |
|---|---|
| Find episodes about a topic | `search episode "X"` |
| Find a podcast by name | `search podcast "X"` |
| See what's trending | `popular` |
| See new episodes from followed shows | `list episodes --latest 7` |
| Explore a specific show's episodes | `drill <podcast-url>` |
| Get a synthesized answer from transcripts | `ask "X"` |
| Summarize an episode | `get summary <url>` |
| Get the full transcript | `get transcript <url>` |
| Export subtitles | `get transcript <url> --format srt` |
| Process a YouTube video | confirm → `process <youtube-url>` |
| Transcribe a local file | confirm → `process <file>` |
| Follow a podcast | `follow <podcast-url>` |
| See listening history | `history listened` |
| See reading history | `history read` |

---

## Common Failure Cases

- **`get` returns "not processed"**: Run `podwise_proxy process <url>` first (confirm with user — quota is consumed).
- **`get` passed a non-Podwise URL**: `get` only accepts Podwise episode URLs. Run `process` on external URLs first.
- **Unsupported local file extension**: Supported formats are `.mp3 .wav .m4a .mp4 .m4v .mov .webm`.
- **`ask` returns a plan-limit error**: Report it directly. Do not fabricate an answer.
- **`process` runs without confirmation**: Always wrong — `process` consumes quota.
