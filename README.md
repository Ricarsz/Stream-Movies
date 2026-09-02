# goMovies

全栈电影网站学习 Demo，参考 FreeCodeCamp Go 全栈教程，使用 Go + Gin + MongoDB 后端 + React 前端。

## 技术栈

后端：
- Go 1.27 + Gin Web Framework
- MongoDB 8.x（go-driver v2）
- JWT 认证（golang-jwt/v5）
- bcrypt 密码加密
- CORS 跨域支持

前端：
- React 19 + Vite
- React Router DOM
- Axios（自动携带 Token）
- 暗色主题 CSS

## 项目结构

```
goMovies/
├── server/MovieServer/
│   ├── cmd/main.go                 # 入口，路由注册
│   ├── controllers/
│   │   ├── movie_controller.go     # 电影 CRUD
│   │   ├── user_controller.go      # 注册 + 登录
│   │   └── auth_middleware.go      # JWT 中间件
│   ├── database/
│   │   ├── connection.go           # MongoDB 连接
│   │   └── interfaces.go           # Collection 接口抽象
│   ├── models/
│   │   ├── movieModel.go           # Movie / Ranking / Genre
│   │   └── modelUser.go            # User
│   ├── repository/
│   │   ├── movie_repository.go     # 电影数据访问
│   │   └── user_repository.go      # 用户数据访问
│   ├── utils/
│   │   ├── tool.go                 # 密码哈希 + 验证
│   │   └── token.go                # JWT 签发 + 校验
│   ├── mocks/                      # 测试 Mock
│   ├── tests/                      # 集成测试
│   ├── .env                        # 环境变量
│   ├── Makefile                    # build / run / test
│   └── go.mod
└── client/
    ├── src/
    │   ├── api.js                  # Axios 封装
    │   ├── context/AuthContext.jsx # 登录状态管理
    │   ├── components/Navbar.jsx   # 导航栏
    │   ├── pages/
    │   │   ├── Home.jsx            # 电影列表
    │   │   ├── MovieDetail.jsx     # 电影详情
    │   │   ├── Login.jsx           # 登录
    │   │   ├── Register.jsx        # 注册
    │   │   ├── AddMovie.jsx        # 添加电影（Admin）
    │   │   └── EditMovie.jsx       # 编辑电影（Admin）
    │   ├── App.jsx                 # 路由
    │   └── App.css                 # 样式
    └── vite.config.js              # Vite + API Proxy
```

## 启动

前置要求：Go 1.27+、MongoDB、fnm/Node.js

```bash
# 1. 启动 MongoDB
mongod --dbpath /tmp/mongodata --fork --logpath /tmp/mongod.log

# 2. 启动后端（端口 8080）
cd server/MovieServer
make run

# 3. 启动前端（端口 5173）
cd client
eval "$(fnm env)"
npm run dev
```

访问 http://localhost:5173

## API 接口

| 方法     | 路径               | 说明         | 认证   |
|----------|--------------------|--------------|--------|
| POST     | /users/register    | 用户注册     | 无     |
| POST     | /users/login       | 用户登录     | 无     |
| GET      | /movies            | 电影列表     | 无     |
| GET      | /movies/:imdb_id   | 电影详情     | 无     |
| POST     | /movies            | 添加电影     | Bearer |
| PUT      | /movies/:imdb_id   | 更新电影     | Bearer |
| DELETE   | /movies/:imdb_id   | 删除电影     | Bearer |

## 测试

```bash
cd server/MovieServer
make test
```

13 个集成测试覆盖电影 CRUD、用户注册、重复邮箱校验、密码哈希、参数验证。

## 环境变量

```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=stream-movies
JWT_SECRET=your-secret-key
```

## 学到什么

- Go + Gin 构建 REST API
- MongoDB 驱动操作与接口抽象
- Repository 模式分层
- JWT 认证流程（签发 / 刷新 / 中间件）
- bcrypt 密码安全存储
- React 路由与 Context 状态管理
- Vite 开发代理与构建
- 集成测试与 Mock
