# Progress Wall 后端开发指南

## 技术栈

- **Go 1.21+**: 主要编程语言
- **Gin**: Web框架
- **GORM**: ORM框架
- **MySQL/SQLite**: 数据库支持

## 项目结构

```
backend/
├── config/                 # 配置管理
│   └── config.go
├── controllers/            # 控制器层
│   └── ...
├── database/              # 数据库连接
|   |-  seed.go            # 初始权限组设定
|   |-  migrate.go         # 数据库迁移（自动建表）
│   └── database.go        # 数据库连接
├── middleware/            # 中间件
│   └── ...
├── models/                # 数据模型
│   └── ...
├── repository/            # 数据访问层
│   └── ...
├── services/              # 业务逻辑层
│   └── ...
├── router/                # 路由配置
│   └── ...
├── go.mod                 # Go模块文件
├── main.go               # 应用入口
└── config.env.example    # 环境变量示例
```

## 架构设计

### 分层架构

1. **Models层**: 定义数据结构和数据库映射
2. **Repository层**: 数据访问抽象，封装数据库操作
3. **Service层**: 业务逻辑处理，包含核心业务规则
4. **Controller层**: HTTP请求处理，参数验证和响应格式化
5. **Router层**: 路由配置和中间件设置

## 开发环境设置

### 1. 环境要求

- Go 1.21 或更高版本
- MySQL 5.7+ 或 SQLite 3

### 2. 安装依赖

```bash
cd backend
go mod tidy
```

### 3. 环境配置

复制环境变量示例文件：

```bash
cp config.env.example .env
```

编辑 `.env` 文件，配置数据库和服务器参数。

### 4. 运行应用

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

### 5. 构建应用/交叉编译

```bash
go build ./
```

交叉编译请参考： https://www.topgoer.com/%E5%85%B6%E4%BB%96/%E8%B7%A8%E5%B9%B3%E5%8F%B0%E4%BA%A4%E5%8F%89%E7%BC%96%E8%AF%91.html

## API 接口文档

详细的API接口文档请参考：[API.md](./docs/API.md)

文档包含：
- 认证相关接口（注册、登录）
- 用户相关接口
- 看板相关接口（CRUD）
- 列相关接口（CRUD）
- 任务相关接口（CRUD）
- 任务拖拽排序接口

所有接口遵循RESTful规范，支持完整的增删改查操作。

## 开发规范

### 代码组织

1. **单一职责**: 每个文件只负责一个功能模块
2. **接口抽象**: 使用接口定义服务契约
3. **依赖注入**: 通过构造函数注入依赖
4. **错误处理**: 统一的错误处理和响应格式

### 命名规范

- **包名**: 小写，简短有意义
- **接口名**: 以 `er` 结尾，如 `UserRepository`
- **结构体**: 大驼峰命名
- **方法名**: 大驼峰命名，动词开头

### 文件大小控制

- 单个文件不超过200行
- 复杂功能拆分为多个文件
- 保持代码模块化和可维护性

## 故障排除

### 常见问题

1. **数据库连接失败**: 检查数据库配置和连接字符串
2. **JWT认证失败**: 检查JWT_SECRET配置
3. **CORS错误**: 检查CORS_ALLOW_ORIGINS配置
4. **端口占用**: 检查SERVER_PORT是否被占用

### 调试技巧

- 使用 `SERVER_MODE=debug` 查看详细日志
- 检查数据库连接和表结构
- 使用Postman或curl测试API接口

## 贡献指南

1. 遵循Go代码规范
2. 保持文件大小在200行以内
3. 添加适当的注释和文档
4. 编写单元测试 (推荐)
5. 提交前运行 `go fmt` 和 `go vet`