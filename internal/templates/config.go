package templates

// ApplicationYML is the main application.yml template
const ApplicationYML = `spring:
  application:
    name: {{ARTIFACT_ID}}
  profiles:
    active: dev

server:
  port: 8080

management:
  endpoints:
    web:
      exposure:
        include: health,info,prometheus,metrics
  metrics:
    export:
      prometheus:
        enabled: true

mybatis-plus:
  configuration:
    map-underscore-to-camel-case: true
    log-impl: org.apache.ibatis.logging.stdout.StdOutImpl
  mapper-locations: classpath*:/mapper/**/*Mapper.xml
`

// ApplicationDevYML is the dev profile application.yml template
const ApplicationDevYML = `spring:
  datasource:
    url: jdbc:mysql://localhost:3306/{{ARTIFACT_ID}}?useSSL=false&serverTimezone=Asia/Shanghai&characterEncoding=utf8
    username: root
    password: root
    driver-class-name: com.mysql.cj.jdbc.Driver
    hikari:
      maximum-pool-size: 20
      minimum-idle: 5
      connection-timeout: 30000

logging:
  level:
    root: INFO
    {{PACKAGE_NAME}}: DEBUG
  pattern:
    console: "%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
`

// README is the README.md template
const README = `# {{PROJECT_NAME}}

{{PROJECT_DESCRIPTION}}

## 快速开始

### 构建项目

` + "```" + `bash
mvn clean install
` + "```" + `

### 运行应用

` + "```" + `bash
java --enable-preview -jar starter/target/starter-{{VERSION}}.jar
` + "```" + `

### 测试

` + "```" + `bash
curl http://localhost:8080/api/health
` + "```" + `

## 项目结构

` + "```" + `
{{ARTIFACT_ID}}/
├── common/              # 公共模块
├── domain/              # 领域层
├── infrastructure/      # 基础设施层
├── adapter/
│   ├── adapter-rest/   # REST接口
│   └── adapter-schedule/ # 定时任务
├── application/
│   └── application-user/ # 用户业务
└── starter/            # 启动模块
` + "```" + `

## 添加新的业务模块

使用 ` + "`phjvgen add <module-name>`" + ` 命令添加新的业务模块：

` + "```" + `bash
phjvgen add payment
` + "```" + `

## 技术栈

- Java 25 LTS
- Spring Boot 4.0.0-RC1
- MyBatis Plus 3.5.8+
- MySQL 8.0+
`

// GitIgnore is the .gitignore template
const GitIgnore = `# Maven
target/
pom.xml.tag
pom.xml.releaseBackup
pom.xml.versionsBackup
pom.xml.next
release.properties
dependency-reduced-pom.xml

# IntelliJ IDEA
.idea/
*.iml
*.iws
*.ipr
out/

# Eclipse
.classpath
.project
.settings/

# VS Code
.vscode/

# Java
*.class
*.jar
*.war
*.ear
hs_err_pid*

# Logs
logs/
*.log

# OS
.DS_Store
Thumbs.db

# Application
application-local.yml
`

// UserTableSQL is the SQL for creating user table
const UserTableSQL = `CREATE TABLE IF NOT EXISTS ` + "`t_user`" + ` (
    ` + "`id`" + ` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    ` + "`username`" + ` VARCHAR(50) NOT NULL COMMENT '用户名',
    ` + "`email`" + ` VARCHAR(100) COMMENT '邮箱',
    ` + "`phone`" + ` VARCHAR(20) COMMENT '手机号',
    ` + "`status`" + ` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    ` + "`create_time`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    ` + "`update_time`" + ` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    ` + "`deleted`" + ` TINYINT NOT NULL DEFAULT 0 COMMENT '删除标记：0-未删除，1-已删除',
    PRIMARY KEY (` + "`id`" + `),
    UNIQUE KEY ` + "`uk_username`" + ` (` + "`username`" + `),
    KEY ` + "`idx_email`" + ` (` + "`email`" + `),
    KEY ` + "`idx_phone`" + ` (` + "`phone`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
`
