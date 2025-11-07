package com.example.demo.domain.repository;

import com.example.demo.domain.model.User;
import java.util.List;
import java.util.Optional;

/**
 * 用户仓储接口
 */
public interface UserRepository {

    /**
     * 根据ID查找用户
     */
    Optional<User> findById(Long id);

    /**
     * 根据用户名查找用户
     */
    Optional<User> findByUsername(String username);

    /**
     * 查找所有用户
     */
    List<User> findAll();

    /**
     * 保存用户
     */
    User save(User user);

    /**
     * 更新用户
     */
    User update(User user);

    /**
     * 删除用户
     */
    void deleteById(Long id);

    /**
     * 检查用户名是否存在
     */
    boolean existsByUsername(String username);
}
