#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "Usage: $0 <podwise-episode-url> [output-dir]" >&2
  exit 1
fi

episode_url="$1"
output_dir="${2:-./podwise-output}"
process_log="${output_dir}/process.log"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

require_cmd podwise
mkdir -p "$output_dir"

echo "Processing: $episode_url"
podwise process "$episode_url" --interval 30s --timeout 45m | tee "$process_log"

echo "Fetching summary -> ${output_dir}/summary.md"
podwise get summary "$episode_url" >"${output_dir}/summary.md"

echo "Fetching transcript (text/srt/vtt)"
podwise get transcript "$episode_url" --format text >"${output_dir}/transcript.txt"
podwise get transcript "$episode_url" --format srt >"${output_dir}/transcript.srt"
podwise get transcript "$episode_url" --format vtt >"${output_dir}/transcript.vtt"

echo "Fetching chapters/qa/mindmap/highlights/keywords"
podwise get chapters "$episode_url" >"${output_dir}/chapters.md"
podwise get qa "$episode_url" >"${output_dir}/qa.md"
podwise get mindmap "$episode_url" >"${output_dir}/mindmap.md"
podwise get highlights "$episode_url" >"${output_dir}/highlights.md"
podwise get keywords "$episode_url" >"${output_dir}/keywords.md"

echo
echo "Done. Files written to: $output_dir"
echo "Process log: $process_log"
