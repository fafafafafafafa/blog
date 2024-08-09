#!/usr/bin/env sh

# 将模板 conf 配置文件注入对应环境变量, 生成到指定文件夹
# -e 表示如果任何命令执行失败（返回非零状态），脚本将立即退出。-u 表示如果有未设置的变量被引用，脚本会抛出错误并退出。这样增加了脚本的稳定性和安全性。
set -eu 
if [ "$USE_HTTPS" == "true" ]; then
	# 在 < 和 > 之间，脚本使用 envsubst 命令替换模板文件中的变量。
	# ${SERVER_NAME}, ${BACKEND_HOST}, 和 ${SERVER_PORT} 这些环境变量会从模板文件中被替换成它们的实际值。
	# 这些环境变量通常在部署时被设置，以指定服务器的名称、后端主机地址和服务器端口号等信息。
	# 然后生成的配置文件将被保存到 /etc/nginx/conf.d/default.conf
	envsubst '${SERVER_NAME} ${BACKEND_HOST} ${SERVER_PORT}' < /etc/nginx/conf.d/default.conf.ssl.template > /etc/nginx/conf.d/default.conf
else
	envsubst '${SERVER_NAME} ${BACKEND_HOST} ${SERVER_PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
fi
# rm /etc/nginx/conf.d/default.conf.template
# rm /etc/nginx/conf.d/default.conf.ssl.template

#  执行传递给脚本的任何参数（通常是 Nginx 服务器的启动命令）。这允许你在生成配置文件后启动 Nginx 服务器或者其他相关的服务。
exec "$@"