package com.example.demo.common.constant;

public interface ErrorCode {
    int SUCCESS = 200;
    int BAD_REQUEST = 400;
    int UNAUTHORIZED = 401;
    int FORBIDDEN = 403;
    int NOT_FOUND = 404;
    int INTERNAL_ERROR = 500;

    String USER_NOT_FOUND = "用户不存在";
    String USER_ALREADY_EXISTS = "用户已存在";
    String INVALID_PARAMETER = "参数错误";
}
