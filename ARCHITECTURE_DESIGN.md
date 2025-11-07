# DDD 分层架构设计理念

本文档详细阐述了 phjvgen 生成项目所采用的领域驱动设计（DDD）分层架构的核心设计理念和最佳实践。

---

## 目录

1. [架构概览](#架构概览)
2. [分层职责详解](#分层职责详解)
3. [Service vs Executor 的区别](#service-vs-executor-的区别)
4. [领域事件机制](#领域事件机制)
5. [对象转换策略](#对象转换策略)
6. [依赖注入最佳实践](#依赖注入最佳实践)
7. [完整请求流程示例](#完整请求流程示例)
8. [设计原则总结](#设计原则总结)

---

## 架构概览

### 分层结构

```
┌─────────────────────────────────────────────────────────┐
│                    Adapter Layer                        │
│  (适配器层 - 对外接口、HTTP/MQ/Job等)                    │
│  - REST Controller                                       │
│  - Request/Response VO                                   │
│  - Filter & Interceptor                                  │
│  - Global Exception Handler                              │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                  Application Layer                      │
│  (应用层 - 用例编排、业务流程)                           │
│  - Application Service (薄层 CRUD)                       │
│  - Use Case Executor (复杂用例编排)                      │
│  - DTO & Command                                         │
│  - Assembler (对象转换)                                  │
│  - Event Listener (事件监听器)                           │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                         │
│  (领域层 - 核心业务逻辑)                                 │
│  - Domain Entity (领域实体)                              │
│  - Domain Service (领域服务)                             │
│  - Repository Interface (仓储接口)                       │
│  - Domain Event (领域事件)                               │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                Infrastructure Layer                     │
│  (基础设施层 - 技术实现)                                 │
│  - Repository Impl (仓储实现)                            │
│  - Data Object (DO)                                      │
│  - MyBatis Mapper                                        │
│  - External Service Gateway                              │
└─────────────────────────────────────────────────────────┘
                          ↓
                    Database / MQ / Cache
```

### 依赖方向原则

**核心原则**：依赖只能从外向内，内层不依赖外层

```
Adapter → Application → Domain ← Infrastructure
   ↓          ↓          ↑            ↑
  依赖      依赖      被依赖      实现接口
```

- **Adapter** 依赖 Application
- **Application** 依赖 Domain
- **Infrastructure** 实现 Domain 定义的接口（依赖倒置）
- **Domain** 不依赖任何外层（最纯粹的业务逻辑）

---

## 分层职责详解

### 1. Domain Layer（领域层）

**定位**：系统的核心，包含业务规则和领域逻辑

**包含组件**：

#### 1.1 Domain Entity（领域实体）
```java
@Data
public class User {
    private Long id;
    private String username;
    private String email;
    private Integer status;

    // 领域方法 - 封装业务规则
    public void enable() {
        this.status = 1;
    }

    public void disable() {
        this.status = 0;
    }

    public boolean isActive() {
        return this.status == 1;
    }
}
```

**职责**：
- 封装核心业务数据
- 提供领域方法（而不是简单的 getter/setter）
- 保证业务不变性（invariants）

#### 1.2 Domain Service（领域服务）
```java
@Service
@RequiredArgsConstructor
public class UserDomainService {
    private final UserRepository userRepository;
    private final ApplicationEventPublisher eventPublisher;

    public User registerUser(String username, String email, String phone) {
        // 1. 验证领域规则
        if (userRepository.existsByUsername(username)) {
            throw new IllegalArgumentException("用户名已存在");
        }

        // 2. 创建领域实体
        User user = new User();
        user.setUsername(username);
        user.enable();

        // 3. 持久化
        user = userRepository.save(user);

        // 4. 发布领域事件
        eventPublisher.publishEvent(new UserCreatedEvent(...));

        return user;
    }
}
```

**职责**：
- 处理跨聚合根的业务逻辑
- 执行复杂的领域规则验证
- 发布领域事件
- **不包含**：技术细节（如事务、DTO转换）

**何时使用**：
- 业务逻辑涉及多个实体
- 复杂的业务规则验证
- 需要发布领域事件的场景

#### 1.3 Repository Interface（仓储接口）
```java
public interface UserRepository {
    User save(User user);
    Optional<User> findById(Long id);
    List<User> findAll();
    boolean existsByUsername(String username);
    void deleteById(Long id);
}
```

**职责**：
- 定义数据访问契约
- 使用领域语言（而非 SQL 语言）
- 由 Infrastructure 层实现

#### 1.4 Domain Event（领域事件）
```java
@Getter
public class UserCreatedEvent {
    private final Long userId;
    private final String username;
    private final String email;
    private final LocalDateTime occurredOn;

    public UserCreatedEvent(Long userId, String username, String email) {
        this.userId = userId;
        this.username = username;
        this.email = email;
        this.occurredOn = LocalDateTime.now();
    }
}
```

**职责**：
- 表示领域中发生的重要事件
- 用于解耦不同聚合根之间的依赖
- 支持异步处理和最终一致性

---

### 2. Application Layer（应用层）

**定位**：编排用例，协调 Domain 和 Infrastructure

**包含组件**：

#### 2.1 Application Service（应用服务）
```java
@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final UserAssembler userAssembler;
    private final RegisterUserExecutor registerUserExecutor;

    // 复杂用例 - 委托给 Executor
    @Transactional
    public UserDTO createUser(CreateUserCommand command) {
        return registerUserExecutor.execute(command);
    }

    // 简单 CRUD - 直接在 Service 中完成
    @Transactional
    public UserDTO updateUser(UpdateUserCommand command) {
        User user = userRepository.findById(command.getId())
            .orElseThrow(() -> new BusinessException("用户不存在"));
        userAssembler.updateEntity(user, command);
        user = userRepository.update(user);
        return userAssembler.toDTO(user);
    }

    // 查询操作 - 简单直接
    public UserDTO getUserById(Long id) {
        User user = userRepository.findById(id)
            .orElseThrow(() -> new BusinessException("用户不存在"));
        return userAssembler.toDTO(user);
    }
}
```

**职责**：
- 提供面向数据的 CRUD 接口
- 薄薄一层，主要做：
  - 参数校验（基础校验）
  - 事务控制（@Transactional）
  - DTO ↔ Entity 转换
  - 调用 Repository
  - 对于复杂用例，委托给 Executor

**Service 应该保持"薄"**：
- ❌ 不要包含复杂的业务流程编排
- ❌ 不要调用多个外部服务
- ❌ 不要包含复杂的业务逻辑
- ✅ 简单 CRUD 直接完成
- ✅ 复杂用例委托给 Executor

#### 2.2 Use Case Executor（用例执行器）
```java
@Component
@RequiredArgsConstructor
public class RegisterUserExecutor {
    private final UserDomainService userDomainService;
    private final UserAssembler userAssembler;
    // 真实项目中还会注入：
    // private final CouponService couponService;
    // private final RiskControlService riskControlService;
    // private final SmsService smsService;

    @Transactional
    public UserDTO execute(CreateUserCommand command) {
        // 步骤1: 调用领域服务执行核心逻辑
        User user = userDomainService.registerUser(
            command.getUsername(),
            command.getEmail(),
            command.getPhone()
        );

        // 在复杂项目中，这里还会：
        // - 调用风控服务验证（同步）
        // - 调用实名认证服务（同步）
        // - 初始化积分账户（异步事件）
        // - 发放新人优惠券（异步事件）
        // - 发送短信验证码（同步）

        return userAssembler.toDTO(user);
    }
}
```

**职责**：
- 编排复杂的业务流程
- 调用多个 Domain Service
- 调用多个 Repository
- 调用外部服务（短信、邮件、第三方 API）
- 处理复杂的条件分支和业务编排

**何时使用 Executor**：
- 业务流程涉及多个步骤
- 需要调用多个服务协同完成
- 有复杂的条件分支和业务规则
- 需要与外部系统交互

#### 2.3 DTO & Command（数据传输对象）
```java
// 查询结果 DTO
@Data
public class UserDTO {
    private Long id;
    private String username;
    private String email;
    private String phone;
    private Integer status;
    private LocalDateTime createTime;
}

// 命令对象
@Data
public class CreateUserCommand {
    @NotBlank
    private String username;
    @Email
    private String email;
    private String phone;
}
```

**职责**：
- 在层与层之间传递数据
- 携带参数校验注解
- 与领域实体解耦

#### 2.4 Assembler（对象转换器）
```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {
    // Entity → DTO
    UserDTO toDTO(User user);

    // Command → Entity
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "status", constant = "1")
    User toEntity(CreateUserCommand command);

    // 更新 Entity（只更新非 null 字段）
    @BeanMapping(nullValuePropertyMappingStrategy =
                 NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "username", ignore = true)
    void updateEntity(@MappingTarget User user, UpdateUserCommand command);
}
```

**职责**：
- 使用 MapStruct 自动生成转换代码
- 避免手写重复的转换逻辑
- 支持复杂的映射规则

#### 2.5 Event Listener（事件监听器）
```java
@Component
@Slf4j
public class UserEventListener {
    @Async
    @EventListener
    public void handleUserCreated(UserCreatedEvent event) {
        // 异步处理后续业务逻辑
        sendWelcomeEmail(event.getEmail());
        recordUserRegistration(event.getUserId());
        // 发放优惠券、初始化积分等
    }
}
```

**职责**：
- 监听领域事件
- 异步执行后续业务逻辑
- 解耦业务流程

---

### 3. Adapter Layer（适配器层）

**定位**：对外提供接口，适配不同的访问方式

**包含组件**：

#### 3.1 REST Controller
```java
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {
    private final UserService userService;
    private final UserControllerAssembler assembler;
    private final ResponseVOAssembler responseAssembler;

    @PostMapping
    public Result<UserResponseVO> createUser(
            @Validated @RequestBody CreateUserRequest request) {
        // Request → Command
        CreateUserCommand command = assembler.toCreateCommand(request);

        // 调用 Application Service
        UserDTO dto = userService.createUser(command);

        // DTO → Response VO（数据脱敏、格式转换）
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);

        return Result.success(vo);
    }
}
```

**职责**：
- 接收 HTTP 请求
- Request → Command 转换
- 调用 Application Service
- DTO → Response VO 转换（脱敏、格式化）
- 统一返回格式封装

#### 3.2 Request/Response VO
```java
// Request - 接收前端请求
@Data
public class CreateUserRequest {
    @NotBlank(message = "用户名不能为空")
    private String username;

    @Email(message = "邮箱格式不正确")
    private String email;

    private String phone;
}

// Response VO - 返回给前端（脱敏、格式化）
@Data
public class UserResponseVO {
    private Long id;
    private String username;
    private String email;
    private String phone;  // 已脱敏：138****5678
    private String statusText;  // 已转换："启用"/"禁用"

    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createTime;
}
```

**职责**：
- Request：接收并校验前端数据
- Response VO：
  - 数据脱敏（手机号、身份证等）
  - 格式转换（时间格式、枚举转文本）
  - 只返回前端需要的字段

**为什么不直接返回 DTO**：
- DTO 是内部数据传输对象，包含完整信息
- Response VO 是视图对象，面向前端：
  - 数据脱敏处理
  - 格式友好转换
  - 字段按需裁剪

#### 3.3 Filter & Interceptor
```java
// Filter - 请求日志
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class LoggingFilter implements Filter {
    @Override
    public void doFilter(ServletRequest request, ServletResponse response,
                         FilterChain chain) {
        long startTime = System.currentTimeMillis();
        log.info("Request started: {} {}", method, uri);

        chain.doFilter(request, response);

        long duration = System.currentTimeMillis() - startTime;
        log.info("Request completed in {}ms", duration);
    }
}

// Interceptor - 权限认证
@Component
public class AuthInterceptor implements HandlerInterceptor {
    @Override
    public boolean preHandle(HttpServletRequest request,
                             HttpServletResponse response,
                             Object handler) {
        String token = request.getHeader("Authorization");
        // 验证 token
        return true;
    }
}
```

**职责**：
- Filter：请求日志、CORS、编码处理
- Interceptor：认证、鉴权、参数预处理

#### 3.4 Global Exception Handler
```java
@RestControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(BusinessException.class)
    public Result<?> handleBusinessException(BusinessException e) {
        return Result.error(e.getCode(), e.getMessage());
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public Result<?> handleValidationException(
            MethodArgumentNotValidException e) {
        String message = e.getBindingResult()
            .getFieldError().getDefaultMessage();
        return Result.error(ErrorCode.PARAM_ERROR, message);
    }
}
```

**职责**：
- 统一异常处理
- 转换为统一的返回格式
- 记录异常日志

---

### 4. Infrastructure Layer（基础设施层）

**定位**：提供技术实现，实现 Domain 定义的接口

**包含组件**：

#### 4.1 Repository Implementation
```java
@Repository
@RequiredArgsConstructor
public class UserRepositoryImpl implements UserRepository {
    private final UserMapper userMapper;

    @Override
    public User save(User user) {
        UserDO userDO = toDataObject(user);
        userMapper.insert(userDO);
        user.setId(userDO.getId());
        return user;
    }

    @Override
    public Optional<User> findById(Long id) {
        UserDO userDO = userMapper.selectById(id);
        return Optional.ofNullable(userDO)
            .map(this::toEntity);
    }

    // Entity ↔ DO 转换
    private UserDO toDataObject(User user) { ... }
    private User toEntity(UserDO userDO) { ... }
}
```

**职责**：
- 实现 Domain 层定义的 Repository 接口
- 封装数据访问细节（MyBatis、JPA 等）
- Entity ↔ DO 转换

#### 4.2 Data Object (DO)
```java
@Data
@TableName("t_user")
public class UserDO {
    @TableId(type = IdType.AUTO)
    private Long id;

    private String username;
    private String email;
    private String phone;
    private Integer status;

    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createTime;

    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updateTime;
}
```

**职责**：
- 映射数据库表结构
- 携带持久化相关注解
- 与领域实体分离

**为什么需要 DO 和 Entity 分离**：
- Entity：纯粹的业务对象，不关心持久化细节
- DO：持久化对象，关注数据库映射
- 解耦业务逻辑和技术实现

---

## Service vs Executor 的区别

这是初学者最容易混淆的概念。

### 核心区别

| 维度 | Application Service | Use Case Executor |
|------|-------------------|------------------|
| **定位** | 面向数据的 CRUD 操作 | 面向业务用例的流程编排 |
| **复杂度** | 简单、薄薄一层 | 复杂、多步骤编排 |
| **职责** | 参数校验、事务控制、DTO转换 | 业务流程编排、服务协同 |
| **调用范围** | 通常调用单个 Repository | 可能调用多个 Domain Service、Repository、外部服务 |
| **业务逻辑** | 不包含复杂业务逻辑 | 包含复杂的业务流程编排 |

### 对比示例

#### ❌ 不好的设计 - 把复杂逻辑放在 Service 中

```java
@Service
public class UserService {
    @Transactional
    public UserDTO createUser(CreateUserCommand command) {
        // Service 变得臃肿，包含太多业务逻辑

        // 风控验证
        RiskCheckResult risk = riskControlService.check(command);
        if (risk.isHighRisk()) throw new BusinessException("风险用户");

        // 验证唯一性
        if (userRepository.existsByUsername(command.getUsername())) {
            throw new BusinessException("用户名已存在");
        }

        // 创建用户
        User user = new User();
        user.setUsername(command.getUsername());
        user = userRepository.save(user);

        // 初始化积分
        pointsService.initAccount(user.getId(), 100);

        // 发放优惠券
        if (isNewChannel(command.getChannel())) {
            couponService.grantNewUserCoupon(user.getId());
        }

        // 发送短信
        smsService.sendWelcome(user.getPhone());

        // 发布事件
        eventPublisher.publishEvent(new UserCreatedEvent(...));

        return userAssembler.toDTO(user);
    }
}
```

问题：
- Service 太臃肿，包含太多业务逻辑
- 难以测试和维护
- 违反单一职责原则

#### ✅ 好的设计 - Service + Executor 分工明确

```java
// Service - 保持简单
@Service
@RequiredArgsConstructor
public class UserService {
    private final RegisterUserExecutor registerUserExecutor;

    @Transactional
    public UserDTO createUser(CreateUserCommand command) {
        // 简单委托，Service 保持薄薄一层
        return registerUserExecutor.execute(command);
    }

    // 简单 CRUD 直接在 Service 中完成
    public UserDTO getUserById(Long id) {
        User user = userRepository.findById(id)
            .orElseThrow(() -> new BusinessException("用户不存在"));
        return userAssembler.toDTO(user);
    }
}

// Executor - 负责复杂的业务编排
@Component
@RequiredArgsConstructor
public class RegisterUserExecutor {
    private final UserDomainService userDomainService;
    private final RiskControlService riskControlService;
    private final PointsService pointsService;
    private final CouponService couponService;
    private final SmsService smsService;

    @Transactional
    public UserDTO execute(CreateUserCommand command) {
        // 步骤1: 风控验证（同步）
        RiskCheckResult risk = riskControlService.check(command);
        if (risk.isHighRisk()) {
            throw new BusinessException("风险用户");
        }

        // 步骤2: 调用领域服务创建用户
        User user = userDomainService.registerUser(
            command.getUsername(),
            command.getEmail(),
            command.getPhone()
        );

        // 步骤3: 初始化积分账户
        pointsService.initAccount(user.getId(), 100);

        // 步骤4: 发放新人优惠券（根据渠道）
        if (isNewChannel(command.getChannel())) {
            couponService.grantNewUserCoupon(user.getId());
        }

        // 步骤5: 发送短信通知
        smsService.sendWelcome(user.getPhone());

        // 事件会由 DomainService 发布
        // 后续的欢迎邮件、统计等由事件监听器异步处理

        return userAssembler.toDTO(user);
    }

    private boolean isNewChannel(String channel) {
        return "APP_DOWNLOAD".equals(channel);
    }
}
```

优势：
- Service 保持简洁，只做技术层面的事
- Executor 专注业务编排，职责清晰
- 易于测试、维护和扩展

### 使用指南

**何时使用 Service（不需要 Executor）**：
- 简单的 CRUD 操作
- 单表查询
- 简单的更新和删除
- 不涉及复杂业务流程

**何时使用 Executor**：
- 业务流程涉及多个步骤
- 需要调用多个服务协同完成
- 需要与外部系统交互（短信、邮件、第三方 API）
- 有复杂的条件分支和业务规则
- 需要编排同步和异步操作

**真实案例**：

| 场景 | 使用 Service | 使用 Executor |
|------|-------------|--------------|
| 查询用户信息 | ✅ | ❌ |
| 更新用户资料 | ✅ | ❌ |
| 用户注册 | ❌ | ✅（涉及风控、积分、优惠券） |
| 下单 | ❌ | ✅（库存、价格、优惠、支付） |
| 退款 | ❌ | ✅（退款、库存、积分、通知） |
| 审批流程 | ❌ | ✅（多步骤审批、通知） |

---

## 领域事件机制

领域事件是 DDD 中解耦业务逻辑的核心机制。

### 为什么需要领域事件？

#### 1. 解耦

❌ **不使用事件（紧耦合）**：
```java
public User registerUser(...) {
    User user = userRepository.save(user);

    // 直接调用，强依赖
    emailService.sendWelcome(user.getEmail());
    statisticsService.record(user.getId());
    couponService.grant(user.getId());

    // 新增功能必须修改这里，违反开闭原则

    return user;
}
```

✅ **使用事件（解耦）**：
```java
public User registerUser(...) {
    User user = userRepository.save(user);

    // 发布事件，解耦
    eventPublisher.publishEvent(new UserCreatedEvent(...));

    // 新增功能只需添加新的监听器，无需修改这里

    return user;
}
```

#### 2. 异步处理

事件监听器可以异步执行，不阻塞主流程：

```java
@Async  // 异步执行
@EventListener
public void handleUserCreated(UserCreatedEvent event) {
    // 即使发邮件失败，也不影响用户注册成功
    emailService.sendWelcome(event.getEmail());
}
```

#### 3. 扩展性

一个事件可以有多个监听器，互不影响：

```java
// 监听器1：发邮件
@Component
public class EmailListener {
    @Async
    @EventListener
    public void sendEmail(UserCreatedEvent event) {
        emailService.sendWelcome(event.getEmail());
    }
}

// 监听器2：发优惠券（新增功能，无需修改原有代码）
@Component
public class CouponListener {
    @Async
    @EventListener
    public void grantCoupon(UserCreatedEvent event) {
        couponService.grantNewUserCoupon(event.getUserId());
    }
}

// 监听器3：记录统计（新增功能，无需修改原有代码）
@Component
public class StatisticsListener {
    @Async
    @EventListener
    public void recordStats(UserCreatedEvent event) {
        statisticsService.increment("user_count");
    }
}
```

#### 4. 审计

事件天然就是审计日志，记录了系统中发生的重要业务事件。

### 事件流转全过程

```
【1. 创建事件】
UserDomainService.registerUser()
    ↓
UserCreatedEvent event = new UserCreatedEvent(
    user.getId(),
    user.getUsername(),
    user.getEmail()
);

【2. 发布事件】
eventPublisher.publishEvent(event);
    ↓
Spring ApplicationEventPublisher（事件总线）
    ↓
【3. 扫描监听器】
找到所有标注 @EventListener 且参数类型匹配的方法
    ↓
【4. 调用监听器】
UserEventListener.handleUserCreated(event)
EmailListener.sendEmail(event)
CouponListener.grantCoupon(event)
...（所有监听器都会被调用）
    ↓
【5. 异步执行】（如果有 @Async）
在独立线程中执行业务逻辑
```

### 事件定义

```java
/**
 * 领域事件应该是不可变对象（所有字段 final）
 */
@Getter
public class UserCreatedEvent {
    private final Long userId;
    private final String username;
    private final String email;
    private final LocalDateTime occurredOn;  // 事件发生时间

    public UserCreatedEvent(Long userId, String username, String email) {
        this.userId = userId;
        this.username = username;
        this.email = email;
        this.occurredOn = LocalDateTime.now();
    }
}
```

**设计原则**：
- 所有字段 `final`（不可变）
- 包含事件发生时间
- 包含事件相关的所有必要数据
- 使用业务语言命名（UserCreatedEvent，而非 UserSavedEvent）

### 事件发布

```java
@Service
@RequiredArgsConstructor
public class UserDomainService {
    private final UserRepository userRepository;
    private final ApplicationEventPublisher eventPublisher;  // 注入事件发布器

    public User registerUser(...) {
        User user = userRepository.save(user);

        // 发布事件
        eventPublisher.publishEvent(new UserCreatedEvent(
            user.getId(),
            user.getUsername(),
            user.getEmail()
        ));

        return user;
    }
}
```

**关键点**：
- 注入 `ApplicationEventPublisher`
- 调用 `publishEvent(event)` 发布事件
- `publishEvent()` 是同步调用，但监听器可以异步处理

### 事件监听

```java
@Component
@Slf4j
public class UserEventListener {

    // 实际项目中注入需要的服务
    // private final EmailService emailService;
    // private final CouponService couponService;

    /**
     * 监听用户创建事件
     */
    @Async  // 异步执行，不阻塞主流程
    @EventListener  // 声明这是一个事件监听器
    public void handleUserCreated(UserCreatedEvent event) {
        log.info("处理用户创建事件: {}", event.getUsername());

        try {
            // 业务逻辑1：发送欢迎邮件
            sendWelcomeEmail(event.getEmail(), event.getUsername());

            // 业务逻辑2：记录统计
            recordUserRegistration(event.getUserId());

            // 可扩展：发放优惠券、初始化积分等

        } catch (Exception e) {
            // 异步监听器应该自己处理异常
            log.error("处理事件失败", e);
        }
    }
}
```

**关键注解**：
- `@EventListener`：声明这是一个事件监听器
- 方法参数类型决定监听哪个事件（这里是 `UserCreatedEvent`）
- `@Async`：在独立线程中异步执行

**注意事项**：
1. 监听器应该处理自己的异常，不要让异常传播
2. 异步监听器中的异常不会影响主流程
3. 如果需要强一致性，不要使用 `@Async`
4. 一个事件可以有多个监听器

### 启用异步支持

在启动类上添加 `@EnableAsync`：

```java
@EnableAsync  // 启用异步支持
@SpringBootApplication
@MapperScan("com.example.demo.infrastructure.persistence.mapper")
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

### 同步 vs 异步事件

| 特性 | 同步事件 | 异步事件 |
|------|---------|---------|
| **注解** | 只有 `@EventListener` | `@EventListener` + `@Async` |
| **执行时机** | 立即执行，阻塞主流程 | 在独立线程执行，不阻塞 |
| **事务** | 与发布者在同一事务中 | 在独立事务中执行 |
| **异常处理** | 异常会影响主流程 | 异常不影响主流程 |
| **一致性** | 强一致性 | 最终一致性 |
| **使用场景** | 需要立即完成的关键操作 | 可以延后的辅助操作 |

**选择建议**：
- **同步事件**：数据验证、关键业务规则检查
- **异步事件**：发邮件、发短信、记录日志、统计分析

---

## 对象转换策略

在分层架构中，不同层使用不同的对象模型，需要进行对象转换。

### 为什么需要多种对象模型？

```
Request VO  →  Command  →  Entity  →  DO  →  Database
    ↓           ↓          ↓         ↓
Response VO ← DTO ← Entity ← DO ← Database
```

| 对象类型 | 所在层 | 职责 | 特点 |
|---------|-------|------|------|
| **Request VO** | Adapter | 接收前端请求 | 包含校验注解、符合前端习惯 |
| **Response VO** | Adapter | 返回给前端 | 数据脱敏、格式转换、字段裁剪 |
| **Command** | Application | 携带命令参数 | 不可变对象，表达用户意图 |
| **DTO** | Application | 跨层传输数据 | 简单数据结构，无业务逻辑 |
| **Entity** | Domain | 领域实体 | 包含业务逻辑和业务规则 |
| **DO** | Infrastructure | 数据库映射 | 包含持久化注解，映射表结构 |

### 使用 MapStruct 自动转换

#### 为什么选择 MapStruct？

- **编译时生成**：在编译期生成转换代码，零反射，性能高
- **类型安全**：编译期检查，避免运行时错误
- **易于维护**：不需要手写重复的转换代码
- **功能强大**：支持复杂映射、嵌套对象、集合转换

#### 基础转换

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {

    // Entity → DTO（简单转换）
    UserDTO toDTO(User user);

    // List<Entity> → List<DTO>（集合转换）
    List<UserDTO> toDTOList(List<User> users);
}
```

#### 自定义映射规则

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {

    // Command → Entity
    @Mapping(target = "id", ignore = true)  // 忽略 id 字段
    @Mapping(target = "status", constant = "1")  // 设置默认值
    @Mapping(target = "createTime", ignore = true)  // 由数据库自动填充
    @Mapping(target = "updateTime", ignore = true)
    User toEntity(CreateUserCommand command);
}
```

#### 更新实体（部分字段）

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {

    /**
     * 更新实体（只更新非 null 字段）
     * @MappingTarget 表示目标对象
     */
    @BeanMapping(nullValuePropertyMappingStrategy =
                 NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)  // 不更新 id
    @Mapping(target = "username", ignore = true)  // 不更新用户名
    @Mapping(target = "createTime", ignore = true)
    @Mapping(target = "updateTime", ignore = true)
    void updateEntity(@MappingTarget User user, UpdateUserCommand command);
}
```

使用示例：
```java
User user = userRepository.findById(id).orElseThrow();
UpdateUserCommand command = new UpdateUserCommand();
command.setEmail("new@email.com");
command.setPhone(null);  // phone 为 null，不会更新

userAssembler.updateEntity(user, command);
// user 的 email 被更新，phone 保持不变
```

#### 自定义转换方法

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface ResponseVOAssembler {

    @Mapping(source = "phone", target = "phone", qualifiedByName = "maskPhone")
    @Mapping(source = "status", target = "statusText", qualifiedByName = "statusToText")
    UserResponseVO toUserResponseVO(UserDTO dto);

    /**
     * 自定义方法：手机号脱敏
     */
    @Named("maskPhone")
    default String maskPhone(String phone) {
        if (phone == null || phone.length() < 11) {
            return phone;
        }
        return phone.substring(0, 3) + "****" + phone.substring(7);
    }

    /**
     * 自定义方法：状态码转文本
     */
    @Named("statusToText")
    default String statusToText(Integer status) {
        return status != null && status == 1 ? "启用" : "禁用";
    }
}
```

转换效果：
```java
UserDTO dto = new UserDTO();
dto.setPhone("13812345678");
dto.setStatus(1);

UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
// vo.getPhone() = "138****5678"  （脱敏）
// vo.getStatusText() = "启用"     （转换）
```

#### Controller 中的多次转换

```java
@RestController
@RequiredArgsConstructor
public class UserController {
    private final UserService userService;
    private final UserControllerAssembler assembler;
    private final ResponseVOAssembler responseAssembler;

    @PostMapping
    public Result<UserResponseVO> createUser(
            @Validated @RequestBody CreateUserRequest request) {

        // 第1次转换：Request → Command
        CreateUserCommand command = assembler.toCreateCommand(request);

        // 调用 Service（内部会进行 Command → Entity → DTO 的转换）
        UserDTO dto = userService.createUser(command);

        // 第2次转换：DTO → Response VO（脱敏、格式化）
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);

        return Result.success(vo);
    }
}
```

### 转换层次总结

```
Controller 层：
  Request VO  ──────────→  Command
                         (UserControllerAssembler)
                              ↓
Application 层：
  Command  ───────────→  Entity
                        (UserAssembler)
       Entity  ────────→  DTO
                        (UserAssembler)
                              ↓
Controller 层：
  DTO  ───────────────→  Response VO
                        (ResponseVOAssembler)
                        （数据脱敏、格式转换）

Infrastructure 层：
  Entity  ────────────→  DO
                        (手动转换或 MapStruct)
       DO  ────────────→  Entity
                        (手动转换或 MapStruct)
```

---

## 依赖注入最佳实践

### 使用构造器注入 + Lombok

❌ **不推荐：字段注入**
```java
@Service
public class UserService {
    @Autowired  // 不推荐
    private UserRepository userRepository;

    @Autowired
    private UserAssembler userAssembler;
}
```

问题：
- 无法创建不可变对象
- 不利于单元测试（难以 mock 依赖）
- 可能出现循环依赖
- IDE 警告

✅ **推荐：构造器注入 + Lombok**
```java
@Service
@RequiredArgsConstructor  // Lombok 自动生成构造器
public class UserService {
    private final UserRepository userRepository;  // final 字段
    private final UserAssembler userAssembler;

    // Lombok 自动生成：
    // public UserService(UserRepository userRepository,
    //                    UserAssembler userAssembler) {
    //     this.userRepository = userRepository;
    //     this.userAssembler = userAssembler;
    // }
}
```

优势：
- 依赖明确，容易看出依赖关系
- 支持 `final` 字段，不可变对象
- 利于单元测试（可以直接 new 对象）
- 避免循环依赖
- 代码简洁（Lombok 自动生成构造器）

### Spring Bean 注册

| 注解 | 用途 | 所在层 |
|------|-----|-------|
| `@RestController` | REST 控制器 | Adapter |
| `@Service` | 服务类 | Application / Domain |
| `@Component` | 通用组件 | Application（Executor、Listener） |
| `@Repository` | 数据访问类 | Infrastructure |
| `@Mapper` | MapStruct 转换器 | Application |

**示例**：
```java
// Controller
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController { ... }

// Application Service
@Service
@RequiredArgsConstructor
public class UserService { ... }

// Domain Service
@Service
@RequiredArgsConstructor
public class UserDomainService { ... }

// Executor
@Component  // 使用 @Component，不是 @Service
@RequiredArgsConstructor
public class RegisterUserExecutor { ... }

// Event Listener
@Component  // 使用 @Component
@RequiredArgsConstructor
public class UserEventListener { ... }

// Repository
@Repository
@RequiredArgsConstructor
public class UserRepositoryImpl implements UserRepository { ... }

// MapStruct Assembler
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler { ... }
```

---

## 完整请求流程示例

让我们通过一个完整的用户注册流程，串联所有的概念。

### 场景：用户注册

**业务需求**：
1. 验证用户名唯一性
2. 创建用户账户
3. 发送欢迎邮件
4. 发放新人优惠券
5. 初始化积分账户
6. 记录注册统计

### 请求流程图

```
【1. 前端请求】
POST /api/users
{
  "username": "zhangsan",
  "email": "zhangsan@example.com",
  "phone": "13812345678"
}
    ↓
【2. Adapter Layer】
UserController.createUser(CreateUserRequest)
  → 校验参数（@Validated）
  → Request → Command (UserControllerAssembler)
    ↓
【3. Application Layer】
UserService.createUser(CreateUserCommand)
  → 委托给 RegisterUserExecutor
    ↓
RegisterUserExecutor.execute(CreateUserCommand)
  → 调用 UserDomainService
    ↓
【4. Domain Layer】
UserDomainService.registerUser(username, email, phone)
  → 验证用户名唯一性（领域规则）
  → 创建 User 实体
  → 调用 UserRepository.save(user)
    ↓
【5. Infrastructure Layer】
UserRepositoryImpl.save(User)
  → Entity → DO 转换
  → UserMapper.insert(userDO)
  → 保存到数据库
    ↓
【6. Domain Layer - 发布事件】
UserDomainService.registerUser()
  → eventPublisher.publishEvent(new UserCreatedEvent(...))
    ↓
【7. Application Layer - 事件监听】
Spring 事件总线分发事件
  → UserEventListener.handleUserCreated(event) 【异步】
      → 发送欢迎邮件
      → 记录注册统计
  → CouponListener.grantCoupon(event) 【异步】
      → 发放新人优惠券
  → PointsListener.initAccount(event) 【异步】
      → 初始化积分账户
    ↓
【8. Application Layer - 返回结果】
RegisterUserExecutor.execute()
  → User → UserDTO (UserAssembler)
  → 返回 UserDTO
    ↓
UserService.createUser()
  → 返回 UserDTO
    ↓
【9. Adapter Layer - 响应】
UserController.createUser()
  → DTO → Response VO (ResponseVOAssembler)
      → 手机号脱敏：138****5678
      → 状态转文本：启用
  → Result.success(UserResponseVO)
    ↓
【10. 前端收到响应】
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "phone": "138****5678",
    "statusText": "启用",
    "createTime": "2025-01-15 10:30:00"
  }
}
```

### 代码实现

#### 1. Adapter Layer

```java
// Request VO
@Data
public class CreateUserRequest {
    @NotBlank(message = "用户名不能为空")
    @Length(min = 3, max = 20, message = "用户名长度为3-20个字符")
    private String username;

    @Email(message = "邮箱格式不正确")
    @NotBlank(message = "邮箱不能为空")
    private String email;

    @Pattern(regexp = "^1[3-9]\\d{9}$", message = "手机号格式不正确")
    private String phone;
}

// Controller
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {
    private final UserService userService;
    private final UserControllerAssembler assembler;
    private final ResponseVOAssembler responseAssembler;

    @PostMapping
    public Result<UserResponseVO> createUser(
            @Validated @RequestBody CreateUserRequest request) {

        // Request → Command
        CreateUserCommand command = assembler.toCreateCommand(request);

        // 调用 Application Service
        UserDTO dto = userService.createUser(command);

        // DTO → Response VO（脱敏、格式化）
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);

        return Result.success(vo);
    }
}
```

#### 2. Application Layer

```java
// Command
@Data
public class CreateUserCommand {
    private String username;
    private String email;
    private String phone;
}

// Application Service
@Service
@RequiredArgsConstructor
public class UserService {
    private final RegisterUserExecutor registerUserExecutor;

    @Transactional(rollbackFor = Exception.class)
    public UserDTO createUser(CreateUserCommand command) {
        // 委托给 Executor 处理复杂的注册流程
        return registerUserExecutor.execute(command);
    }
}

// Use Case Executor
@Component
@RequiredArgsConstructor
public class RegisterUserExecutor {
    private final UserDomainService userDomainService;
    private final UserAssembler userAssembler;

    @Transactional(rollbackFor = Exception.class)
    public UserDTO execute(CreateUserCommand command) {
        log.info("执行用户注册用例: {}", command.getUsername());

        // 调用领域服务执行核心业务逻辑
        User user = userDomainService.registerUser(
            command.getUsername(),
            command.getEmail(),
            command.getPhone()
        );

        // Entity → DTO
        return userAssembler.toDTO(user);
    }
}
```

#### 3. Domain Layer

```java
// Domain Entity
@Data
public class User {
    private Long id;
    private String username;
    private String email;
    private String phone;
    private Integer status;  // 0-禁用, 1-启用
    private LocalDateTime createTime;
    private LocalDateTime updateTime;

    // 领域方法
    public void enable() {
        this.status = 1;
    }

    public void disable() {
        this.status = 0;
    }

    public boolean isActive() {
        return this.status == 1;
    }
}

// Domain Service
@Service
@RequiredArgsConstructor
@Slf4j
public class UserDomainService {
    private final UserRepository userRepository;
    private final ApplicationEventPublisher eventPublisher;

    public User registerUser(String username, String email, String phone) {
        log.info("注册新用户: {}", username);

        // 1. 验证领域规则 - 用户名唯一性
        if (userRepository.existsByUsername(username)) {
            throw new IllegalArgumentException("用户名已存在: " + username);
        }

        // 2. 创建用户实体
        User user = new User();
        user.setUsername(username);
        user.setEmail(email);
        user.setPhone(phone);
        user.enable();  // 使用领域方法

        // 3. 持久化
        user = userRepository.save(user);

        // 4. 发布领域事件
        eventPublisher.publishEvent(new UserCreatedEvent(
            user.getId(),
            user.getUsername(),
            user.getEmail()
        ));

        log.info("用户注册成功, id: {}, 事件已发布", user.getId());
        return user;
    }
}

// Domain Event
@Getter
public class UserCreatedEvent {
    private final Long userId;
    private final String username;
    private final String email;
    private final LocalDateTime occurredOn;

    public UserCreatedEvent(Long userId, String username, String email) {
        this.userId = userId;
        this.username = username;
        this.email = email;
        this.occurredOn = LocalDateTime.now();
    }
}
```

#### 4. Infrastructure Layer

```java
// Data Object
@Data
@TableName("t_user")
public class UserDO {
    @TableId(type = IdType.AUTO)
    private Long id;

    private String username;
    private String email;
    private String phone;
    private Integer status;

    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createTime;

    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updateTime;
}

// Repository Implementation
@Repository
@RequiredArgsConstructor
public class UserRepositoryImpl implements UserRepository {
    private final UserMapper userMapper;

    @Override
    public User save(User user) {
        UserDO userDO = toDataObject(user);
        userMapper.insert(userDO);
        user.setId(userDO.getId());
        user.setCreateTime(userDO.getCreateTime());
        user.setUpdateTime(userDO.getUpdateTime());
        return user;
    }

    @Override
    public boolean existsByUsername(String username) {
        return userMapper.selectCount(
            new LambdaQueryWrapper<UserDO>()
                .eq(UserDO::getUsername, username)
        ) > 0;
    }

    private UserDO toDataObject(User user) {
        UserDO userDO = new UserDO();
        userDO.setId(user.getId());
        userDO.setUsername(user.getUsername());
        userDO.setEmail(user.getEmail());
        userDO.setPhone(user.getPhone());
        userDO.setStatus(user.getStatus());
        return userDO;
    }

    private User toEntity(UserDO userDO) {
        User user = new User();
        user.setId(userDO.getId());
        user.setUsername(userDO.getUsername());
        user.setEmail(userDO.getEmail());
        user.setPhone(userDO.getPhone());
        user.setStatus(userDO.getStatus());
        user.setCreateTime(userDO.getCreateTime());
        user.setUpdateTime(userDO.getUpdateTime());
        return user;
    }
}
```

#### 5. Event Listeners（异步处理）

```java
// 邮件监听器
@Component
@Slf4j
public class UserEventListener {

    @Async
    @EventListener
    public void handleUserCreated(UserCreatedEvent event) {
        log.info("处理用户创建事件: {}", event.getUsername());

        try {
            // 发送欢迎邮件
            sendWelcomeEmail(event.getEmail(), event.getUsername());

            // 记录注册统计
            recordUserRegistration(event.getUserId());

        } catch (Exception e) {
            log.error("处理用户创建事件失败", e);
        }
    }

    private void sendWelcomeEmail(String email, String username) {
        log.info("发送欢迎邮件到: {}", email);
        // 实际调用邮件服务
    }

    private void recordUserRegistration(Long userId) {
        log.info("记录用户注册统计, userId: {}", userId);
        // 实际调用统计服务
    }
}

// 优惠券监听器
@Component
@Slf4j
public class CouponEventListener {

    @Async
    @EventListener
    public void grantNewUserCoupon(UserCreatedEvent event) {
        log.info("为新用户发放优惠券, userId: {}", event.getUserId());

        try {
            // 调用优惠券服务
            // couponService.grantNewUserCoupon(event.getUserId());

        } catch (Exception e) {
            log.error("发放优惠券失败", e);
        }
    }
}

// 积分监听器
@Component
@Slf4j
public class PointsEventListener {

    @Async
    @EventListener
    public void initPointsAccount(UserCreatedEvent event) {
        log.info("初始化用户积分账户, userId: {}", event.getUserId());

        try {
            // 调用积分服务，赠送新人积分
            // pointsService.initAccount(event.getUserId(), 100);

        } catch (Exception e) {
            log.error("初始化积分账户失败", e);
        }
    }
}
```

### 时序总结

```
时间轴          同步操作                    异步操作
------+---------------------------+---------------------------
 t0   | Controller 接收请求
      | ↓
 t1   | 参数校验
      | ↓
 t2   | Request → Command
      | ↓
 t3   | Service.createUser()
      | ↓
 t4   | Executor.execute()
      | ↓
 t5   | DomainService.registerUser()
      | ↓
 t6   | 验证用户名唯一性
      | ↓
 t7   | 保存到数据库 ✓
      | ↓
 t8   | 发布事件 ────────────────→ Spring 事件总线接收
      | ↓                          ↓
 t9   | 返回 User Entity           事件分发给监听器
      | ↓                          ↓
t10   | User → DTO                 UserEventListener 【线程1】
      | ↓                          → 发送欢迎邮件
t11   | DTO → Response VO          → 记录统计
      | ↓                          ↓
t12   | 返回给前端 ✓               CouponEventListener 【线程2】
      |                            → 发放优惠券
------+                            ↓
                                   PointsEventListener 【线程3】
                                   → 初始化积分账户
                                   ↓
                                   所有异步任务完成 ✓
```

**关键点**：
- t12 时前端已收到成功响应，用户体验好
- 后续的邮件、优惠券、积分等异步处理，不阻塞主流程
- 即使异步任务失败，也不影响用户注册成功

---

## 设计原则总结

### 1. 单一职责原则（SRP）

每个类只负责一件事：
- **Controller**：接收请求、参数校验、格式转换
- **Service**：简单 CRUD、事务控制
- **Executor**：复杂用例编排
- **Domain Service**：领域业务逻辑
- **Repository**：数据访问

### 2. 开闭原则（OCP）

对扩展开放，对修改关闭：
- 使用**事件**而非直接调用，新增功能只需添加监听器
- 使用**接口**定义契约，实现可替换

### 3. 依赖倒置原则（DIP）

高层不依赖低层，都依赖抽象：
- Domain 定义 `Repository` 接口
- Infrastructure 实现接口
- 依赖方向：`Application → Domain ← Infrastructure`

### 4. 接口隔离原则（ISP）

客户端不应依赖它不使用的接口：
- Repository 接口只定义必要的方法
- 不同的查询需求可以定义不同的查询接口

### 5. 最小知识原则（LoD）

一个对象应该对其他对象有最少的了解：
- Controller 不直接访问 Repository
- Service 不直接访问 Mapper
- 通过层次调用，保持松耦合

### 6. 关注点分离

不同的关注点放在不同的层：
- **业务逻辑** → Domain Layer
- **用例编排** → Application Layer
- **技术实现** → Infrastructure Layer
- **接口适配** → Adapter Layer

### 7. 不要重复自己（DRY）

使用工具和框架减少重复：
- **MapStruct**：自动生成对象转换代码
- **Lombok**：自动生成 getter/setter/constructor
- **Spring**：依赖注入、事务管理、事件机制

---

## 快速参考

### 何时使用什么组件？

| 场景 | 使用组件 |
|------|---------|
| 简单 CRUD | Application Service |
| 复杂用例编排 | Use Case Executor |
| 业务规则验证 | Domain Service |
| 跨聚合根逻辑 | Domain Service + Event |
| 数据访问 | Repository |
| 异步处理 | Event + Listener |
| 对象转换 | MapStruct Assembler |
| 数据脱敏 | Response VO Assembler |

### 对象模型选择

| 场景 | 使用对象 |
|------|---------|
| 接收前端请求 | Request VO |
| 返回给前端 | Response VO（脱敏、格式化） |
| 携带命令参数 | Command |
| 跨层传输 | DTO |
| 业务逻辑处理 | Domain Entity |
| 数据库映射 | Data Object (DO) |

### 注解使用

| 注解 | 用途 |
|------|-----|
| `@RestController` | REST 控制器 |
| `@Service` | Application/Domain Service |
| `@Component` | Executor、Listener |
| `@Repository` | Repository 实现 |
| `@Transactional` | 事务控制 |
| `@EventListener` | 事件监听 |
| `@Async` | 异步执行 |
| `@Mapper` | MapStruct 转换器 |
| `@RequiredArgsConstructor` | 构造器注入 |

---

## 总结

本架构设计遵循领域驱动设计（DDD）的核心理念：

1. **清晰的分层**：Adapter、Application、Domain、Infrastructure 各司其职
2. **依赖倒置**：依赖只能从外向内，Domain 是核心
3. **事件驱动**：使用领域事件解耦业务逻辑
4. **用例驱动**：使用 Executor 编排复杂业务流程
5. **自动化工具**：使用 MapStruct、Lombok 减少重复代码
6. **最佳实践**：构造器注入、不可变对象、异步处理

这套架构适用于中大型项目，能够有效应对业务复杂度的增长，保持代码的可维护性和可扩展性。

---

**版本**：1.0.0
**作者**：phixia team
**更新时间**：2025-01-15
