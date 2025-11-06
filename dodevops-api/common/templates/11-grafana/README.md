# Grafana Docker 部署模板

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | grafana |
| GRAFANA_ADMIN_USER | 管理员用户名 | admin |
| GRAFANA_ADMIN_PASSWORD | 管理员密码 | admin123 |
| GRAFANA_PORT | Web端口 | 3000 |
| GRAFANA_ROOT_URL | 根URL | http://localhost:3000 |
| GRAFANA_PLUGINS | 插件列表 | 空 |
| DATA_DIR | 数据目录 | /data/grafana |

## 快速部署

```bash
GRAFANA_ADMIN_PASSWORD=MySecurePass123 docker-compose -f grafana-latest-docker-compose.yml up -d
```

访问: `http://your-ip:3000`

## 常用插件

在环境变量中指定插件：
```bash
GRAFANA_PLUGINS=grafana-piechart-panel,grafana-worldmap-panel
```

推荐插件：
- `grafana-piechart-panel` - 饼图
- `grafana-worldmap-panel` - 世界地图
- `grafana-clock-panel` - 时钟
- `alexanderzobnin-zabbix-app` - Zabbix集成

## 数据源配置

自动配置数据源，创建 `provisioning/datasources/datasource.yml`:

```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
```

## Dashboard 自动导入

创建 `provisioning/dashboards/dashboard.yml`:

```yaml
apiVersion: 1
providers:
  - name: 'default'
    orgId: 1
    folder: ''
    type: file
    options:
      path: /var/lib/grafana/dashboards
```

将 dashboard JSON 文件放入 `${DATA_DIR}/dashboards/` 目录。
