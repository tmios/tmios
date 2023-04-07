package errors

import (
	"tmios/lib/errors"
)

func New(msg string) error {
	return errors.InternalNew(msg, msg)
}

var (
	ErrParam                 = errors.BadRequest(400000, "参数错误:")
	ErrNoTokenFound          = errors.BadRequest(400010, "No token found in headers")
	ErrInvalidAppToken       = errors.BadRequest(400011, "Invalid app token")
	ErrNoAuth                = errors.BadRequest(400020, "没有登录")
	ErrLocked                = errors.BadRequest(400021, "账号已锁定")
	ErrInvalidUserOrPassword = errors.BadRequest(400040, "账号或密码错误")
	ErrInvalidLoginType      = errors.BadRequest(400050, "Invalid login type")
	ErrInvalidOriPassword    = errors.BadRequest(400060, "原始密码错误")
	ErrParseFormFile         = errors.BadRequest(400100, "Parse FormFile failed")
	ErrInvalidRequest        = errors.BadRequest(400101, "非法请求")
	ErrInvalidResponse       = errors.BadRequest(400103, "非法返回")

	ErrNotFound       = errors.Conflict(400404, "记录不存在:")
	ErrNoPermission   = errors.Conflict(409010, "没有权限")
	ErrInvalidRole    = errors.Conflict(409011, "错误的角色")
	ErrDuplicateEntry = errors.Conflict(400102, "记录重复:")

	ErrUserExisted      = errors.Conflict(409110, "该用户已存在")
	ErrUserDisabled     = errors.Conflict(409112, "用户被禁用")
	ErrDeleteRole       = errors.Conflict(409160, "删除角色失败:")
	ErrAddRoleRight     = errors.Conflict(409170, "添加权限失败:")
	ErrRightNotInRole   = errors.Conflict(409171, "角色范围不包含该权限:")
	ErrRoleScope        = errors.Conflict(409172, "RoleScope错误:")
	ErrDeleteDepartment = errors.Conflict(409300, "删除部门失败:")

	ErrClientConn      = errors.Conflict(410100, "客户端连接错误:")
	ErrClientBrowerDir = errors.Conflict(410110, "客户端浏览目录错误:")
	ErrClientAuth      = errors.Conflict(410120, "客户端未授权:")

	ErrDBCurd = errors.Conflict(410210, "数据库错误:")
)
