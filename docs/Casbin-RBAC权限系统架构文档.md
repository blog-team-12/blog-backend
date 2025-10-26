# Casbin + RBAC 权限系统架构文档

## 1. 系统概述

本项目采用 **Casbin + RBAC（基于角色的访问控制）** 模式实现权限管理系统。Casbin 作为权限执行引擎，RBAC 作为权限模型，两者结合提供了灵活、高效的权限控制方案。

## 2. 核心组件架构

### 2.1 权限模型定义 (model.conf)

```ini
[request_definition]
r = sub, obj

[policy_definition]  
p = sub, obj

[role_definition]
g = _,_

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj
```

**解释：**
- `sub`: 主体（用户ID或角色代码）
- `obj`: 对象（资源，如API路径或菜单代码）
- `g`: 角色继承关系（用户 → 角色）
- `p`: 权限策略（角色 → 资源）

### 2.2 数据库表结构

```
用户表 (users)
├── id (用户ID)
├── username
├── email
└── ... (注意：不包含role_id字段)

角色表 (roles)  
├── id (角色ID)
├── code (角色代码，如 "admin", "user")
├── name (角色名称)
└── ...

菜单表 (menus)
├── id (菜单ID) 
├── code (菜单代码)
├── name (菜单名称)
└── ...

API表 (apis)
├── id (API ID)
├── path (API路径)
├── method (HTTP方法)
└── ...

关系表（核心）：
├── user_roles (用户-角色关系) ← 用户角色通过此表管理
├── role_menus (角色-菜单关系)  
└── menu_apis (菜单-API关系)
```

## 3. 权限控制逻辑关系

### 3.1 三层权限模型

```
用户 (User) → 角色 (Role) → 菜单 (Menu) → API
     ↓           ↓           ↓         ↓
   userID    roleCode    menuCode   path:method
```

**权限传递链：**
1. **用户 → 角色**：用户拥有一个或多个角色
2. **角色 → 菜单**：角色可以访问特定的前端菜单页面
3. **菜单 → API**：菜单关联后端API接口

### 3.2 Casbin 策略存储

Casbin 将权限关系存储为策略规则：

```
# 用户角色关系 (g策略)
g, 1, admin        # 用户1拥有admin角色
g, 2, user         # 用户2拥有user角色

# 角色菜单权限 (p策略)  
p, admin, user_manage, read    # admin角色可以访问用户管理菜单
p, user, dashboard, read       # user角色可以访问仪表板菜单

# 菜单API权限 (p策略)
p, user_manage, /api/users:GET, access     # 用户管理菜单可以调用用户列表API
p, dashboard, /api/stats:GET, access       # 仪表板菜单可以调用统计API
```

## 4. 权限验证流程

### 4.1 API权限验证中间件

```go
// 权限验证流程
func (p *PermissionMiddleware) CheckPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 检查白名单（登录、注册等公开接口）
        if isWhiteList(path, method) {
            c.Next()
            return
        }
        
        // 2. 提取用户信息（从JWT Token）
        userID := extractUserID(c)
        
        // 3. 检查超级管理员（直接放行）
        if isSuperAdmin(userID) {
            c.Next() 
            return
        }
        
        // 4. 验证API权限
        resource := fmt.Sprintf("%s:%s", path, method)
        hasPermission := casbin.Enforce(userID, resource, "access")
        
        if !hasPermission {
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 4.2 权限验证逻辑

```
请求 /api/users (GET) 
    ↓
1. 用户ID: "1"
    ↓  
2. Casbin查询: Enforce("1", "/api/users:GET", "access")
    ↓
3. Casbin内部推理:
   - 用户1有admin角色: g("1", "admin") ✓
   - admin角色有user_manage菜单权限: p("admin", "user_manage", "read") ✓  
   - user_manage菜单有/api/users:GET权限: p("user_manage", "/api/users:GET", "access") ✓
    ↓
4. 返回: true (允许访问)
```

## 5. 用户注册与默认角色分配

### 5.1 注册流程

```go
func (u *UserService) Register(ctx *gin.Context, req *request.RegisterReq) (*entity.User, error) {
    // 1. 验证用户信息（邮箱、用户名唯一性）
    
    // 2. 创建用户记录
    user := &entity.User{
        Username: req.Username,
        Password: util.BcryptHash(req.Password),
        Email:    req.Email,
        UUID:     uuid.Must(uuid.NewV4()),
        // 注意：用户表不包含role_id字段，角色通过关系表管理
    }
    err = u.userRepo.Create(ctx, user)
    
    // 3. 分配默认角色（从配置获取）
    defaultRoleCode := global.Config.System.DefaultRoleCode
    if defaultRoleCode == "" {
        defaultRoleCode = "user" // 兜底默认值
    }
    
    // 4. 查找默认角色
    defaultRole, err := u.roleRepo.GetByCode(ctx, defaultRoleCode)
    
    // 5. 通过权限服务分配角色（同时更新数据库和Casbin）
    err = u.permissionService.AssignRoleToUser(ctx, user.ID, defaultRole.ID)
    
    return user, nil
}
```

### 5.2 角色分配的双重操作

```go
func (p *PermissionService) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
    // 获取角色信息
    role, err := roleRepo.GetByID(ctx, roleID)
    
    // 1. 数据库操作：在user_roles表中插入关系记录
    err = roleRepo.AssignRoleToUser(ctx, userID, roleID)
    
    // 2. Casbin操作：添加用户角色关系到内存策略
    userIDStr := strconv.FormatUint(uint64(userID), 10)
    _, err = p.casbinSvc.Enforcer.AddRoleForUser(userIDStr, role.Code)
    
    // 如果Casbin操作失败，回滚数据库操作（保证数据一致性）
    if err != nil {
        roleRepo.RemoveRoleFromUser(ctx, userID, roleID)
        return err
    }
    
    return nil
}
```

## 6. 权限同步机制

### 6.1 系统启动时的权限同步

```go
func (p *PermissionService) SyncAllPermissionsToCasbin(ctx context.Context) error {
    // 1. 清空现有权限
    p.ClearAllPermission(ctx)
    
    // 2. 同步用户角色关系
    p.SyncUserRolesToCasbin(ctx)
    
    // 3. 同步角色菜单权限  
    p.SyncRoleMenusToCasbin(ctx)
    
    // 4. 同步菜单API权限
    p.SyncMenuAPIsToCasbin(ctx)
    
    return nil
}
```

### 6.2 数据一致性保证

- **双写策略**：每次权限变更同时更新数据库和Casbin
- **事务回滚**：Casbin操作失败时回滚数据库操作
- **定期同步**：系统启动时从数据库重新同步所有权限到Casbin

## 7. 配置化默认角色

### 7.1 配置文件 (configs.yaml)

```yaml
system:
  host: "0.0.0.0"
  port: 8080
  env: "debug"
  default_role_code: "user"        # 默认角色代码
  default_role_name: "普通用户"     # 默认角色名称
```

### 7.2 配置结构体

```go
type System struct {
    Host            string `yaml:"host"`
    Port            int    `yaml:"port"`
    Env             string `yaml:"env"`
    DefaultRoleCode string `yaml:"default_role_code"` // 新用户默认角色
    DefaultRoleName string `yaml:"default_role_name"` // 默认角色显示名
}
```

## 8. 关键特性

### 8.1 优势

1. **灵活性**：支持复杂的权限模型和动态权限变更
2. **性能**：Casbin内存执行，权限验证高效
3. **一致性**：数据库与Casbin双重保证数据一致性
4. **可配置**：默认角色可通过配置文件灵活调整
5. **可扩展**：支持多种权限模型（RBAC、ABAC等）

### 8.2 安全保障

1. **白名单机制**：公开接口无需权限验证
2. **超级管理员**：特殊角色拥有所有权限
3. **权限传递**：通过角色继承实现权限传递
4. **事务保证**：权限变更的原子性操作

## 9. 使用示例

### 9.1 新用户注册

```
1. 用户提交注册信息
2. 系统创建用户记录
3. 从配置获取默认角色 "user"
4. 分配角色：
   - 数据库：INSERT INTO user_roles (user_id, role_id) VALUES (1, 2)
   - Casbin：AddRoleForUser("1", "user")
5. 用户获得默认权限，可以访问普通用户菜单和API
```

### 9.2 权限验证

```
1. 用户访问 /api/users
2. 中间件提取用户ID: "1"  
3. Casbin验证：Enforce("1", "/api/users:GET", "access")
4. 内部推理：用户1 → user角色 → dashboard菜单 → /api/users:GET
5. 返回验证结果
```

## 10. 总结

本权限系统通过 **Casbin + RBAC** 的组合，实现了：

- **分层权限控制**：用户 → 角色 → 菜单 → API 的四层权限模型
- **配置化管理**：默认角色可通过配置文件灵活调整
- **数据一致性**：数据库与Casbin的双重同步机制
- **高性能验证**：内存级别的权限验证执行
- **安全可靠**：完善的错误处理和事务回滚机制

这种架构既保证了权限控制的灵活性和性能，又确保了系统的安全性和可维护性。


**版本**: v0.5
**最后更新**: 2025-10-26 
**维护者**: 王得贤

# 后期需要升级的地方
无法控制前端的界面/按钮展示

不麻烦的！改动其实不大，而且我已经为你分析好了具体的实现方案。

## 📊 改动评估

### 🟢 **改动量：小到中等**
- **新增文件**：1-2个（按钮权限实体、前端权限API控制器）
- **修改文件**：3-4个（菜单实体扩展、权限服务、Repository层）
- **数据库变更**：1个新表 + 1个字段扩展

### 🚀 **实现难度：简单**
因为你的项目架构已经很完善了：
1. **RBAC基础**：Casbin + 用户角色系统已经完整
2. **分层架构**：Repository-Service-Controller 层次清晰
3. **权限框架**：菜单权限查询逻辑已存在，只需扩展

### 💡 **后期实现步骤预览**
```
1️⃣ 扩展Menu实体 (5分钟)
   - 添加 ButtonPermissions 字段

2️⃣ 创建Button实体 (10分钟)  
   - 按钮权限表结构

3️⃣ 新增API接口 (15分钟)
   - GET /api/user/permissions/frontend
   - 返回用户的菜单+按钮权限

4️⃣ Repository层扩展 (10分钟)
   - 查询用户按钮权限方法

5️⃣ 测试验证 (10分钟)
```

### ✨ **优势**
- **渐进式**：不影响现有功能
- **向后兼容**：现有API继续正常工作
- **架构一致**：遵循你现有的设计模式

**总结**：大概1小时左右的工作量，改动不大，而且可以分步骤实现，不会影响现有系统的稳定性。你的架构设计得很好，为后期扩展留了很好的空间！ 🎯
        