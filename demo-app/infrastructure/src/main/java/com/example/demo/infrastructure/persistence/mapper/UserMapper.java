package com.example.demo.infrastructure.persistence.mapper;

import com.example.demo.infrastructure.persistence.dataobject.UserDO;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Mapper;

/**
 * 用户Mapper
 */
@Mapper
public interface UserMapper extends BaseMapper<UserDO> {
}
