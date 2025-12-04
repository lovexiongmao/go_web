# Go Web 后端服务

基于 Gin + GORM + Dig + Audit + Logrus 的 Web 后端服务框架。

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **依赖注入**: Dig (Uber)
- **审计**: 自定义审计中间件
- **日志**: Logrus
- **数据库**: MySQL

## 项目结构

```
.
├── cmd/
│   └── server/          # 服务入口
│       └── main.go
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   ├── handler/         # HTTP处理器
│   ├── logger/          # 日志模块
│   ├── middleware/      # 中间件（日志、审计）
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── router/          # 路由配置
│   └── service/         # 业务逻辑层
├── pkg/
│   └── dig/             # 依赖注入容器
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置数据库连接信息等。

### 3. 创建数据库

在MySQL中创建数据库：

```sql
CREATE DATABASE testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 健康检查

```
GET /health
```

### 用户管理

- `POST /api/v1/users` - 创建用户
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

### 创建用户示例

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "password": "123456"
  }'
```

## 功能特性

- ✅ RESTful API 设计
- ✅ 依赖注入（Dig）
- ✅ 数据库连接池管理
- ✅ 请求日志记录（Logrus）
- ✅ 审计日志中间件
- ✅ 优雅关闭服务
- ✅ 自动数据库迁移
- ✅ 环境变量配置

## 开发说明

### 添加新的API

1. 在 `internal/model/` 中定义数据模型
2. 在 `internal/repository/` 中实现数据访问层
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/handler/` 中实现HTTP处理器
5. 在 `internal/router/router.go` 中注册路由
6. 在 `pkg/dig/container.go` 中注册依赖

### 日志级别

- `debug`: 详细调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 审计日志

所有API请求都会自动记录审计日志，包括：
- 请求方法
- 请求路径
- 客户端IP
- User-Agent

## 许可证

MIT

