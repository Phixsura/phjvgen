package com.example.demo.adapter.rest.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

/**
 * 认证拦截器示例
 */
@Slf4j
@Component
public class AuthInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        String token = request.getHeader("Authorization");

        // 健康检查接口跳过认证
        if (request.getRequestURI().startsWith("/api/health")) {
            return true;
        }

        // 这里可以添加实际的认证逻辑
        log.debug("Auth token: {}", token);

        return true;
    }
}
