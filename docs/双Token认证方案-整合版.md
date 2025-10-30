# 双Token认证方案（整合版）

## 概述
- 本项目采用 Access Token + Refresh Token 的双Token认证机制。经过安全复盘与改造，现方案为：Refresh 通过 HttpOnly Cookie 传递，Access 仍通过请求头传递。
- 目标：在最小改动下，降低 XSS 风险，控制 CSRF 风险，简化前端刷新逻辑，并保持与多点登录/黑名单机制的兼容。

## 内容总览
- 旧的双Token实现方法（Header双Token）
- 风险分析：XSS 与 CSRF
- 方案演进与取舍
- 现方案：Refresh 用 HttpOnly Cookie + Access 用 Header（本项目落地）
- 前端使用示例与注意事项
- 兼容性与迁移影响
- 参考与致谢

---

## 旧的双Token实现方法（已废弃）

### Token类型
- Access Token：短期有效，用于 API 访问验证。
- Refresh Token：长期有效，用于刷新 Access Token。

### 传输方式（旧约定）
```http
x-access-token: <ACCESS_TOKEN>
x-refresh-token: <REFRESH_TOKEN>
```

### 核心流程（旧）
1. 登录：后端签发双token并放入响应体；前端保存两者。
2. API 调用：前端在请求头携带 `x-access-token`。
3. Token 刷新：`access_token` 过期时，使用 `x-refresh-token` 刷新获取新的双token。
4. 多点登录控制：新登录将旧 Refresh Token 加入黑名单。

### 优点（旧）
- 实现简单，前端易调试；无需考虑跨域 Cookie 配置。

### 缺点（旧）
- XSS 风险：若在可读环境（localStorage、JS变量）持有 Refresh Token，易被脚本窃取。
- 前端复杂：需管理两类 Token 的存取、同步与重试。
- CSRF 风险虽可控，但策略分散、实现复杂。

---

## 风险分析

### 什么是 XSS
- 跨站脚本攻击。攻击者向页面注入恶意脚本，在受害者浏览器执行以窃取会话、Cookie、个人信息等。
- 常见类型：
  - 存储型：恶意脚本存入服务端（DB等），他人访问即触发。
  - 反射型：通过 URL 等传入，服务器回显时触发。

### 什么是 CSRF
- 跨站请求伪造。攻击者诱导用户在登录态下访问恶意页面，由浏览器自动携带 Cookie 向目标站点发起伪造请求，目标站错误认为是用户自愿行为。

### CSRF 与 XSS 的区别（要点）
- 目标：CSRF 利用用户登录态伪造请求；XSS 窃取数据/会话。
- 依赖：CSRF 依赖用户已登录目标站；XSS 依赖站点存在输入输出漏洞。
- 防护：CSRF 关注来源校验（SameSite、Origin/Referer、CSRF Token 等）；XSS 关注输入输出过滤与内容安全策略（CSP）。

---

## 方案演进与取舍
1. 旧思路：Access 与 Refresh 均在请求头中传递（实现简单，但 Refresh Token 暴露面大，受 XSS 影响）。
2. 尝试思路：Access 与 Refresh 均放 Cookie 且 `HttpOnly=true`（可减少 XSS 暴露，但 API 调用层面引入更强 CSRF 风险控制需求）。
3. 现方案（折中）：
   - Refresh 用 HttpOnly Cookie（降低 XSS 面）；
   - Access 用 Header（便于前端显式控制与避免非幂等接口被被动伪造）；
   - 配合 SameSite/Origin 校验等手段控制 CSRF 风险。

---

## 现方案：Refresh 用 HttpOnly Cookie + Access 用 Header（本项目落地）

### 客户端约定
- 请求头携带：`x-access-token: <ACCESS_TOKEN>`。
- 刷新时：浏览器自动携带 `x-refresh-token` Cookie（HttpOnly）。

### 后端行为（已实现）
- 登录成功由服务端写入 Cookie：
  - Cookie 名：`x-refresh-token`；属性：`HttpOnly`、`Path=/`、`Secure` 随是否 HTTPS 设置；`Max-Age` 与后端刷新令牌过期一致；建议配置 `SameSite`。
- 刷新接口：仅从 Cookie 读取，不再接受 `x-refresh-token` 请求头。
- 多点登录/黑名单：沿用原有 Redis/黑名单策略，新登录将旧 Refresh Token 加入黑名单。

### 实施细节（代码改动纪要）
- 刷新凭证只读 Cookie：
  - `pkg/jwt/claims.go`：`GetRefreshToken` 仅从 Cookie 读取 `x-refresh-token`。
- 登录响应移除刷新字段：
  - `internal/model/dto/response/LoginResponse` 移除 `refresh_token`、`refresh_token_expires_at` 字段。
- Service 返回签名调整：
  - `internal/service/system/jwtSvc.go`：`IssueLoginTokens` 改为返回 `(LoginResponse, refreshToken, refreshExpiresAt, error)`，不再把 Refresh 写入响应体。
- Controller 设置 Cookie：
  - `internal/controller/system/userCtrl.go`：`TokenNext` 收到 `refreshToken` 与 `refreshExpiresAt` 后，设置 `x-refresh-token` HttpOnly Cookie，并确保响应体仅包含用户信息与 `access_token` 相关字段。
- 中间件/访问令牌校验：沿用原有逻辑，不受影响。

### 优点（现）
- 显著降低 XSS 面：Refresh Token 不在前端可读环境中出现。
- CSRF 可控：结合 SameSite、Origin 校验、HTTPS 与严格路由策略。
- 前端更简：无需存/传 Refresh Token，仅管理 Access Token 与自动刷新。

### 注意点（现）
- 跨域需开启凭证：前端 `withCredentials: true`；后端 CORS 需 `Access-Control-Allow-Credentials: true` 且明确 `Access-Control-Allow-Origin`。
- Cookie `SameSite`：同站推荐 `Lax`；跨站推荐 `None` 并强制 `Secure`。

---

## 前端使用示例

### 基础调用（Axios）
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.API_BASE,
  withCredentials: true, // 确保跨域时携带 Cookie
});

// 请求拦截：附加 Access Token（示例：从内存/安全容器中读取）
api.interceptors.request.use((config) => {
  const accessToken = window.__ACCESS_TOKEN__;
  if (accessToken) {
    config.headers['x-access-token'] = accessToken;
  }
  return config;
});

// 响应拦截：遇 401 时尝试刷新
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    if (err.response?.status === 401) {
      try {
        // 刷新接口仅依赖 Cookie，无需显式传 Refresh Token
        const r = await api.get('/refreshToken');
        if (r.data?.data?.access_token) {
          window.__ACCESS_TOKEN__ = r.data.data.access_token;
          // 复制原请求并重试
          const cfg = err.config;
          cfg.headers['x-access-token'] = window.__ACCESS_TOKEN__;
          return api.request(cfg);
        }
      } catch (_) {}
    }
    return Promise.reject(err);
  }
);

export default api;
```

---

## 兼容性与迁移影响
- 破坏性变化：
  - 登录响应体不再包含 `refresh_token` 及其过期时间。
  - 客户端不得再发送 `x-refresh-token` 请求头。
- 不受影响：
  - 访问令牌校验、中间件拦截；
  - 多点登录、黑名单、Redis 存储策略。

### 迁移建议
- 前端移除对于 `refresh_token` 的存储与传输逻辑。
- 确认跨域时已启用 `withCredentials: true` 并配置后端 CORS。
- 在生产环境启用 HTTPS，并将 Cookie 设为 `Secure`。
- 视业务选择在刷新成功时重写 Refresh Cookie（延长会话）或启用 Refresh 轮换（Rotation）。

---

## FAQ（常见问题）
- Q：接口路径、返回结构是否变化？
  - A：刷新接口路径与返回结构保持不变，仅从 Cookie 读取刷新凭证。
- Q：是否支持登出清除 Cookie？
  - A：可通过后端提供登出接口在响应中清除 `x-refresh-token` Cookie（推荐）。
- Q：如何进一步降低 CSRF 风险？
  - A：结合 `SameSite`、Origin/Referer 校验、仅允许安全方法自动携带 Cookie、对高危操作引入一次性 CSRF Token 等。

---

## 参考与致谢
- 项目原有文档与实现：
  - 《双Token认证方案》（旧）
  - 《双Token方案升级与对比（简版）》
- 相关背景与阐述参考：
  - 双Token的风险与改进思路（博文，CC BY-SA 4.0）：https://blog.csdn.net/2302_80067378/article/details/154023425

> 注：本文为项目内整合版说明，结合现有代码落地情况进行整理，便于团队统一理解与使用。
 
**版本**: v0.5
**最后更新**: 2025-10-28
**维护者**: 王得贤