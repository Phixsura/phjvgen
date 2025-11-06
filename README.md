# phjvgen - Java 25 LTS Project Generator

一个基于 Go 语言开发的 Java 25 LTS 分层架构项目生成器。完全替代 shell 脚本实现，支持 macOS 和 Linux。

## 特性

- **纯 Go 实现**：不依赖任何 shell 命令，跨平台兼容
- **DDD 分层架构**：生成完整的领域驱动设计分层结构
- **Maven 多模块**：自动生成标准的 Maven 多模块项目
- **Spring Boot 4.0**：支持最新的 Spring Boot 4.0.0-RC1
- **Java 25 LTS**：使用 Java 25 的最新特性
- **完整示例**：可选生成完整的 CRUD 示例代码
- **模块化扩展**：轻松添加新的业务模块

## 安装

### 从源码构建

```bash
cd phjvgen
go build -o phjvgen .
```

### 安装到系统

构建完成后，可以使用以下命令安装到系统：

```bash
# 方法1：使用 install 命令
./phjvgen install

# 方法2：手动复制
# macOS/Linux
cp phjvgen ~/.local/bin/
chmod +x ~/.local/bin/phjvgen

# 确保 ~/.local/bin 在你的 PATH 中
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## 使用

### 生成新项目

交互式生成新项目：

```bash
phjvgen generate
# 或使用别名
phjvgen gen
phjvgen g
```

### 生成示例项目

快速生成预配置的示例项目：

```bash
phjvgen example
```

这会生成一个位于 `./demo-app` 的示例项目，配置如下：
- Group ID: com.example.demo
- Artifact ID: demo-app
- Version: 1.0.0

### 生成 CRUD 示例代码

在现有项目中生成完整的 User 模块 CRUD 示例：

```bash
cd your-project
phjvgen demo
```

这会生成：
- Domain 层：User 实体、UserRepository 接口
- Infrastructure 层：UserDO、UserMapper、UserRepositoryImpl
- Application 层：UserDTO、UserService、UserAssembler、Commands
- Adapter 层：UserController、Request/Response、全局异常处理
- 数据库脚本：用户表创建 SQL

### 添加新的业务模块

在现有项目中添加新的 Application 模块：

```bash
cd your-project
phjvgen add payment        # 创建 application-payment 模块
phjvgen add order          # 创建 application-order 模块
phjvgen add user-profile   # 创建 application-user-profile 模块
```

### 查看版本

```bash
phjvgen version
```

### 查看帮助

```bash
phjvgen --help
phjvgen generate --help
phjvgen add --help
```

## 项目结构

生成的项目结构如下：

```
your-project/
├── pom.xml                      # 父 POM
├── common/                      # 公共模块
│   └── src/main/java/.../common/
│       ├── exception/           # 异常类
│       ├── response/            # 响应封装
│       ├── constant/            # 常量
│       └── utils/              # 工具类
├── domain/                      # 领域层
│   └── src/main/java/.../domain/
│       ├── model/              # 领域实体
│       ├── repository/         # 仓储接口
│       ├── service/            # 领域服务
│       └── event/              # 领域事件
├── infrastructure/              # 基础设施层
│   └── src/main/java/.../infrastructure/
│       ├── persistence/        # 持久化
│       │   ├── mapper/        # MyBatis Mapper
│       │   └── impl/          # Repository 实现
│       ├── cache/             # 缓存
│       ├── mq/                # 消息队列
│       └── gateway/           # 外部网关
├── adapter/                    # 适配器层
│   ├── adapter-rest/          # REST 适配器
│   │   └── src/main/java/.../adapter/rest/
│   │       ├── controller/    # 控制器
│   │       ├── request/       # 请求 VO
│   │       ├── response/      # 响应 VO
│   │       └── advice/        # 全局异常处理
│   └── adapter-schedule/      # 定时任务适配器
│       └── src/main/java/.../adapter/schedule/
│           └── job/           # 定时任务
├── application/                # 应用层
│   └── application-user/      # 用户业务模块
│       └── src/main/java/.../application/user/
│           ├── service/       # 应用服务
│           ├── dto/           # DTO
│           ├── assembler/     # 对象转换器
│           └── executor/      # 执行器
└── starter/                    # 启动模块
    └── src/main/java/.../     # Application 主类
```

## 完整工作流程

### 1. 生成新项目

```bash
# 生成项目
phjvgen generate

# 输入项目信息
Group ID: com.mycompany
Artifact ID: my-app
Version: 1.0.0
Project Name: My Application
Description: My awesome application
Output Directory: ./my-app
```

### 2. 添加 CRUD 示例（可选）

```bash
cd my-app
phjvgen demo
```

### 3. 配置数据库

编辑 `starter/src/main/resources/application-dev.yml`：

```yaml
spring:
  datasource:
    url: jdbc:mysql://localhost:3306/my_app
    username: root
    password: your_password
```

### 4. 创建数据库并执行 SQL（如果生成了 demo）

```bash
mysql -u root -p
> CREATE DATABASE my_app;
> USE my_app;
> source infrastructure/src/main/resources/db/migration/V1__create_user_table.sql;
```

### 5. 构建和运行

```bash
# 构建项目
mvn clean install

# 运行应用
java --enable-preview -jar starter/target/starter-1.0.0.jar

# 测试
curl http://localhost:8080/api/health
curl http://localhost:8080/api/users
```

### 6. 添加新模块

```bash
# 添加支付模块
phjvgen add payment

# 重新构建
mvn clean install
```

## 技术栈

### 生成器本身
- Go 1.24+
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Color](https://github.com/fatih/color) - 终端颜色输出

### 生成的项目
- Java 25 LTS
- Spring Boot 4.0.0-RC1
- MyBatis Plus 3.5.8
- MySQL 8.0
- Lombok 1.18.42
- MapStruct 1.6.0
- Maven 3.x

## 开发

### 项目结构

```
phjvgen/
├── cmd/                        # Cobra 命令定义
│   ├── root.go                # 根命令
│   ├── generate.go            # generate 命令
│   ├── demo.go                # demo 命令
│   ├── add.go                 # add 命令
│   ├── example.go             # example 命令
│   ├── install.go             # install 命令
│   └── version.go             # version 命令
├── internal/
│   ├── generator/             # 生成器逻辑
│   │   ├── config.go         # 配置管理
│   │   ├── project.go        # 项目生成
│   │   ├── demo.go           # CRUD 示例生成
│   │   └── module.go         # 模块添加
│   ├── templates/             # 模板文件
│   │   ├── pom.go            # POM 模板
│   │   ├── java.go           # Java 代码模板
│   │   └── config.go         # 配置文件模板
│   └── utils/                 # 工具函数
│       ├── color.go          # 颜色输出
│       ├── file.go           # 文件操作
│       └── input.go          # 用户输入
├── main.go                    # 主入口
├── go.mod                     # Go 模块定义
└── README.md                  # 本文件
```

### 构建

```bash
# 开发构建
go build -o phjvgen .

# 生产构建（优化）
go build -ldflags="-s -w" -o phjvgen .

# 交叉编译
# Linux
GOOS=linux GOARCH=amd64 go build -o phjvgen-linux .

# macOS
GOOS=darwin GOARCH=amd64 go build -o phjvgen-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o phjvgen-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o phjvgen.exe .
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/generator/...

# 带覆盖率
go test -cover ./...

# 详细输出
go test -v ./...
```

## 与原 Shell 脚本的区别

### 优势

1. **跨平台**：不依赖 bash、sed、awk 等 Unix 工具
2. **性能更好**：Go 编译后的二进制执行效率更高
3. **错误处理**：更完善的错误处理和提示
4. **代码维护**：Go 代码比复杂的 shell 脚本更易维护
5. **类型安全**：编译时类型检查，减少运行时错误
6. **单一二进制**：所有功能打包在一个可执行文件中

### 功能对等

| Shell 脚本 | Go 实现 | 说明 |
|-----------|---------|------|
| generate-project.sh | `phjvgen generate` | 生成项目 |
| generate-demo-code.sh | `phjvgen demo` | 生成 CRUD 示例 |
| add-application-module.sh | `phjvgen add` | 添加模块 |
| example.sh | `phjvgen example` | 生成示例项目 |
| install.sh | `phjvgen install` | 安装到系统 |

## 常见问题

### Q: 如何修改生成的项目结构？

A: 编辑 `internal/templates/` 目录下的模板文件，然后重新构建 phjvgen。

### Q: 支持 Windows 吗？

A: 是的，phjvgen 是纯 Go 实现，完全支持 Windows。

### Q: 生成的项目可以用其他数据库吗？

A: 可以，生成后修改 pom.xml 中的数据库依赖和 application.yml 中的配置即可。

### Q: 如何添加自定义模板？

A: 在 `internal/templates/` 中添加新的模板常量，然后在相应的生成器中使用。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 作者

phixia team
