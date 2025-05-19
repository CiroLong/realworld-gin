package service

// 职责：
// 实现核心业务逻辑（如用户注册的密码加密、关注逻辑）
// 协调多个 Repository 操作（事务管理）
// 业务规则校验（如用户名唯一性检查）

// 只做事务控制、业务规则校验、调用 repository，不再操作 Gin、HTTP 或 DB。
