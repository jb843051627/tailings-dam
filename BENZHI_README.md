# 尾矿库安全监测系统 (Tailings Dam Safety Monitoring System)

## 项目简介

尾矿库安全监测系统是一个用于监测和管理尾矿坝安全状态的 Web 服务。系统提供尾矿坝信息管理、监测点管理、渗流/位移/孔隙水压力读数采集、告警管理、巡检管理和排水系统管理等功能。

## 技术栈

- **语言**: Go 1.22
- **数据库**: SQLite (modernc.org/sqlite, 纯 Go 实现，无需 CGO)
- **架构**: 分层架构 (Model → Store → Service → Handler)

## 项目结构

```
tailings-dam/
├── main.go                      # 程序入口
├── go.mod / go.sum              # Go 模块文件
├── benzhi.Dockerfile            # Docker 构建文件
├── build_benzhi_docker.sh       # Docker 构建脚本
├── BENZHI_README.md             # 本文档
├── .gitignore
├── .dockerignore
└── internal/
    ├── model/                   # 数据模型层
    │   ├── dam.go               # 尾矿坝模型
    │   ├── monitoring_point.go  # 监测点模型
    │   ├── reading.go           # 读数模型（渗流/位移/孔隙水压力）
    │   ├── alert.go             # 告警模型
    │   ├── inspection.go        # 巡检模型
    │   └── drainage.go          # 排水系统模型
    ├── store/                   # 数据存储层
    │   ├── store.go             # Store 结构体和数据库初始化
    │   ├── dam_store.go         # 尾矿坝 CRUD
    │   ├── monitoring_point_store.go  # 监测点 CRUD
    │   ├── reading_store.go     # 读数 CRUD 和批量导入
    │   ├── alert_store.go       # 告警 CRUD
    │   ├── inspection_store.go  # 巡检 CRUD
    │   └── drainage_store.go    # 排水系统 CRUD
    ├── service/                 # 业务逻辑层
    │   ├── dam_service.go       # 尾矿坝业务逻辑
    │   ├── monitoring_service.go     # 监测点业务逻辑
    │   ├── reading_service.go   # 读数业务逻辑
    │   ├── alert_service.go     # 告警业务逻辑和阈值检查
    │   ├── inspection_service.go # 巡检业务逻辑
    │   └── drainage_service.go   # 排水系统业务逻辑
    └── handler/                 # HTTP 处理层
        ├── router.go            # 路由配置和中间件
        ├── dam_handler.go       # 尾矿坝 HTTP 处理器
        ├── monitoring_handler.go     # 监测点 HTTP 处理器
        ├── reading_handler.go   # 读数 HTTP 处理器
        ├── alert_handler.go     # 告警 HTTP 处理器
        ├── inspection_handler.go # 巡检 HTTP 处理器
        ├── drainage_handler.go  # 排水系统 HTTP 处理器
        └── dashboard_handler.go # 仪表盘 HTTP 处理器
```

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/dams | 列出所有尾矿坝 |
| POST | /api/dams | 创建尾矿坝 |
| GET | /api/dams/{id} | 获取尾矿坝详情 |
| PUT | /api/dams/{id} | 更新尾矿坝 |
| GET | /api/monitoring-points | 列出监测点 |
| POST | /api/monitoring-points | 创建监测点 |
| GET | /api/readings/seepage | 列出渗流读数 |
| POST | /api/readings/seepage | 创建渗流读数 |
| GET | /api/readings/displacement | 列出位移读数 |
| POST | /api/readings/displacement | 创建位移读数 |
| GET | /api/readings/pore-pressure | 列出孔隙水压力读数 |
| POST | /api/readings/pore-pressure | 创建孔隙水压力读数 |
| POST | /api/readings/batch | 批量导入读数 |
| GET | /api/alerts | 列出告警 |
| PUT | /api/alerts/{id}/resolve | 解决告警 |
| GET | /api/inspections | 列出巡检 |
| POST | /api/inspections | 创建巡检 |
| PUT | /api/inspections/{id}/complete | 完成巡检 |
| GET | /api/drainage | 列出排水系统 |
| GET | /api/dashboard | 仪表盘汇总 |

## 本地运行

```bash
# 设置环境变量
export GOTOOLCHAIN=local
export GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod tidy

# 编译
go build -o tailings-dam .

# 运行
./tailings-dam

# 或指定端口和数据库路径
PORT=9090 DB_PATH=/tmp/test.db ./tailings-dam
```

## Docker 运行

```bash
# 构建镜像
bash build_benzhi_docker.sh

# 运行容器
docker run -p 8080:8080 -v tailings-dam-data:/data tailings-dam:latest
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | HTTP 监听端口 |
| DB_PATH | tailings-dam.db | SQLite 数据库文件路径 |

## 数据实体

### Dam (尾矿坝)
- 名称、位置、经纬度、容量、当前水位
- 危险等级 (low/medium/high/critical)
- 状态 (active/inactive/closed/overflow_risk/emergency)

### MonitoringPoint (监测点)
- 关联坝体、名称、编号、类型
- 类型: seepage/displacement/pore_pressure/phreatic_surface
- GPS 位置、高程、状态

### Reading (读数)
- SeepageReading: 流量、浊度、pH、温度
- DisplacementReading: 水平/垂直位移、累计位移
- PorePressureReading: 压力、深度、水位

### Alert (告警)
- 等级 (info/warning/danger/critical)
- 状态 (active/acknowledged/resolved/suppressed)
- 阈值和当前值

### Inspection (巡检)
- 巡检人、计划日期、完成日期
- 状态 (pending/in_progress/completed/overdue/cancelled)
- 优先级 (low/normal/high/urgent)

### DrainageSystem (排水系统)
- 类型 (main/secondary/emergency/seepage_collection/discharge)
- 状态 (normal/blocked/overflow/damaged/maintenance/offline)
- 设计流量、实际流量、管径、长度、材质
