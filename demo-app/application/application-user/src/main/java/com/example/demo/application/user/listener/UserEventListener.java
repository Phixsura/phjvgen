package com.example.demo.application.user.listener;

import com.example.demo.domain.event.UserCreatedEvent;
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
