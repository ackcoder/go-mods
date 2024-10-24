
### golang 唯一ID三方包比较

[来源](https://mp.weixin.qq.com/s/8UdvCM9udqoRcVmrG03lCg)

| 库             | 特性                 | 有序性 | 长度            | 适用场景         |
| -------------- | -------------------- | ------ | --------------- | ---------------- |
| [UUID][1]      | 全局唯一             | 无序   | 128, 36(string) | 分布式、标识符   |
| [ULID][2]      | 全局唯一             | 有序   | 26(string)      | 日志ID、消息队列 |
| [Snowflake][3] | 全局唯一             | 有序   | 64              | 分布式、自增ID   |
| [ShortID][4]   | 简短唯一、含特殊字符 | 无序   | 7~14(string)    | 短链接、验证码   |
| [XID][5]       | 全局唯一             | 有序   | 20(string)      | 分布式数据库主键 |
| [KSUID][6]     | 全局唯一             | 有序   | 27(string)      | 日志ID、消息队列 |
| [Sonyflake][7] | 全局唯一             | 有序   | 64              | 分布式、日志ID   |


[1]: https://github.com/google/uuid
[2]: https://github.com/oklog/ulid
[3]: https://github.com/bwmarrin/snowflake
[4]: https://github.com/teris-io/shortid
[5]: https://github.com/rs/xid
[6]: https://github.com/segmentio/ksuid
[7]: https://github.com/sony/sonyflake
