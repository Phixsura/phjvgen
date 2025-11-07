package com.example.demo.adapter.rest.controller;

import com.example.demo.adapter.rest.assembler.ResponseVOAssembler;
import com.example.demo.adapter.rest.assembler.UserControllerAssembler;
import com.example.demo.adapter.rest.request.CreateUserRequest;
import com.example.demo.adapter.rest.request.UpdateUserRequest;
import com.example.demo.adapter.rest.response.UserResponseVO;
import com.example.demo.application.user.dto.CreateUserCommand;
import com.example.demo.application.user.dto.UpdateUserCommand;
import com.example.demo.application.user.dto.UserDTO;
import com.example.demo.application.user.service.UserService;
import com.example.demo.common.response.Result;
import lombok.RequiredArgsConstructor;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import java.util.List;
import java.util.stream.Collectors;

/**
 * 用户控制器
 */
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;
    private final UserControllerAssembler assembler;
    private final ResponseVOAssembler responseAssembler;

    /**
     * 创建用户
     */
    @PostMapping
    public Result<UserResponseVO> createUser(@Validated @RequestBody CreateUserRequest request) {
        CreateUserCommand command = assembler.toCreateCommand(request);
        UserDTO dto = userService.createUser(command);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 更新用户
     */
    @PutMapping("/{id}")
    public Result<UserResponseVO> updateUser(@PathVariable Long id,
                                             @Validated @RequestBody UpdateUserRequest request) {
        UpdateUserCommand command = assembler.toUpdateCommand(id, request);
        UserDTO dto = userService.updateUser(command);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 查询用户
     */
    @GetMapping("/{id}")
    public Result<UserResponseVO> getUser(@PathVariable Long id) {
        UserDTO dto = userService.getUserById(id);
        UserResponseVO vo = responseAssembler.toUserResponseVO(dto);
        return Result.success(vo);
    }

    /**
     * 查询所有用户
     */
    @GetMapping
    public Result<List<UserResponseVO>> getAllUsers() {
        List<UserDTO> dtos = userService.getAllUsers();
        List<UserResponseVO> vos = dtos.stream()
                .map(responseAssembler::toUserResponseVO)
                .collect(Collectors.toList());
        return Result.success(vos);
    }

    /**
     * 删除用户
     */
    @DeleteMapping("/{id}")
    public Result<Void> deleteUser(@PathVariable Long id) {
        userService.deleteUser(id);
        return Result.success();
    }
}
