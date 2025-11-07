package com.example.demo.adapter.rest.advice;

import com.example.demo.common.exception.BusinessException;
import com.example.demo.common.response.Result;
import lombok.extern.slf4j.Slf4j;
import org.springframework.validation.BindException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(BusinessException.class)
    public Result<?> handleBusinessException(BusinessException e) {
        log.error("Business exception: {}", e.getMessage());
        return Result.fail(e.getCode(), e.getMessage());
    }

    @ExceptionHandler({MethodArgumentNotValidException.class, BindException.class})
    public Result<?> handleValidationException(Exception e) {
        FieldError fieldError;
        if (e instanceof MethodArgumentNotValidException) {
            fieldError = ((MethodArgumentNotValidException) e).getBindingResult().getFieldError();
        } else {
            fieldError = ((BindException) e).getBindingResult().getFieldError();
        }
        String message = fieldError != null ? fieldError.getDefaultMessage() : "参数校验失败";
        log.error("Validation exception: {}", message);
        return Result.fail(400, message);
    }

    @ExceptionHandler(Exception.class)
    public Result<?> handleException(Exception e) {
        log.error("Unexpected exception", e);
        return Result.fail("系统异常，请稍后重试");
    }
}
