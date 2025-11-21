# Progress Wall API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: Bearer Token (JWT)
- **Content-Type**: `application/json`

## 认证

大部分接口需要认证，请求头中需要包含：

```
Authorization: Bearer <accessToken>
```

获取 token 的方式：通过 `/api/auth/login` 接口登录后获取。

---

## 认证相关接口

### 1. 用户注册

**POST** `/api/auth/register`

**请求体**:
```json
{
  "username": "string (必填, 3-20个字符, 只能包含字母、数字、下划线)",
  "email": "string (必填, 有效邮箱格式)",
  "password": "string (必填, 至少6位, 包含字母和数字)",
  "nickname": "string (可选)"
}
```

**响应** (201 Created):
```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "created_at": "2025-11-20T10:00:00Z"
  }
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误
- `409 Conflict`: 用户名或邮箱已存在

---

### 2. 用户登录

**POST** `/api/auth/login`

**请求体**:
```json
{
  "username": "string (必填, 支持用户名或邮箱)",
  "password": "string (必填)"
}
```

**响应** (200 OK):
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户"
  }
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 用户名或密码错误
- `403 Forbidden`: 账户已被禁用

---

## 用户相关接口

### 3. 获取当前用户信息

**GET** `/api/user/profile`

**需要认证**: 是

**响应** (200 OK):
```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "phone": "",
    "status": 1,
    "created_at": "2025-11-20T10:00:00Z"
  }
}
```

**错误响应**:
- `401 Unauthorized`: 未认证或token无效
- `404 Not Found`: 用户不存在

---

## 看板相关接口

### 4. 获取用户的所有看板

**GET** `/api/boards`

**需要认证**: 是

**响应** (200 OK):
```json
{
  "boards": [
    {
      "id": 1,
      "name": "项目看板",
      "description": "项目描述",
      "color": "#3498db",
      "status": 1,
      "project_id": 1,
      "owner_id": 1,
      "position": 0,
      "created_at": "2025-11-20T10:00:00Z",
      "updated_at": "2025-11-20T10:00:00Z"
    }
  ]
}
```

---

### 5. 创建看板

**POST** `/api/boards`

**需要认证**: 是

**请求体**:
```json
{
  "name": "string (必填)",
  "description": "string (可选)",
  "color": "string (可选, 十六进制颜色, 默认#3498db)",
  "project_id": "number (必填)"
}
```

**响应** (201 Created):
```json
{
  "id": 1,
  "name": "项目看板",
  "description": "项目描述",
  "color": "#3498db",
  "status": 1,
  "project_id": 1,
  "owner_id": 1,
  "position": 0,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z"
}
```

---

### 6. 获取单个看板（包含嵌套的列和任务）

**GET** `/api/boards/:boardId`

**需要认证**: 是

**路径参数**:
- `boardId`: 看板ID (number)

**响应** (200 OK):
```json
{
  "id": 1,
  "name": "项目看板",
  "description": "项目描述",
  "color": "#3498db",
  "status": 1,
  "project_id": 1,
  "owner_id": 1,
  "position": 0,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z",
  "owner": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  },
  "columns": [
    {
      "id": 1,
      "name": "待办",
      "description": "待处理的任务",
      "color": "#95a5a6",
      "position": 0,
      "board_id": 1,
      "status": 1,
      "created_at": "2025-11-20T10:00:00Z",
      "tasks": [
        {
          "id": 1,
          "title": "任务标题",
          "description": "任务描述",
          "priority": 2,
          "status": 1,
          "position": 0,
          "column_id": 1,
          "creator_id": 1,
          "assignee_id": 2,
          "project_id": 1,
          "created_at": "2025-11-20T10:00:00Z",
          "assignee": {
            "id": 2,
            "username": "assignee",
            "email": "assignee@example.com"
          },
          "creator": {
            "id": 1,
            "username": "testuser",
            "email": "test@example.com"
          }
        }
      ]
    }
  ]
}
```

**错误响应**:
- `400 Bad Request`: 无效的看板ID
- `404 Not Found`: 看板不存在

---

### 7. 更新看板

**PUT** `/api/boards/:boardId`

**需要认证**: 是

**路径参数**:
- `boardId`: 看板ID (number)

**请求体** (所有字段可选):
```json
{
  "name": "string (可选)",
  "description": "string (可选)",
  "color": "string (可选)",
  "status": "number (可选, 1=活跃, 2=归档)",
  "position": "number (可选)"
}
```

**响应** (200 OK):
```json
{
  "message": "更新成功"
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 看板不存在

---

### 8. 删除看板

**DELETE** `/api/boards/:boardId`

**需要认证**: 是

**路径参数**:
- `boardId`: 看板ID (number)

**响应** (200 OK):
```json
{
  "message": "删除成功"
}
```

**错误响应**:
- `400 Bad Request`: 无效的看板ID
- `404 Not Found`: 看板不存在

---

## 列相关接口

### 9. 获取看板的所有列

**GET** `/api/boards/:boardId/columns`

**需要认证**: 是

**路径参数**:
- `boardId`: 看板ID (number)

**响应** (200 OK):
```json
{
  "columns": [
    {
      "id": 1,
      "name": "待办",
      "description": "待处理的任务",
      "color": "#95a5a6",
      "position": 0,
      "board_id": 1,
      "status": 1,
      "created_at": "2025-11-20T10:00:00Z",
      "updated_at": "2025-11-20T10:00:00Z"
    }
  ]
}
```

---

### 10. 创建列

**POST** `/api/boards/:boardId/columns`

**需要认证**: 是

**路径参数**:
- `boardId`: 看板ID (number)

**请求体**:
```json
{
  "name": "string (必填)",
  "description": "string (可选)",
  "color": "string (可选, 十六进制颜色, 默认#95a5a6)"
}
```

**响应** (201 Created):
```json
{
  "id": 1,
  "name": "待办",
  "description": "待处理的任务",
  "color": "#95a5a6",
  "position": 0,
  "board_id": 1,
  "status": 1,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z"
}
```

**说明**: 新创建的列会自动设置 `position` 为当前看板中列的最大位置+1

---

### 11. 获取单个列

**GET** `/api/columns/:columnId`

**需要认证**: 是

**路径参数**:
- `columnId`: 列ID (number)

**响应** (200 OK):
```json
{
  "id": 1,
  "name": "待办",
  "description": "待处理的任务",
  "color": "#95a5a6",
  "position": 0,
  "board_id": 1,
  "status": 1,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z",
  "tasks": [
    {
      "id": 1,
      "title": "任务标题",
      "description": "任务描述",
      "priority": 2,
      "status": 1,
      "position": 0,
      "column_id": 1,
      "created_at": "2025-11-20T10:00:00Z"
    }
  ]
}
```

**错误响应**:
- `400 Bad Request`: 无效的列ID
- `404 Not Found`: 列不存在

---

### 12. 更新列

**PUT** `/api/columns/:columnId`

**需要认证**: 是

**路径参数**:
- `columnId`: 列ID (number)

**请求体** (所有字段可选):
```json
{
  "name": "string (可选)",
  "description": "string (可选)",
  "color": "string (可选)",
  "status": "number (可选, 1=正常, 2=禁用)",
  "position": "number (可选)"
}
```

**响应** (200 OK):
```json
{
  "message": "更新成功"
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 列不存在

---

### 13. 删除列

**DELETE** `/api/columns/:columnId`

**需要认证**: 是

**路径参数**:
- `columnId`: 列ID (number)

**响应** (200 OK):
```json
{
  "message": "删除成功"
}
```

**错误响应**:
- `400 Bad Request`: 无效的列ID
- `404 Not Found`: 列不存在

---

## 任务相关接口

### 14. 获取列的所有任务

**GET** `/api/columns/:columnId/tasks`

**需要认证**: 是

**路径参数**:
- `columnId`: 列ID (number)

**响应** (200 OK):
```json
{
  "tasks": [
    {
      "id": 1,
      "title": "任务标题",
      "description": "任务描述",
      "priority": 2,
      "status": 1,
      "position": 0,
      "column_id": 1,
      "creator_id": 1,
      "assignee_id": 2,
      "project_id": 1,
      "due_date": "2025-12-01T00:00:00Z",
      "start_date": null,
      "end_date": null,
      "estimated_hours": 8.0,
      "actual_hours": null,
      "created_at": "2025-11-20T10:00:00Z",
      "updated_at": "2025-11-20T10:00:00Z"
    }
  ]
}
```

---

### 15. 创建任务

**POST** `/api/columns/:columnId/tasks`

**需要认证**: 是

**路径参数**:
- `columnId`: 列ID (number)

**请求体**:
```json
{
  "title": "string (必填)",
  "description": "string (可选)",
  "priority": "number (可选, 1=低, 2=中, 3=高, 4=紧急, 默认2)",
  "due_date": "string (可选, ISO 8601格式)",
  "start_date": "string (可选, ISO 8601格式)",
  "estimated_hours": "number (可选)",
  "assignee_id": "number (可选)",
  "project_id": "number (必填)"
}
```

**响应** (201 Created):
```json
{
  "id": 1,
  "title": "任务标题",
  "description": "任务描述",
  "priority": 2,
  "status": 1,
  "position": 0,
  "column_id": 1,
  "creator_id": 1,
  "assignee_id": 2,
  "project_id": 1,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z"
}
```

**说明**: 
- 新创建的任务会自动设置 `position` 为当前列中任务的最大位置+1
- `status` 默认为 1 (待办)
- `creator_id` 自动设置为当前登录用户ID

---

### 16. 获取单个任务

**GET** `/api/tasks/:taskId`

**需要认证**: 是

**路径参数**:
- `taskId`: 任务ID (number)

**响应** (200 OK):
```json
{
  "id": 1,
  "title": "任务标题",
  "description": "任务描述",
  "priority": 2,
  "status": 1,
  "position": 0,
  "column_id": 1,
  "creator_id": 1,
  "assignee_id": 2,
  "project_id": 1,
  "due_date": "2025-12-01T00:00:00Z",
  "start_date": null,
  "end_date": null,
  "estimated_hours": 8.0,
  "actual_hours": null,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z",
  "assignee": {
    "id": 2,
    "username": "assignee",
    "email": "assignee@example.com"
  },
  "creator": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  },
  "column": {
    "id": 1,
    "name": "待办",
    "board_id": 1
  }
}
```

**错误响应**:
- `400 Bad Request`: 无效的任务ID
- `404 Not Found`: 任务不存在

---

### 17. 更新任务

**PUT** `/api/tasks/:taskId`

**需要认证**: 是

**路径参数**:
- `taskId`: 任务ID (number)

**请求体** (所有字段可选):
```json
{
  "title": "string (可选)",
  "description": "string (可选)",
  "priority": "number (可选, 1=低, 2=中, 3=高, 4=紧急)",
  "status": "number (可选, 1=待办, 2=进行中, 3=已完成, 4=已取消)",
  "due_date": "string (可选, ISO 8601格式)",
  "start_date": "string (可选, ISO 8601格式)",
  "end_date": "string (可选, ISO 8601格式)",
  "estimated_hours": "number (可选)",
  "actual_hours": "number (可选)",
  "assignee_id": "number (可选, null表示取消分配)"
}
```

**响应** (200 OK):
```json
{
  "message": "更新成功"
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 任务不存在

---

### 18. 删除任务

**DELETE** `/api/tasks/:taskId`

**需要认证**: 是

**路径参数**:
- `taskId`: 任务ID (number)

**响应** (200 OK):
```json
{
  "message": "删除成功"
}
```

**错误响应**:
- `400 Bad Request`: 无效的任务ID
- `404 Not Found`: 任务不存在

---

### 19. 移动任务（拖拽排序）

**PATCH** `/api/tasks/:taskId/move`

**需要认证**: 是

**路径参数**:
- `taskId`: 任务ID (number)

**请求体**:
```json
{
  "newColumnId": "number (必填, 目标列ID)",
  "newOrder": "number (必填, 在目标列中的新位置索引, 从0开始)"
}
```

**响应** (200 OK):
```json
{
  "message": "移动成功"
}
```

**功能说明**:
- 支持跨列移动：将任务从一个列移动到另一个列
- 支持同列内移动：在同一列内调整任务顺序
- 自动更新相关任务的 `position` 值，保证排序正确
- 使用数据库事务保证操作的一致性

**示例**:
```json
// 将任务移动到列2的第0个位置
{
  "newColumnId": 2,
  "newOrder": 0
}

// 在同一列内将任务移动到第3个位置
{
  "newColumnId": 1,
  "newOrder": 3
}
```

**错误响应**:
- `400 Bad Request`: 请求参数错误（缺少 newColumnId 或 newOrder）
- `404 Not Found`: 任务不存在

---

## 数据模型说明

### Board (看板)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | number | 看板ID |
| name | string | 看板名称 |
| description | string | 看板描述 |
| color | string | 看板颜色（十六进制） |
| status | number | 看板状态（1=活跃, 2=归档） |
| project_id | number | 所属项目ID |
| owner_id | number | 所有者用户ID |
| position | number | 在看板列表中的排序位置 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

### Column (列)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | number | 列ID |
| name | string | 列名称 |
| description | string | 列描述 |
| color | string | 列颜色（十六进制） |
| position | number | 在看板中的排序位置 |
| board_id | number | 所属看板ID |
| status | number | 列状态（1=正常, 2=禁用） |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

### Task (任务)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | number | 任务ID |
| title | string | 任务标题 |
| description | string | 任务描述 |
| priority | number | 优先级（1=低, 2=中, 3=高, 4=紧急） |
| status | number | 任务状态（1=待办, 2=进行中, 3=已完成, 4=已取消） |
| position | number | 在列中的排序位置 |
| column_id | number | 所属列ID |
| creator_id | number | 创建者用户ID |
| assignee_id | number | 分配给的用户ID（可选） |
| project_id | number | 所属项目ID |
| due_date | string | 截止日期（ISO 8601格式） |
| start_date | string | 开始日期（ISO 8601格式） |
| end_date | string | 结束日期（ISO 8601格式） |
| estimated_hours | number | 预估工时（小时） |
| actual_hours | number | 实际工时（小时） |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

---

## 错误码说明

| HTTP状态码 | 说明 |
|-----------|------|
| 200 OK | 请求成功 |
| 201 Created | 资源创建成功 |
| 400 Bad Request | 请求参数错误 |
| 401 Unauthorized | 未认证或认证失败 |
| 403 Forbidden | 无权限访问 |
| 404 Not Found | 资源不存在 |
| 409 Conflict | 资源冲突（如用户名已存在） |
| 500 Internal Server Error | 服务器内部错误 |

---

## 注意事项

1. **认证**: 除登录和注册接口外，所有接口都需要在请求头中携带有效的 JWT token
2. **时间格式**: 所有日期时间字段使用 ISO 8601 格式（例如: `2025-11-20T10:00:00Z`）
3. **排序**: 列和任务默认按 `position` 字段升序排列
4. **软删除**: 删除操作使用软删除，数据不会真正从数据库中删除
5. **嵌套结构**: `GET /api/boards/:boardId` 返回的看板对象包含完整的嵌套结构（列和任务）
6. **拖拽排序**: 使用 `PATCH /api/tasks/:taskId/move` 接口进行任务拖拽排序，系统会自动处理位置更新

---

## 更新日志

- **2025-11-20**: 初始版本，包含看板、列、任务的完整CRUD接口和拖拽排序功能

