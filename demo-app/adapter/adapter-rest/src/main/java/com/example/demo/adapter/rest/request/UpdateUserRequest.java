package com.example.demo.adapter.rest.request;

import lombok.Data;
import jakarta.validation.constraints.Email;

@Data
public class UpdateUserRequest {

    @Email(message = "邮箱格式不正确")
    private String email;

    private String phone;

    private Integer status;
}
