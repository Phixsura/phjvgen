package com.example.demo.domain.service;

import com.example.demo.domain.event.UserCreatedEvent;
import com.example.demo.domain.model.User;
import com.example.demo.domain.repository.UserRepository;
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
