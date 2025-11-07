package com.example.demo.adapter.rest.response;

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
