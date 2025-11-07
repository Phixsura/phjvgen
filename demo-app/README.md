# Demo Application

A demo application for testing

## 快速开始

### 构建项目

```bash
mvn clean install
```

### 运行应用

```bash
java --enable-preview -jar starter/target/starter-1.0.0.jar
```

### 测试

```bash
curl http://localhost:8080/api/health
```

## 项目结构

```
demo-app/
├── common/              # 公共模块
├── domain/              # 领域层
├── infrastructure/      # 基础设施层
├── adapter/
│   ├── adapter-rest/   # REST接口
│   └── adapter-schedule/ # 定时任务
├── application/
│   └── application-user/ # 用户业务
└── starter/            # 启动模块
```

## 添加新的业务模块

使用 `phjvgen add <module-name>` 命令添加新的业务模块：

```bash
phjvgen add payment
```

## 技术栈

- Java 25 LTS
- Spring Boot 4.0.0-RC1
- MyBatis Plus 3.5.8+
- MySQL 8.0+
