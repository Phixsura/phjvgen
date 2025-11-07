package com.example.demo.application.user.assembler;

import com.example.demo.application.user.dto.CreateUserCommand;
import com.example.demo.application.user.dto.UpdateUserCommand;
import com.example.demo.application.user.dto.UserDTO;
import com.example.demo.domain.model.User;
import org.mapstruct.*;

/**
 * 用户对象转换器
 * 使用MapStruct自动生成实现代码
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserAssembler {

    /**
     * 领域实体转DTO
     */
    UserDTO toDTO(User user);

    /**
     * 创建命令转领域实体
     * 默认设置状态为启用(1)
     */
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "status", constant = "1")
    @Mapping(target = "createTime", ignore = true)
    @Mapping(target = "updateTime", ignore = true)
    User toEntity(CreateUserCommand command);

    /**
     * 更新命令应用到领域实体
     * 只更新非null字段
     */
    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "username", ignore = true)
    @Mapping(target = "createTime", ignore = true)
    @Mapping(target = "updateTime", ignore = true)
    void updateEntity(@MappingTarget User user, UpdateUserCommand command);
}
