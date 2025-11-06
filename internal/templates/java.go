package templates

// ApplicationMain is the Spring Boot application main class template
const ApplicationMain = `package {{PACKAGE_NAME}};

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * 应用启动类
 */
@SpringBootApplication
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

// UserAssembler is the User assembler template
const UserAssembler = `package {{PACKAGE_NAME}}.application.user.assembler;

import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.domain.model.User;
import org.springframework.stereotype.Component;

/**
 * 用户对象转换器
 */
@Component
public class UserAssembler {

    public UserDTO toDTO(User user) {
        if (user == null) {
            return null;
        }
        UserDTO dto = new UserDTO();
        dto.setId(user.getId());
        dto.setUsername(user.getUsername());
        dto.setEmail(user.getEmail());
        dto.setPhone(user.getPhone());
        dto.setStatus(user.getStatus());
        dto.setCreateTime(user.getCreateTime());
        dto.setUpdateTime(user.getUpdateTime());
        return dto;
    }

    public User toEntity(CreateUserCommand command) {
        if (command == null) {
            return null;
        }
        User user = new User();
        user.setUsername(command.getUsername());
        user.setEmail(command.getEmail());
        user.setPhone(command.getPhone());
        user.setStatus(1); // 默认启用
        return user;
    }

    public void updateEntity(User user, UpdateUserCommand command) {
        if (command.getEmail() != null) {
            user.setEmail(command.getEmail());
        }
        if (command.getPhone() != null) {
            user.setPhone(command.getPhone());
        }
        if (command.getStatus() != null) {
            user.setStatus(command.getStatus());
        }
    }
}
`

// UserService is the User service template
const UserService = `package {{PACKAGE_NAME}}.application.user.service;

import {{PACKAGE_NAME}}.application.user.assembler.UserAssembler;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
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
 * 用户业务服务
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;
    private final UserAssembler userAssembler;

    /**
     * 创建用户
     */
    @Transactional(rollbackFor = Exception.class)
    public UserDTO createUser(CreateUserCommand command) {
        log.info("Creating user: {}", command.getUsername());

        // 检查用户名是否已存在
        if (userRepository.existsByUsername(command.getUsername())) {
            throw new BusinessException(ErrorCode.BAD_REQUEST, ErrorCode.USER_ALREADY_EXISTS);
        }

        // 转换并保存
        User user = userAssembler.toEntity(command);
        user = userRepository.save(user);

        log.info("User created successfully, id: {}", user.getId());
        return userAssembler.toDTO(user);
    }

    /**
     * 更新用户
     */
    @Transactional(rollbackFor = Exception.class)
    public UserDTO updateUser(UpdateUserCommand command) {
        log.info("Updating user: {}", command.getId());

        User user = userRepository.findById(command.getId())
                .orElseThrow(() -> new BusinessException(ErrorCode.NOT_FOUND, ErrorCode.USER_NOT_FOUND));

        userAssembler.updateEntity(user, command);
        user = userRepository.update(user);

        log.info("User updated successfully, id: {}", user.getId());
        return userAssembler.toDTO(user);
    }

    /**
     * 根据ID查询用户
     */
    public UserDTO getUserById(Long id) {
        log.info("Getting user by id: {}", id);

        User user = userRepository.findById(id)
                .orElseThrow(() -> new BusinessException(ErrorCode.NOT_FOUND, ErrorCode.USER_NOT_FOUND));

        return userAssembler.toDTO(user);
    }

    /**
     * 查询所有用户
     */
    public List<UserDTO> getAllUsers() {
        log.info("Getting all users");

        return userRepository.findAll().stream()
                .map(userAssembler::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * 删除用户
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

// UserController is the User controller template
const UserController = `package {{PACKAGE_NAME}}.adapter.rest.controller;

import {{PACKAGE_NAME}}.adapter.rest.request.CreateUserRequest;
import {{PACKAGE_NAME}}.adapter.rest.request.UpdateUserRequest;
import {{PACKAGE_NAME}}.application.user.dto.CreateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UpdateUserCommand;
import {{PACKAGE_NAME}}.application.user.dto.UserDTO;
import {{PACKAGE_NAME}}.application.user.service.UserService;
import {{PACKAGE_NAME}}.common.response.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import java.util.List;

/**
 * 用户控制器
 */
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    /**
     * 创建用户
     */
    @PostMapping
    public Result<UserDTO> createUser(@Validated @RequestBody CreateUserRequest request) {
        CreateUserCommand command = new CreateUserCommand();
        command.setUsername(request.getUsername());
        command.setEmail(request.getEmail());
        command.setPhone(request.getPhone());

        UserDTO user = userService.createUser(command);
        return Result.success(user);
    }

    /**
     * 更新用户
     */
    @PutMapping("/{id}")
    public Result<UserDTO> updateUser(@PathVariable Long id,
                                      @Validated @RequestBody UpdateUserRequest request) {
        UpdateUserCommand command = new UpdateUserCommand();
        command.setId(id);
        command.setEmail(request.getEmail());
        command.setPhone(request.getPhone());
        command.setStatus(request.getStatus());

        UserDTO user = userService.updateUser(command);
        return Result.success(user);
    }

    /**
     * 查询用户
     */
    @GetMapping("/{id}")
    public Result<UserDTO> getUser(@PathVariable Long id) {
        UserDTO user = userService.getUserById(id);
        return Result.success(user);
    }

    /**
     * 查询所有用户
     */
    @GetMapping
    public Result<List<UserDTO>> getAllUsers() {
        List<UserDTO> users = userService.getAllUsers();
        return Result.success(users);
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
