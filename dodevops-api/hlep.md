  1. 编译项目

  # 本地编译（Mac/Windows开发环境）
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o devops main.go
  # 构建镜像
  docker build -t crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/deviops-api:v1.0 .
# 推送镜像
  docker push  crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/deviops-api:v1.0


  # 或在Ubuntu系统上直接编译
  go build -o dodevops-api-linux main.go

  2. 准备配置文件

  将config.yaml配置文件单独放置，不打包进二进制文件。

  3. 启动方式

  指定配置文件启动：
  ./dodevops-api-linux -c /path/to/config.yaml

  使用默认配置文件启动（当前目录下的config.yaml）：
  ./dodevops-api-linux

  查看帮助信息：
  ./dodevops-api-linux -h

  4. Ubuntu 20.4系统部署示例

  # 创建部署目录
  mkdir -p /opt/dodevops-api
  cd /opt/dodevops-api

  # 上传文件
  # - dodevops-api-linux (二进制文件)
  # - config.yaml (配置文件)

  # 赋予执行权限
  chmod +x dodevops-api-linux

  # 启动服务
  ./dodevops-api-linux -c config.yaml

  # 或者创建systemd服务
  sudo tee /etc/systemd/system/dodevops-api.service > /dev/null <<EOF
  [Unit]
  Description=DevOps API Service
  After=network.target mysql.service redis.service

  [Service]
  Type=simple
  User=www-data
  WorkingDirectory=/opt/dodevops-api
  ExecStart=/opt/dodevops-api/dodevops-api-linux -c /opt/dodevops-api/config.yaml
  Restart=always
  RestartSec=5

  [Install]
  WantedBy=multi-user.target
  EOF

  # 启动服务
  sudo systemctl daemon-reload
  sudo systemctl enable dodevops-api
  sudo systemctl start dodevops-api


#go mod init admin-api
#go get github.com/gin-gonic/gin@v1.8.1
#GOPROXY=https://goproxy.io
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/sirupsen/logrus
go get github.com/lestrrat-go/file-rotatelogs
go get github.com/rifflock/lfshook
go get github.com/go-redis/redis/v8@v8.11.5
go get github.com/mojocn/base64Captcha@v1.3.1
go get github.com/dgrijalva/jwt-go
go get gopkg.in/yaml.v3
go get -u github.com/wenlng/go-user-agent
go get github.com/gogf/gf
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger

 go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o devops main.go
lsof -i:8000 | grep LISTEN | awk '{print $2}' | xargs kill -9
 /Users/apple/go/bin/swag  init
 go run main.go

