package com.example.demo.domain.event;

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
