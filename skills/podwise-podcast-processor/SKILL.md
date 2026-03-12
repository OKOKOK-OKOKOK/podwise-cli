---
name: podwise-podcast-processor
description: 使用 podwise 对播客和 Youtube 视频进行转录总结处理的全流程：搜索节目、处理播客/YouTube/小宇宙链接、轮询处理状态、获取 transcript/summary/chapters/Q&A/mind map/highlights/keywords，并导出结构化结果。用户提出“处理播客”“整理节目内容”“提取字幕”“生成摘要”“YouTube 转文字”“小宇宙总结”“podcast transcript/summary”等请求时使用。
license: MIT
metadata:
  author: Podwise
  version: "1.0"
---

# Podwise 播客节目处理助手

使用这个 skill，把原始节目链接转成可直接消费的结构化内容。

## 工作目标

1. 先确认 `podwise` 可用且 API Key 有效。
2. 根据输入选择 `search`、`process` 或 `get` 路径。
3. 只有`process`退出, 且状态到 `done` 才拉取正文结果。
4. 输出时始终带上 episode URL 和状态说明。

## 第一步：环境检查

运行：

```bash
podwise --help
podwise config show
```

## 第二步：选择流程

- 用户只给关键词或节目名：运行 `podwise search`。
- 用户给 YouTube / 小宇宙链接：运行 `podwise process <url>`（会自动导入）。
- 用户给 Podwise episode URL：如果不确定已处理完成，先 `podwise process <episode-url>`。
- 用户只要已知节目的某类结果：直接 `podwise get <type> <episode-url>`。

## 第三步：执行命令

### 搜索节目

```bash
podwise search "Hard Fork"
podwise search "AI agent" --json
```

需要给下游程序解析时，使用 `--json`。

### 处理节目或视频

```bash
podwise process https://podwise.ai/dashboard/episodes/7360326
podwise process https://www.youtube.com/watch?v=d0-Gn_Bxf8s
podwise process https://youtu.be/d0-Gn_Bxf8s
podwise process https://www.xiaoyuzhoufm.com/episode/abc123
```

`process` 过程会持续自动轮询处理进度和状态，直到结束才退出返回。

### 获取 AI 结果

```bash
podwise get transcript <episode-url>
podwise get summary <episode-url>
podwise get qa <episode-url>
podwise get chapters <episode-url>
podwise get mindmap <episode-url>
podwise get highlights <episode-url>
podwise get keywords <episode-url>
```

注意：`get` 获取结果只能接受 podwise episode url 作为参数，不能是 Youtube 和小宇宙链接。`get` 结果将直接打印到标准输出。

## 中文请求到命令的映射

- “帮我处理这个 YouTube 并输出字幕和摘要”：`process` + `get transcript` + `get summary`。
- “按主题找几期播客”：`search "<主题>" --limit <n>`。
- “导出字幕文件”：`get transcript --format srt` 或 `--format vtt`。
- “给我结构化复盘”：`get summary` + `get chapters` + `get highlights` + `get keywords`。

## 一键脚本（推荐）

使用 [scripts/podwise_pipeline.sh](scripts/podwise_pipeline.sh) 自动完成：

1. 处理 URL (URL 可以是 podwise episode URL, Youtube video URL 和 Xiaoyuzhou episode URL)。
2. 导出 transcript/summary/chapters/qa/mindmap/highlights/keywords。

运行：

```bash
bash scripts/podwise_pipeline.sh "<episode-url or video-url>" "<output-dir>"
```

默认输出目录：`./podwise-output`。

## 常见错误处理

- 找不到 `podwise` cli 工具或未正确配置，提示用户并直接终止流程执行。安装文档： https://github.com/hardhackerlabs/podwise-cli

## 输出格式约定

返回结果时至少包含：

1. 解析后的 episode URL。
2. 当前处理状态。
3. 用户请求的内容块（summary/transcript 等）。
4. 缺失内容明确标记为 unavailable。

需要完整命令清单时，读取 [references/commands.md](references/commands.md)。
