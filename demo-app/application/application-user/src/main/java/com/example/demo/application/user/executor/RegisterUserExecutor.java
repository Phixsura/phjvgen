package com.example.demo.application.user.executor;

import com.example.demo.application.user.assembler.UserAssembler;
import com.example.demo.application.user.dto.CreateUserCommand;
import com.example.demo.application.user.dto.UserDTO;
import com.example.demo.domain.model.User;
import com.example.demo.domain.service.UserDomainService;
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
