# podwise_proxy.ps1
# 设置代理环境变量（只对当前脚本生效）
$env:HTTP_PROXY="http://127.0.0.1:7897"
$env:HTTPS_PROXY="http://127.0.0.1:7897"

# 执行 Podwise 命令，传递所有参数
podwise.exe $args