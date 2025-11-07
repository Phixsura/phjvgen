package com.example.demo.application.user.dto;

import lombok.Data;

@Data
public class UpdateUserCommand {
    private Long id;
    private String email;
    private String phone;
    private Integer status;
}
