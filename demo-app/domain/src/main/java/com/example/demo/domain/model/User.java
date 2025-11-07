package com.example.demo.domain.model;

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
