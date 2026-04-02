package errorx

// 错误码规范（6位数字）:
//   1000xx → 通用系统错误
//   1001xx → 用户模块
//   1002xx → 认证模块
//   1003xx → 消息模块
//   1004xx → 群组模块
//   1005xx → 文件/媒体模块

const (
	CodeSuccess      = 0
	CodeServerError  = 100001 // 服务器内部错误
	CodeInvalidParam = 100002 // 请求参数错误
	CodeUnauthorized = 100003 // 未登录或 token 失效
	CodeForbidden    = 100004 // 无权限
	CodeNotFound     = 100005 // 资源不存在
	CodeTooManyReqs  = 100006 // 请求过于频繁

	CodeUserNotFound      = 100101 // 用户不存在
	CodeUserAlreadyExists = 100102 // 用户已存在
	CodePasswordWrong     = 100103 // 密码错误
	CodeUserDisabled      = 100104 // 账户已被禁用

	CodeTokenInvalid  = 100201 // token 无效
	CodeTokenExpired  = 100202 // token 已过期
	CodeTokenMissing  = 100203 // token 缺失
	CodeRefreshFailed = 100204 // 刷新 token 失败

	CodeMsgSendFailed    = 100301 // 消息发送失败
	CodeMsgNotFound      = 100302 // 消息不存在
	CodeMsgTooLong       = 100303 // 消息内容过长
	CodeMsgUnsupportType = 100304 // 不支持的消息类型

	CodeGroupNotFound  = 100401 // 群组不存在
	CodeGroupFull      = 100402 // 群组已满
	CodeNotGroupMember = 100403 // 非群成员
	CodeNotGroupOwner  = 100404 // 非群主，无权操作
	CodeAlreadyInGroup = 100405 // 已在群组中
)

// defaultMessages 错误码对应的默认用户提示
var defaultMessages = map[int]string{
	CodeSuccess:           "ok",
	CodeServerError:       "服务器内部错误",
	CodeInvalidParam:      "请求参数错误",
	CodeUnauthorized:      "请先登录",
	CodeForbidden:         "无权限执行此操作",
	CodeNotFound:          "资源不存在",
	CodeTooManyReqs:       "操作过于频繁，请稍后再试",
	CodeUserNotFound:      "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodePasswordWrong:     "密码错误",
	CodeUserDisabled:      "账户已被禁用",
	CodeTokenInvalid:      "登录状态无效，请重新登录",
	CodeTokenExpired:      "登录已过期，请重新登录",
	CodeTokenMissing:      "请先登录",
	CodeRefreshFailed:     "刷新登录状态失败",
	CodeMsgSendFailed:     "消息发送失败",
	CodeMsgNotFound:       "消息不存在",
	CodeMsgTooLong:        "消息内容超出限制",
	CodeMsgUnsupportType:  "不支持的消息类型",
	CodeGroupNotFound:     "群组不存在",
	CodeGroupFull:         "群组人数已达上限",
	CodeNotGroupMember:    "您不在该群组中",
	CodeNotGroupOwner:     "仅群主可执行此操作",
	CodeAlreadyInGroup:    "您已在该群组中",
}

// Message 根据错误码获取默认提示
func Message(code int) string {
	if msg, ok := defaultMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
