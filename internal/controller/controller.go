package controller

// 职责：
//
// 接收 HTTP 请求，解析参数（Path/Query/JSON）
// 调用 Service 层处理业务逻辑
// 处理 HTTP 响应（成功/错误格式统一）
// 请求参数校验（结合 go-playground/validator）

// 只做请求参数绑定、调用 service、返回统一格式
