package com.example.demo.application.user.dto;

import lombok.Data;

@Data
public class CreateUserCommand {
    private String username;
    private String email;
    private String phone;
}
