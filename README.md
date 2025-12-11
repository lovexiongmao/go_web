# Go Web 后端服务

基于 Gin + GORM + Dig + Audit + Logrus 的 Web 后端服务框架。

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **依赖注入**: Dig (Uber)
- **审计**: HTTP 审计中间件 + GORM 数据库审计插件
- **日志**: Logrus
- **数据库**: MySQL
- **API文档**: Swagger/OpenAPI

## 项目结构

```
.
├── cmd/
│   └── server/          # 服务入口
│       └── main.go
├── docs/
│   └── swagger/         # Swagger API 文档
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接和审计插件
│   ├── handler/         # HTTP处理器
│   ├── logger/          # 日志模块
│   ├── middleware/      # 中间件（日志、审计）
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── router/          # 路由配置
│   ├── service/         # 业务逻辑层
│   └── util/            # 工具函数（统一响应格式等）
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

### Swagger API 文档

项目集成了 Swagger UI，启动服务后访问：

```
http://localhost:8080/swagger/index.html
```

Swagger 文档提供了完整的 API 接口说明和在线测试功能。

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

### 核心功能

- ✅ RESTful API 设计
- ✅ 依赖注入（Dig）
- ✅ **数据库连接池管理** - 自动配置连接池参数，优化数据库性能
- ✅ 请求日志记录（Logrus）
- ✅ 优雅关闭服务
- ✅ 自动数据库迁移
- ✅ **环境变量配置** - 支持通过 `.env` 文件配置所有参数

### 审计功能

项目实现了**双层审计机制**：

1. **HTTP 层面审计**（中间件）
   - 记录所有 API 请求信息
   - 包括：请求方法、路径、客户端 IP、User-Agent、用户 ID
   - 记录到日志文件

2. **数据库层面审计**（GORM 插件）
   - 自动记录所有表的增删改操作
   - 记录操作前后的数据变化（JSON 格式）
   - 记录操作者信息（用户 ID、IP）
   - 存储在 `audit_logs` 表中
   - **适用于所有模型和表**，无需额外配置

### API 文档

- ✅ Swagger/OpenAPI 集成
- ✅ 自动生成 API 文档
- ✅ 在线 API 测试界面

### 统一响应格式

- ✅ 统一的 HTTP 响应结构
- ✅ 标准化的错误处理
- ✅ 支持自定义响应消息

## 开发说明

### 添加新的API

1. 在 `internal/model/` 中定义数据模型
2. 在 `internal/repository/` 中实现数据访问层
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/handler/` 中实现HTTP处理器（添加 Swagger 注释）
5. 在 `internal/router/router.go` 中注册路由
6. 在 `pkg/dig/container.go` 中注册依赖
7. 运行 `swag init` 重新生成 Swagger 文档

### Swagger 文档

项目使用 [swaggo/swag](https://github.com/swaggo/swag) 生成 API 文档。

**生成文档**：
```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go
```

**访问文档**：
启动服务后访问 `http://localhost:8080/swagger/index.html`

**添加 API 注释示例**：
```go
// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "用户信息"
// @Success      201   {object}  util.Response{data=UserResponse}
// @Failure      400   {object}  util.Response
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### 数据库连接池管理

项目在 `internal/database/database.go` 中实现了数据库连接池的配置：

```go:30:38:internal/database/database.go
// 配置连接池
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}

sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数：10
sqlDB.SetMaxOpenConns(100)       // 最大打开连接数：100
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间：1小时
```

**配置说明**：
- `SetMaxIdleConns(10)`: 设置连接池中空闲连接的最大数量，保持一定数量的连接可以快速复用
- `SetMaxOpenConns(100)`: 设置数据库打开连接的最大数量，防止连接数过多导致数据库压力过大
- `SetConnMaxLifetime(time.Hour)`: 设置连接的最大生存时间，超过时间的连接会被关闭并重新创建，避免长时间连接导致的网络问题

这些配置可以有效优化数据库连接的性能和稳定性。

### 环境变量配置

项目使用 `godotenv` 库支持通过 `.env` 文件配置所有参数，配置加载逻辑在 `internal/config/config.go` 中：

```go:39:64:internal/config/config.go
func LoadConfig() (*Config, error) {
    // 加载.env文件（如果存在）
    _ = godotenv.Load()

    config := &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Mode: getEnv("GIN_MODE", "debug"),
        },
        Database: DatabaseConfig{
            Host:        getEnv("DB_HOST", "localhost"),
            Port:        getEnv("DB_PORT", "3306"),
            User:        getEnv("DB_USER", "root"),
            Password:    getEnv("DB_PASSWORD", ""),
            DBName:      getEnv("DB_NAME", "testdb"),
            AutoMigrate: getEnv("DB_AUTO_MIGRATE", "true") == "true",
        },
        Log: LogConfig{
            Level:     getEnv("LOG_LEVEL", "info"),
            Format:    getEnv("LOG_FORMAT", "text"),
            Output:    getEnv("LOG_OUTPUT", "stdout"),
            LogFile:   getEnv("APP_LOG_FILE", "logs/app.log"),
            AuditFile: getEnv("AUDIT_LOG_FILE", "logs/audit.log"),
        },
    }
    // ...
}
```

**支持的环境变量**：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `SERVER_PORT` | 服务端口 | `8080` |
| `SERVER_HOST` | 服务地址 | `0.0.0.0` |
| `GIN_MODE` | Gin 模式（debug/release/test） | `debug` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | 空 |
| `DB_NAME` | 数据库名称 | `testdb` |
| `DB_AUTO_MIGRATE` | 是否自动迁移 | `true` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `LOG_FORMAT` | 日志格式（json/text） | `text` |
| `LOG_OUTPUT` | 日志输出（stdout/file/both） | `stdout` |
| `APP_LOG_FILE` | 应用日志文件路径 | `logs/app.log` |
| `AUDIT_LOG_FILE` | 审计日志文件路径 | `logs/audit.log` |

**使用方式**：
1. 创建 `.env` 文件（项目根目录）
2. 设置需要的环境变量，例如：
   ```bash
   SERVER_PORT=8080
   DB_HOST=localhost
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=testdb
   ```
3. 项目启动时会自动加载 `.env` 文件中的配置
4. 如果环境变量未设置，会使用代码中的默认值

**注意**：`.env` 文件通常不应提交到版本控制系统，建议添加到 `.gitignore` 中。

### 日志级别

- `debug`: 详细调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 审计日志

#### HTTP 层面审计（中间件）

所有 API 请求都会自动记录审计日志到日志文件，包括：
- 请求方法
- 请求路径
- 客户端 IP
- User-Agent
- 用户 ID（如果已认证）

#### 数据库层面审计（GORM 插件）

所有数据库的增删改操作都会自动记录到 `audit_logs` 表，包括：
- 表名（`table_name`）
- 记录 ID（`record_id`）
- 操作类型（`action`: create/update/delete）
- 操作前的数据（`old_values`，JSON 格式）
- 操作后的数据（`new_values`，JSON 格式）
- 操作者用户 ID（`user_id`）
- 操作者 IP（`ip`）
- 操作时间（`created_at`）

**特点**：
- 自动适用于所有模型和表，无需额外配置
- 使用独立的数据库连接，不影响原事务
- 自动跳过 `audit_logs` 表自身的操作，避免递归
- 通过 Context 传递用户信息，支持在 Handler/Service 层设置

**使用示例**：

```go
// 在 Handler 中，用户信息已通过中间件自动设置到 Context
// 数据库操作会自动记录审计日志

// 创建用户
db.Create(&user)  // ✅ 自动记录创建审计日志

// 更新用户
db.Save(&user)    // ✅ 自动记录更新审计日志（包含旧值和新值）

// 删除用户
db.Delete(&user)  // ✅ 自动记录删除审计日志（包含旧值）
```

### 统一响应格式

项目提供了统一的 HTTP 响应格式（`internal/util/response.go`）：

```go
// 成功响应
util.Success(c, data)
util.SuccessWithMessage(c, "自定义消息", data)

// 创建成功响应
util.Created(c, data)
util.CreatedWithMessage(c, "创建成功", data)

// 错误响应
util.BadRequest(c, "错误消息")
util.NotFound(c, "资源不存在")
util.InternalServerError(c, "服务器错误")
```

响应格式：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": { ... }
}
```

## 许可证

MIT

