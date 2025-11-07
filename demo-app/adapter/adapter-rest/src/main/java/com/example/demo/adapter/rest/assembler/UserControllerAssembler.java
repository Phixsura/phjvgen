package com.example.demo.adapter.rest.assembler;

import com.example.demo.adapter.rest.request.CreateUserRequest;
import com.example.demo.adapter.rest.request.UpdateUserRequest;
import com.example.demo.application.user.dto.CreateUserCommand;
import com.example.demo.application.user.dto.UpdateUserCommand;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingConstants;

/**
 * Controller层对象转换器
 * 使用MapStruct自动生成实现代码
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserControllerAssembler {

    /**
     * 转换创建用户请求为命令
     */
    CreateUserCommand toCreateCommand(CreateUserRequest request);

    /**
     * 转换更新用户请求为命令
     * @param id 用户ID
     * @param request 更新请求
     * @return 更新命令
     */
    @Mapping(source = "id", target = "id")
    @Mapping(source = "request.email", target = "email")
    @Mapping(source = "request.phone", target = "phone")
    @Mapping(source = "request.status", target = "status")
    UpdateUserCommand toUpdateCommand(Long id, UpdateUserRequest request);
}
