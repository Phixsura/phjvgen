package com.example.demo.infrastructure.persistence.impl;

import com.example.demo.domain.model.User;
import com.example.demo.domain.repository.UserRepository;
import com.example.demo.infrastructure.persistence.dataobject.UserDO;
import com.example.demo.infrastructure.persistence.mapper.UserMapper;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Repository;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * 用户仓储实现
 */
@Repository
@RequiredArgsConstructor
public class UserRepositoryImpl implements UserRepository {

    private final UserMapper userMapper;

    @Override
    public Optional<User> findById(Long id) {
        UserDO userDO = userMapper.selectById(id);
        return Optional.ofNullable(userDO).map(this::toEntity);
    }

    @Override
    public Optional<User> findByUsername(String username) {
        LambdaQueryWrapper<UserDO> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserDO::getUsername, username);
        UserDO userDO = userMapper.selectOne(wrapper);
        return Optional.ofNullable(userDO).map(this::toEntity);
    }

    @Override
    public List<User> findAll() {
        return userMapper.selectList(null).stream()
                .map(this::toEntity)
                .collect(Collectors.toList());
    }

    @Override
    public User save(User user) {
        UserDO userDO = toDO(user);
        userMapper.insert(userDO);
        return toEntity(userDO);
    }

    @Override
    public User update(User user) {
        UserDO userDO = toDO(user);
        userMapper.updateById(userDO);
        return toEntity(userDO);
    }

    @Override
    public void deleteById(Long id) {
        userMapper.deleteById(id);
    }

    @Override
    public boolean existsByUsername(String username) {
        LambdaQueryWrapper<UserDO> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserDO::getUsername, username);
        return userMapper.selectCount(wrapper) > 0;
    }

    private User toEntity(UserDO userDO) {
        if (userDO == null) {
            return null;
        }
        User user = new User();
        user.setId(userDO.getId());
        user.setUsername(userDO.getUsername());
        user.setEmail(userDO.getEmail());
        user.setPhone(userDO.getPhone());
        user.setStatus(userDO.getStatus());
        user.setCreateTime(userDO.getCreateTime());
        user.setUpdateTime(userDO.getUpdateTime());
        return user;
    }

    private UserDO toDO(User user) {
        if (user == null) {
            return null;
        }
        UserDO userDO = new UserDO();
        userDO.setId(user.getId());
        userDO.setUsername(user.getUsername());
        userDO.setEmail(user.getEmail());
        userDO.setPhone(user.getPhone());
        userDO.setStatus(user.getStatus());
        userDO.setCreateTime(user.getCreateTime());
        userDO.setUpdateTime(user.getUpdateTime());
        return userDO;
    }
}
