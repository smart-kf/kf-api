### API 接入

## 基础路由

```shell
    https://api.smartkf.top/api
```

## 基本请求

* 所有Post请求都是： json 格式，header 头设置 `Content-Type: application/json`
* 鉴权接口的请求方式为：需要在 header 中设置：`Authorization:Bearer $jwtToken` 或者 `Authorization: $jwtToken` 

## 基本响应格式

* http响应码都为: `200`
* 所有响应格式都是: `{"msg":"错误提示内容","code":200,"data": { key:value },"debug_info": "debug信息，前端勿用"}`
* 错误码：

|  错误码  |              错误提示               |
|:-----:|:-------------------------------:|
|  200  |              请求成功               |
|  400  |             参数格式错误              |
|  401  |             用户认证失败              |
|  402  | 业务逻辑错误，前端可以弹出后端返回的 msg 作为 toast |
| >=500 |             服务器内部错误             |

## 计费账号：

```shell
    username: kfadmin
    password: kfadmin123
```

## 计费

登录: /api/bill/login

## 客服后台

登录: /api/kf-be/login


## 前台


## 通用接口: 
* 获取验证码id: /api/public/captchaId
* 显示验证码: `<img src="/api/public/$captchaId.png" onclick="reload()">` , `reload` 函数需要设置 `src=/api/public/$captchaId.png?reload=$nowTime`
* 表单提交后，返回了业务逻辑错误，则需要重新获取验证码