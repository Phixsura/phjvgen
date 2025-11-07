package com.example.demo.application.user.service;

import com.example.demo.application.user.assembler.UserAssembler;
import com.example.demo.application.user.dto.CreateUserCommand;
import com.example.demo.application.user.dto.UpdateUserCommand;
import com.example.demo.application.user.dto.UserDTO;
import com.example.demo.application.user.executor.RegisterUserExecutor;
import com.example.demo.common.constant.ErrorCode;
import com.example.demo.common.exception.BusinessException;
import com.example.demo.domain.model.User;
import com.example.demo.domain.repository.UserRepository;
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
