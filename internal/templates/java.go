package templates

// ApplicationMain is the Spring Boot application main class template
const ApplicationMain = `package {{PACKAGE_NAME}};

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

/**
 * 应用启动类
 */
@EnableAsync
@SpringBootApplication
@MapperScan("{{PACKAGE_NAME}}.infrastructure.persistence.mapper")
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
`

// ResultClass is the Result response wrapper template
const ResultClass = `package {{PACKAGE_NAME}}.common.response;

import lombok.Data;
import java.io.Serializable;

/**
 * 统一响应结果
 */
@Data
public class Result<T> implements Serializable {
    private Integer code;
    private String message;
    private T data;
    private Long timestamp;

    public Result() {
        this.timestamp = System.currentTimeMillis();
    }

    public Result(Integer code, String message, T data) {
        this.code = code;
        this.message = message;
        this.data = data;
        this.timestamp = System.currentTimeMillis();
    }

    public static <T> Result<T> success() {
        return success(null);
    }

    public static <T> Result<T> success(T data) {
        return new Result<>(200, "success", data);
    }

    public static <T> Result<T> success(String message, T data) {
        return new Result<>(200, message, data);
    }

    public static <T> Result<T> fail(String message) {
        return new Result<>(500, message, null);
    }

    public static <T> Result<T> fail(Integer code, String message) {
        return new Result<>(code, message, null);
    }

    public boolean isSuccess() {
        return this.code != null && this.code == 200;
    }
}
`

// BusinessExceptionClass is the business exception template
const BusinessExceptionClass = `package {{PACKAGE_NAME}}.common.exception;

import lombok.Getter;

/**
 * 业务异常
 */
@Getter
public class BusinessException extends RuntimeException {
    private final Integer code;

    public BusinessException(String message) {
        super(message);
        this.code = 500;
    }

    public BusinessException(Integer code, String message) {
        super(message);
        this.code = code;
    }

    public BusinessException(String message, Throwable cause) {
        super(message, cause);
        this.code = 500;
    }
}
`

// ErrorCodeClass is the error code constants template
const ErrorCodeClass = `package {{PACKAGE_NAME}}.common.constant;

public interface ErrorCode {
    int SUCCESS = 200;
    int BAD_REQUEST = 400;
    int UNAUTHORIZED = 401;
    int FORBIDDEN = 403;
    int NOT_FOUND = 404;
    int INTERNAL_ERROR = 500;

    String USER_NOT_FOUND = "用户不存在";
    String USER_ALREADY_EXISTS = "用户已存在";
    String INVALID_PARAMETER = "参数错误";
}
`

// HealthControllerClass is the health check controller template
const HealthControllerClass = `package {{PACKAGE_NAME}}.adapter.rest.controller;

import {{PACKAGE_NAME}}.common.response.Result;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * 健康检查控制器
 */
@RestController
@RequestMapping("/api/health")
public class HealthController {

    @GetMapping
    public Result<String> health() {
        return Result.success("Service is running");
    }
}
`

// UserCreatedEvent is the user created domain event template
const UserCreatedEvent = `package {{PACKAGE_NAME}}.domain.event;

import lombok.Getter;
import java.time.LocalDateTime;

/**
 * 用户创建领域事件
 *
 * 领域事件（Domain Event）用于在领域内部或跨领域传播状态变化
 *
 * 为什么需要领域事件？
 * 1. 解耦：避免直接调用其他模块，降低耦合
 * 2. 异步：事件监听器可以异步处理，提升性能
 * 3. 扩展性：新增功能只需添加新的监听器，不修改原有代码
 * 4. 审计：事件天然就是审计日志
 *
 * 使用示例：
 * 1. 发布事件（在 Domain Service 中）：
 *    eventPublisher.publishEvent(new UserCreatedEvent(...));
 *
 * 2. 监听事件（在 Application 层的 Listener 中）：
 *    @EventListener
 *    @Async
 *    public void handleUserCreated(UserCreatedEvent event) {
 *        // 处理逻辑：发邮件、送优惠券等
 *    }
 *
 * 事件流转过程：
 * UserDomainService（发布）
 *   → ApplicationEventPublisher（Spring事件总线）
 *     → UserEventListener（监听）
 *       → 执行业务逻辑（异步）
 */
@Getter
public class UserCreatedEvent {

    /**
     * 用户ID
     */
    private final Long userId;

    /**
     * 用户名
     */
    private final String username;

    /**
     * 邮箱
     */
    private final String email;

    /**
     * 事件发生时间
     */
    private final LocalDateTime occurredOn;

    public UserCreatedEvent(Long userId, String username, String email) {
        this.userId = userId;
        this.username = username;
        this.email = email;
        this.occurredOn = LocalDateTime.now();
    }
}
`

// UserDomainService is the user domain service template
const UserDomainService = `package {{PACKAGE_NAME}}.domain.service;

import {{PACKAGE_NAME}}.domain.event.UserCreatedEvent;
import {{PACKAGE_NAME}}.domain.model.User;
import {{PACKAGE_NAME}}.domain.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Service;

/**
 * 用户领域服务
 * 处理跨聚合根的领域逻辑
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserDomainService {

    private final UserRepository userRepository;
    private final ApplicationEventPublisher eventPublisher;

    /**
     * 注册新用户（领域逻辑）
     *
     * 包含完整的用户注册领域逻辑和事件发布
     */
    public User registerUser(String username, String email, String phone) {
        log.info("Registering new user: {}", username);

        // 步骤1: 验证领域规则 - 用户名唯一性
        if (userRepository.existsByUsername(username)) {
            throw new IllegalArgumentException("用户名已存在: " + username);
        }

        // 步骤2: 创建用户实体
        User user = new User();
        user.setUsername(username);
        user.setEmail(email);
        user.setPhone(phone);
        user.enable(); // 使用领域方法设置状态

        // 步骤3: 持久化到数据库
        user = userRepository.save(user);

        // 步骤4: 发布领域事件
        // ==================== 事件发布详解 ====================
        // 1. 创建事件对象
        UserCreatedEvent event = new UserCreatedEvent(
            user.getId(),
            user.getUsername(),
            user.getEmail()
        );

        // 2. 通过 Spring 的 ApplicationEventPublisher 发布事件
        //    - Spring 会自动将事件分发给所有监听该事件的 @EventListener
        //    - 这是一个同步调用，但监听器可以使用 @Async 异步处理
        eventPublisher.publishEvent(event);

        // 3. 事件发布后，UserEventListener.handleUserCreated() 会被触发
        //    监听器会异步执行：
        //    - 发送欢迎邮件
        //    - 记录用户注册日志
        //    - 其他后续业务逻辑（送优惠券、初始化积分等）
        //
        // 4. 为什么用事件而不是直接调用？
        //    - 解耦：Domain Service 不需要知道有哪些后续操作
        //    - 扩展：新增功能只需添加新的监听器，无需修改这里
        //    - 异步：不阻塞主流程，提升性能
        //    - 事务：监听器可以在独立的事务中执行
        // ===================================================

        log.info("User registered successfully, id: {}, event published", user.getId());
        return user;
    }

    /**
     * 检查用户是否可以被删除
     */
    public boolean canDelete(Long userId) {
        // 这里可以添加复杂的业务规则
        // 例如：检查用户是否有未完成的订单、是否欠款等
        return true;
    }
}
`

// UserEventListener is the event listener template
const UserEventListener = `package {{PACKAGE_NAME}}.application.user.listener;

import {{PACKAGE_NAME}}.domain.event.UserCreatedEvent;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.EventListener;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Component;

/**
 * 用户事件监听器
 *
 * 作用：监听 Domain 层发布的领域事件，执行 Application 层的业务逻辑
 *
 * 事件接收详解：
 * 1. 当 UserDomainService 调用 eventPublisher.publishEvent(event) 时
 * 2. Spring 会扫描所有标注 @EventListener 的方法
 * 3. 找到参数类型匹配的监听器方法（这里是 UserCreatedEvent）
 * 4. 调用该方法，传入事件对象
 *
 * @Async 注解说明：
 * - 使监听器在独立线程中异步执行
 * - 不阻塞主业务流程（用户注册）
 * - 需要在 Application.java 中添加 @EnableAsync 启用异步支持
 * - 异步执行意味着：即使发邮件失败，也不影响用户注册成功
 *
 * 监听器的职责：
 * - 发送通知（邮件、短信、推送）
 * - 记录日志和统计
 * - 调用其他模块的服务（积分、优惠券等）
 * - 数据同步（同步到 ES、缓存等）
 *
 * 注意事项：
 * 1. 监听器应该处理自己的异常，不要让异常传播出去
 * 2. 异步监听器中的异常不会影响主流程
 * 3. 如果需要强一致性，不要使用 @Async
 * 4. 一个事件可以有多个监听器，它们会依次执行
 */
@Slf4j
@Component
public class UserEventListener {

    // 实际项目中应该注入需要的服务：
    // private final EmailService emailService;
    // private final StatisticsService statisticsService;
    // private final CouponService couponService;

    /**
     * 处理用户创建事件
     *
     * 执行时机：
     * - UserDomainService.registerUser() 中调用 publishEvent() 后
     * - 在独立线程中异步执行（因为有 @Async）
     *
     * 执行内容：
     * - 发送欢迎邮件
     * - 记录注册统计
     * - 可扩展：发放新人优惠券、初始化积分账户等
     *
     * @param event 用户创建事件，包含用户ID、用户名、邮箱等信息
     */
    @Async  // 异步执行，不阻塞主流程
    @EventListener  // 监听 UserCreatedEvent 事件
    public void handleUserCreated(UserCreatedEvent event) {
        log.info("========== 开始处理用户创建事件 ==========");
        log.info("用户ID: {}, 用户名: {}, 邮箱: {}",
            event.getUserId(), event.getUsername(), event.getEmail());

        try {
            // 业务逻辑1: 发送欢迎邮件
            sendWelcomeEmail(event.getEmail(), event.getUsername());

            // 业务逻辑2: 记录用户注册统计
            recordUserRegistration(event.getUserId());

            // 在真实项目中，这里还可以：
            // - 发放新人优惠券: couponService.grantNewUserCoupon(event.getUserId());
            // - 初始化积分账户: pointsService.initAccount(event.getUserId(), 100);
            // - 发送短信通知: smsService.sendWelcome(event.getPhone());
            // - 同步到搜索引擎: elasticsearchService.indexUser(event.getUserId());

            log.info("========== 用户创建事件处理完成 ==========");
        } catch (Exception e) {
            // 异步监听器应该自己处理异常，避免影响其他监听器
            log.error("处理用户创建事件失败, userId: {}", event.getUserId(), e);
            // 可以在这里记录到数据库，方便后续人工处理或重试
        }
    }

    /**
     * 发送欢迎邮件
     *
     * 实际项目中应该注入 EmailService 并调用真实的发邮件接口
     */
    private void sendWelcomeEmail(String email, String username) {
        // 模拟发送邮件
        log.info("→ 发送欢迎邮件到: {} (用户名: {})", email, username);

        // 实际代码示例：
        // EmailTemplate template = new EmailTemplate()
        //     .setTo(email)
        //     .setSubject("欢迎加入我们！")
        //     .setContent("亲爱的 " + username + "，欢迎注册...");
        // emailService.send(template);
    }

    /**
     * 记录用户注册统计
     *
     * 实际项目中应该写入统计表或调用统计服务
     */
    private void recordUserRegistration(Long userId) {
        // 模拟记录统计
        log.info("→ 记录用户注册统计, userId: {}", userId);

        // 实际代码示例：
        // statisticsService.increment("user_registration_count");
        // statisticsService.recordEvent("user_registered", userId);
    }
}
`

// UserEntity is the User domain entity template
const UserEntity = `package {{PACKAGE_NAME}}.domain.model;

import lombok.Data;
import java.time.LocalDateTime;

/**
 * 用户领域实体
 */
@Data
public class User {
    /**
     * 用户ID
     */
    private Long id;

    /**
     * 用户名
     */
    private String username;

    /**
     * 邮箱
     */
    private String email;

    /**
     * 手机号
     */
    private String phone;

    /**
     * 状态：0-禁用，1-启用
     */
    private Integer status;

    /**
     * 创建时间
     */
    private LocalDateTime createTime;

    /**
     * 更新时间
     */
    private LocalDateTime updateTime;

    /**
     * 是否启用
     */
    public boolean isEnabled() {
        return status != null && status == 1;
    }

    /**
     * 启用用户
     */
    public void enable() {
        this.status = 1;
    }

    /**
     * 禁用用户
     */
    public void disable() {
        this.status = 0;
    }
}
`

// UserRepository is the User repository interface template
const UserRepository = `package {{PACKAGE_NAME}}.domain.repository;

import {{PACKAGE_NAME}}.domain.model.User;
import java.util.List;
import java.util.Optional;

/**
 * 用户仓储接口
 */
public interface UserRepository {

    /**
     * 根据ID查找用户
     */
    Optional<User> findById(Long id);

    /**
     * 根据用户名查找用户
     */
    Optional<User> findByUsername(String username);

    /**
     * 查找所有用户
     */
    List<User> findAll();

    /**
     * 保存用户
     */
    User save(User user);

    /**
     * 更新用户
     */
    User update(User user);

    /**
     * 删除用户
     */
    void deleteById(Long id);

    /**
     * 检查用户名是否存在
     */
    boolean existsByUsername(String username);
}
`

// UserDO is the User data object template
const UserDO = `package {{PACKAGE_NAME}}.infrastructure.persistence.dataobject;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;
import java.time.LocalDateTime;

/**
 * 用户数据对象
 */
@Data
@TableName("t_user")
public class UserDO {

    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    @TableField("username")
    private String username;

    @TableField("email")
    private String email;

    @TableField("phone")
    private String phone;

    @TableField("status")
    private Integer status;

    @TableField(value = "create_time", fill = FieldFill.INSERT)
    private LocalDateTime createTime;

    @TableField(value = "update_time", fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updateTime;

    @TableLogic
    @TableField("deleted")
    private Integer deleted;
}
`

// UserMapper is the User mapper interface template
const UserMapper = `package {{PACKAGE_NAME}}.infrastructure.persistence.mapper;

import {{PACKAGE_NAME}}.infrastructure.persistence.dataobject.UserDO;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Mapper;

/**
 * 用户Mapper
 */
@Mapper
public interface UserMapper extends BaseMapper<UserDO> {
}
`

// UserRepositoryImpl is the User repository implementation template
const UserRepositoryImpl = `package {{PACKAGE_NAME}}.infrastructure.persistence.impl;

import {{PACKAGE_NAME}}.domain.model.User;
import {{PACKAGE_NAME}}.domain.repository.UserRepository;
import {{PACKAGE_NAME}}.infrastructure.persistence.dataobject.UserDO;
import {{PACKAGE_NAME}}.infrastructure.persistence.mapper.UserMapper;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Repository;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * 用户仓储实现
 */
@Repository
@RequiredArgsConstructor
public class UserRepositoryImpl implements UserRepository {

    private final UserMapper userMapper;

    @Override
    public Optional<User> findById(Long id) {
        UserDO userDO = userMapper.selectById(id);
        return Optional.ofNullable(userDO).map(this::toEntity);
    }

    @Override
    public Optional<User> findByUsername(String username) {
        LambdaQueryWrapper<UserDO> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserDO::getUsername, username);
        UserDO userDO = userMapper.selectOne(wrapper);
        return Optional.ofNullable(userDO).map(this::toEntity);
    }

    @Override
    public List<User> findAll() {
        return userMapper.selectList(null).stream()
                .map(this::toEntity)
                .collect(Collectors.toList());
    }

    @Override
    public User save(User user) {
        UserDO userDO = toDO(user);
        userMapper.insert(userDO);
        return toEntity(userDO);
    }

    @Override
    public User update(User user) {
        UserDO userDO = toDO(user);
        userMapper.updateById(userDO);
        return toEntity(userDO);
    }

    @Override
    public void deleteById(Long id) {
        userMapper.deleteById(id);
    }

    @Override
    public boolean existsByUsername(String username) {
        LambdaQueryWrapper<UserDO> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserDO::getUsername, username);
        return userMapper.selectCount(wrapper) > 0;
    }

    private User toEntity(UserDO userDO) {
        if (userDO == null) {
            return null;
        }
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

    private UserDO toDO(User user) {
        if (user == null) {
            return null;
        }
        UserDO userDO = new UserDO();
        userDO.setId(user.getId());
        userDO.setUsername(user.getUsername());
        userDO.setEmail(user.getEmail());
        userDO.setPhone(user.getPhone());
        userDO.setStatus(user.getStatus());
        userDO.setCreateTime(user.getCreateTime());
        userDO.setUpdateTime(user.getUpdateTime());
        return userDO;
    }
}
`

// UserDTO is the User DTO template
const UserDTO = `package {{PACKAGE_NAME}}.application.user.dto;

import lombok.Data;
import java.time.LocalDateTime;

@Data
public class UserDTO {
    private Long id;
    private String username;
    private String email;
    private String phone;
    private Integer status;
    private LocalDateTime createTime;
    private LocalDateTime updateTime;
}
`

// CreateUserCommand is the create user command template
const CreateUserCommand = `package {{PACKAGE_NAME}}.application.user.dto;

import lombok.Data;

@Data
public class CreateUserCommand {
    private String username;
    private String email;
    private String phone;
}
`

// UpdateUserCommand is the update user command template
const UpdateUserCommand = `package {{PACKAGE_NAME}}.application.user.dto;

import lombok.Data;

@Data
public class UpdateUserCommand {
    private Long id;
    private String email;
    private String phone;
    private Integer status;
}
`

// UserAssembler is the User assembler template using MapStruct
const UserAssembler = `package {{PACKAGE_NAME}}.application.user.assembler;

import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.domain.model.User;
import org.mapstruct.*;

/**
 * 用户对象转换器
 * 使用MapStruct自动生成实现代码
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {

    /**
     * 领域实体转DTO
     */
    UserDTO toDTO(User user);

    /**
     * 创建命令转领域实体
     * 默认设置状态为启用(1)
     */
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "status", constant = "1")
    @Mapping(target = "createTime", ignore = true)
    @Mapping(target = "updateTime", ignore = true)
    User toEntity(CreateUserCommand command);

    /**
     * 更新命令应用到领域实体
     * 只更新非null字段
     */
    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "username", ignore = true)
    @Mapping(target = "createTime", ignore = true)
    @Mapping(target = "updateTime", ignore = true)
    void updateEntity(@MappingTarget User user, UpdateUserCommand command);
}
`

// RegisterUserExecutor is the register user executor template
const RegisterUserExecutor = `package {{PACKAGE_NAME}}.application.user.executor;

import {{PACKAGE_NAME}}.application.user.assembler.UserAssembler;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.domain.model.User;
import {{PACKAGE_NAME}}.domain.service.UserDomainService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

/**
 * 用户注册用例执行器
 *
 * Executor vs Service 的区别：
 * - Service: 面向数据的 CRUD 操作，薄薄一层，主要做参数校验、事务控制、DTO转换
 * - Executor: 面向业务用例，编排复杂流程，可能涉及：
 *   1. 调用多个 Domain Service（跨聚合根的业务逻辑）
 *   2. 调用多个 Repository（多个数据源）
 *   3. 调用外部服务（如发短信、调用第三方API）
 *   4. 复杂的业务流程编排（有条件分支、循环、重试等）
 *
 * 举例说明：
 * - Service.createUser(): 简单的创建用户（校验参数 -> 保存 -> 返回）
 * - RegisterUserExecutor.execute(): 完整的注册流程
 *   1. 调用 UserDomainService 验证业务规则并创建用户
 *   2. 发布领域事件（UserCreatedEvent）
 *   3. 异步发送欢迎邮件（通过事件监听器）
 *   4. 异步赠送新人优惠券（通过事件监听器）
 *   5. 异步记录用户行为日志（通过事件监听器）
 *
 * 在更复杂的场景中，Executor 可能还会：
 * - 调用风控服务检查用户是否存在风险
 * - 调用实名认证服务
 * - 调用积分服务初始化积分账户
 * - 调用会员服务创建会员档案
 * - 根据用户来源渠道分配不同的优惠
 *
 * 这些复杂的业务编排不应该放在 Service 中，Service 应该保持简单
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class RegisterUserExecutor {

    private final UserDomainService userDomainService;
    private final UserAssembler userAssembler;

    // 在真实项目中，这里可能还会注入：
    // private final CouponService couponService;          // 优惠券服务
    // private final RiskControlService riskControlService; // 风控服务
    // private final SmsService smsService;                // 短信服务
    // private final PointsService pointsService;          // 积分服务

    /**
     * 执行完整的用户注册业务流程
     *
     * 这是一个完整的业务用例，包含了注册相关的所有业务逻辑编排
     * 如果只是简单的 CRUD，直接用 Service 就够了
     * 但注册是一个复杂的业务场景，需要编排多个步骤
     */
    @Transactional(rollbackFor = Exception.class)
    public UserDTO execute(CreateUserCommand command) {
        log.info("执行用户注册用例: {}", command.getUsername());

        // 步骤1: 调用领域服务执行核心领域逻辑
        // - 验证用户名唯一性（领域规则）
        // - 创建用户实体
        // - 发布 UserCreatedEvent 领域事件
        User user = userDomainService.registerUser(
            command.getUsername(),
            command.getEmail(),
            command.getPhone()
        );

        // 步骤2: 后续业务编排由事件监听器异步完成：
        // - UserEventListener 监听到 UserCreatedEvent 后：
        //   a) 发送欢迎邮件
        //   b) 记录用户注册日志
        //
        // 在更复杂的项目中，这里还可以添加同步的业务逻辑：
        // - 调用风控服务验证（同步，必须等待结果）
        // - 调用实名认证服务（同步）
        // - 初始化积分账户（可以异步）
        // - 发放新人优惠券（可以异步）

        // 示例：如果需要风控验证（实际项目中的代码）
        // RiskCheckResult riskResult = riskControlService.checkNewUser(user);
        // if (riskResult.isHighRisk()) {
        //     user.markAsRisky();
        //     userRepository.update(user);
        //     throw new BusinessException("用户存在风险，注册失败");
        // }

        log.info("用户注册用例执行完成, userId: {}", user.getId());
        return userAssembler.toDTO(user);
    }
}
`

// UserService is the User service template
const UserService = `package {{PACKAGE_NAME}}.application.user.service;

import {{PACKAGE_NAME}}.application.user.assembler.UserAssembler;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.application.user.executor.RegisterUserExecutor;
import {{PACKAGE_NAME}}.common.constant.ErrorCode;
import {{PACKAGE_NAME}}.common.exception.BusinessException;
import {{PACKAGE_NAME}}.domain.model.User;
import {{PACKAGE_NAME}}.domain.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 用户应用服务（Application Service）
 *
 * Service 的职责定位：
 * 1. 提供面向数据的 CRUD 操作接口
 * 2. 作为 Controller 和 Domain 之间的薄薄一层
 * 3. 主要职责：
 *    - 参数校验（基础校验，复杂业务规则在 Domain Service）
 *    - 事务控制（@Transactional）
 *    - DTO 和 Entity 之间的转换
 *    - 调用 Repository 进行数据持久化
 *    - 对于复杂用例，委托给 Executor 执行
 *
 * 何时使用 Service vs Executor：
 * - 简单 CRUD：直接在 Service 中完成（如 getUserById、updateUser）
 * - 复杂用例：委托给 Executor（如 createUser 委托给 RegisterUserExecutor）
 *
 * Service 应该保持"薄"：
 * - 不要在 Service 中编排复杂的业务流程
 * - 不要在 Service 中调用多个外部服务
 * - 不要在 Service 中包含复杂的业务逻辑
 * - 这些都应该由 Executor 或 Domain Service 负责
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;
    private final UserAssembler userAssembler;
    private final RegisterUserExecutor registerUserExecutor;

    /**
     * 创建用户
     *
     * 注意：这是一个复杂的业务用例（注册），所以委托给 Executor
     * 如果只是简单的创建，可以直接在这里完成
     */
    @Transactional(rollbackFor = Exception.class)
    public UserDTO createUser(CreateUserCommand command) {
        log.info("Creating user: {}", command.getUsername());

        // 委托给 Executor 处理完整的注册流程
        // Executor 会调用 DomainService、发布事件、编排业务流程
        return registerUserExecutor.execute(command);
    }

    /**
     * 更新用户
     *
     * 这是一个简单的 CRUD 操作，直接在 Service 中完成
     * 不需要复杂的业务编排，所以不需要 Executor
     */
    @Transactional(rollbackFor = Exception.class)
    public UserDTO updateUser(UpdateUserCommand command) {
        log.info("Updating user: {}", command.getId());

        // 1. 查询用户
        User user = userRepository.findById(command.getId())
                .orElseThrow(() -> new BusinessException(ErrorCode.NOT_FOUND, ErrorCode.USER_NOT_FOUND));

        // 2. 更新实体（使用 MapStruct 自动映射）
        userAssembler.updateEntity(user, command);

        // 3. 保存
        user = userRepository.update(user);

        log.info("User updated successfully, id: {}", user.getId());
        return userAssembler.toDTO(user);
    }

    /**
     * 根据ID查询用户
     *
     * 简单的查询操作，直接在 Service 中完成
     */
    public UserDTO getUserById(Long id) {
        log.info("Getting user by id: {}", id);

        User user = userRepository.findById(id)
                .orElseThrow(() -> new BusinessException(ErrorCode.NOT_FOUND, ErrorCode.USER_NOT_FOUND));

        return userAssembler.toDTO(user);
    }

    /**
     * 查询所有用户
     *
     * 简单的查询操作，直接在 Service 中完成
     */
    public List<UserDTO> getAllUsers() {
        log.info("Getting all users");

        return userRepository.findAll().stream()
                .map(userAssembler::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * 删除用户
     *
     * 这是一个简单的删除操作
     * 如果删除用户涉及复杂的业务逻辑（如：删除前需要检查订单、积分清零、通知等）
     * 则应该创建一个 DeleteUserExecutor 来编排这些流程
     */
    @Transactional(rollbackFor = Exception.class)
    public void deleteUser(Long id) {
        log.info("Deleting user: {}", id);

        if (!userRepository.findById(id).isPresent()) {
            throw new BusinessException(ErrorCode.NOT_FOUND, ErrorCode.USER_NOT_FOUND);
        }

        userRepository.deleteById(id);
        log.info("User deleted successfully, id: {}", id);
    }
}
`

// CreateUserRequest is the create user request template
const CreateUserRequest = `package {{PACKAGE_NAME}}.adapter.rest.request;

import lombok.Data;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Email;

@Data
public class CreateUserRequest {

    @NotBlank(message = "用户名不能为空")
    private String username;

    @Email(message = "邮箱格式不正确")
    private String email;

    private String phone;
}
`

// UpdateUserRequest is the update user request template
const UpdateUserRequest = `package {{PACKAGE_NAME}}.adapter.rest.request;

import lombok.Data;
import jakarta.validation.constraints.Email;

@Data
public class UpdateUserRequest {

    @Email(message = "邮箱格式不正确")
    private String email;

    private String phone;

    private Integer status;
}
`

// GlobalExceptionHandler is the global exception handler template
const GlobalExceptionHandler = `package {{PACKAGE_NAME}}.adapter.rest.advice;

import {{PACKAGE_NAME}}.common.exception.BusinessException;
import {{PACKAGE_NAME}}.common.response.Result;
import lombok.extern.slf4j.Slf4j;
import org.springframework.validation.BindException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(BusinessException.class)
    public Result<?> handleBusinessException(BusinessException e) {
        log.error("Business exception: {}", e.getMessage());
        return Result.fail(e.getCode(), e.getMessage());
    }

    @ExceptionHandler({MethodArgumentNotValidException.class, BindException.class})
    public Result<?> handleValidationException(Exception e) {
        FieldError fieldError;
        if (e instanceof MethodArgumentNotValidException) {
            fieldError = ((MethodArgumentNotValidException) e).getBindingResult().getFieldError();
        } else {
            fieldError = ((BindException) e).getBindingResult().getFieldError();
        }
        String message = fieldError != null ? fieldError.getDefaultMessage() : "参数校验失败";
        log.error("Validation exception: {}", message);
        return Result.fail(400, message);
    }

    @ExceptionHandler(Exception.class)
    public Result<?> handleException(Exception e) {
        log.error("Unexpected exception", e);
        return Result.fail("系统异常，请稍后重试");
    }
}
`

// UserControllerAssembler is the controller assembler template using MapStruct
const UserControllerAssembler = `package {{PACKAGE_NAME}}.adapter.rest.assembler;

import {{PACKAGE_NAME}}.adapter.rest.request.CreateUserRequest;
import {{PACKAGE_NAME}}.adapter.rest.request.UpdateUserRequest;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingConstants;

/**
 * Controller层对象转换器
 * 使用MapStruct自动生成实现代码
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserControllerAssembler {

    /**
     * 转换创建用户请求为命令
     */
    CreateUserCommand toCreateCommand(CreateUserRequest request);

    /**
     * 转换更新用户请求为命令
     * @param id 用户ID
     * @param request 更新请求
     * @return 更新命令
     */
    @Mapping(source = "id", target = "id")
    @Mapping(source = "request.email", target = "email")
    @Mapping(source = "request.phone", target = "phone")
    @Mapping(source = "request.status", target = "status")
    UpdateUserCommand toUpdateCommand(Long id, UpdateUserRequest request);
}
`

// UserResponseVO is the user response VO template
const UserResponseVO = `package {{PACKAGE_NAME}}.adapter.rest.response;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Data;
import java.time.LocalDateTime;

/**
 * 用户响应VO
 * 专门用于返回给前端的视图对象
 */
@Data
public class UserResponseVO {

    private Long id;

    private String username;

    private String email;

    /**
     * 手机号脱敏显示
     */
    private String phone;

    private String statusText;

    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createTime;

    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime updateTime;
}
`

// LoggingFilter is the logging filter template
const LoggingFilter = `package {{PACKAGE_NAME}}.adapter.rest.filter;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import java.io.IOException;

/**
 * 请求日志过滤器
 */
@Slf4j
@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class LoggingFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        HttpServletRequest httpRequest = (HttpServletRequest) request;
        long startTime = System.currentTimeMillis();

        log.info("Request started: {} {}", httpRequest.getMethod(), httpRequest.getRequestURI());

        try {
            chain.doFilter(request, response);
        } finally {
            long duration = System.currentTimeMillis() - startTime;
            log.info("Request completed: {} {} in {}ms",
                httpRequest.getMethod(), httpRequest.getRequestURI(), duration);
        }
    }
}
`

// AuthInterceptor is the authentication interceptor template
const AuthInterceptor = `package {{PACKAGE_NAME}}.adapter.rest.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

/**
 * 认证拦截器示例
 */
@Slf4j
@Component
public class AuthInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        String token = request.getHeader("Authorization");

        // 健康检查接口跳过认证
        if (request.getRequestURI().startsWith("/api/health")) {
            return true;
        }

        // 这里可以添加实际的认证逻辑
        log.debug("Auth token: {}", token);

        return true;
    }
}
`

// WebMvcConfig is the web mvc configuration template
const WebMvcConfig = `package {{PACKAGE_NAME}}.adapter.rest.config;

import {{PACKAGE_NAME}}.adapter.rest.interceptor.AuthInterceptor;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

/**
 * Web MVC 配置
 */
@Configuration
@RequiredArgsConstructor
public class WebMvcConfig implements WebMvcConfigurer {

    private final AuthInterceptor authInterceptor;

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(authInterceptor)
                .addPathPatterns("/api/**")
                .excludePathPatterns("/api/health");
    }

    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/api/**")
                .allowedOrigins("*")
                .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
                .allowedHeaders("*")
                .maxAge(3600);
    }
}
`

// ResponseVOAssembler is the response VO assembler template
const ResponseVOAssembler = `package {{PACKAGE_NAME}}.adapter.rest.assembler;

import {{PACKAGE_NAME}}.adapter.rest.response.UserResponseVO;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingConstants;
import org.mapstruct.Named;

/**
 * 响应VO转换器
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface ResponseVOAssembler {

    /**
     * DTO转ResponseVO
     * 可以在这里进行数据脱敏、格式化等
     */
    @Mapping(source = "phone", target = "phone", qualifiedByName = "maskPhone")
    @Mapping(source = "status", target = "statusText", qualifiedByName = "statusToText")
    UserResponseVO toUserResponseVO(UserDTO dto);

    /**
     * 手机号脱敏
     */
    @Named("maskPhone")
    default String maskPhone(String phone) {
        if (phone == null || phone.length() < 11) {
            return phone;
        }
        return phone.substring(0, 3) + "****" + phone.substring(7);
    }

    /**
     * 状态码转文本
     */
    @Named("statusToText")
    default String statusToText(Integer status) {
        return status != null && status == 1 ? "启用" : "禁用";
    }
}
`

// UserController is the User controller template
const UserController = `package {{PACKAGE_NAME}}.adapter.rest.controller;

import {{PACKAGE_NAME}}.adapter.rest.assembler.ResponseVOAssembler;
import {{PACKAGE_NAME}}.adapter.rest.assembler.UserControllerAssembler;
import {{PACKAGE_NAME}}.adapter.rest.request.CreateUserRequest;
import {{PACKAGE_NAME}}.adapter.rest.request.UpdateUserRequest;
import {{PACKAGE_NAME}}.adapter.rest.response.UserResponseVO;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.application.user.service.UserService;
import {{PACKAGE_NAME}}.common.response.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 用户控制器
 */
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;
    private final UserControllerAssembler assembler;
    private final ResponseVOAssembler responseAssembler;

    /**
     * 创建用户
     */
    @PostMapping
    public Result<UserResponseVO> createUser(@Validated @RequestBody CreateUserRequest request) {
        CreateUserCommand command = assembler.toCreateCommand(request);
        UserDTO dto = userService.createUser(command);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 更新用户
     */
    @PutMapping("/{id}")
    public Result<UserResponseVO> updateUser(@PathVariable Long id,
                                             @Validated @RequestBody UpdateUserRequest request) {
        UpdateUserCommand command = assembler.toUpdateCommand(id, request);
        UserDTO dto = userService.updateUser(command);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 查询用户
     */
    @GetMapping("/{id}")
    public Result<UserResponseVO> getUser(@PathVariable Long id) {
        UserDTO dto = userService.getUserById(id);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 查询所有用户
     */
    @GetMapping
    public Result<List<UserResponseVO>> getAllUsers() {
        List<UserDTO> dtos = userService.getAllUsers();
        List<UserResponseVO> vos = dtos.stream()
                .map(responseAssembler::toUserResponseVO)
                .collect(Collectors.toList());
        return Result.success(vos);
    }

    /**
     * 删除用户
     */
    @DeleteMapping("/{id}")
    public Result<Void> deleteUser(@PathVariable Long id) {
        userService.deleteUser(id);
        return Result.success();
    }
}
`
