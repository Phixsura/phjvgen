package com.example.demo.adapter.rest.assembler;

import com.example.demo.adapter.rest.response.UserResponseVO;
import com.example.demo.application.user.dto.UserDTO;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingConstants;
import org.mapstruct.Named;

/**
 * 响应VO转换器
 */
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface ResponseVOAssembler {

    /**
     * DTO转ResponseVO
     * 可以在这里进行数据脱敏、格式化等
     */
    @Mapping(source = "phone", target = "phone", qualifiedByName = "maskPhone")
    @Mapping(source = "status", target = "statusText", qualifiedByName = "statusToText")
    UserResponseVO toUserResponseVO(UserDTO dto);

    /**
     * 手机号脱敏
     */
    @Named("maskPhone")
    default String maskPhone(String phone) {
        if (phone == null || phone.length() < 11) {
            return phone;
        }
        return phone.substring(0, 3) + "****" + phone.substring(7);
    }

    /**
     * 状态码转文本
     */
    @Named("statusToText")
    default String statusToText(Integer status) {
        return status != null && status == 1 ? "启用" : "禁用";
    }
}
