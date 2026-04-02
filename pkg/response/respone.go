package response

import (
	"easy-im/pkg/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 统一响应体结构
// 所有接口返回格式：{"code":0,"msg":"ok","data":{}}
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// OK 成功响应，携带数据
func OK(w http.ResponseWriter, data any) {
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: errorx.CodeSuccess,
		Msg:  "ok",
		Data: data,
	})
}

// OKWithMsg 成功响应，自定义提示（如"创建成功"）
func OKWithMsg(w http.ResponseWriter, msg string, data any) {
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: errorx.CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Fail 失败响应，传入 AppError
// handler 层统一调用此方法，禁止直接写 http.Error
func Fail(w http.ResponseWriter, err error) {
	if appErr, ok := errorx.IsAppError(err); ok {
		httpx.WriteJson(w, httpStatusFromCode(appErr.Code), &Body{
			Code: appErr.Code,
			Msg:  appErr.Msg,
		})
		return
	}
	// 非业务错误，统一返回服务器错误，不暴露内部细节
	httpx.WriteJson(w, http.StatusInternalServerError, &Body{
		Code: errorx.CodeServerError,
		Msg:  errorx.Message(errorx.CodeServerError),
	})
}

// ParamError 参数校验失败的快捷响应
func ParamError(w http.ResponseWriter, msg string) {
	httpx.WriteJson(w, http.StatusBadRequest, &Body{
		Code: errorx.CodeInvalidParam,
		Msg:  msg,
	})
}

// httpStatusFromCode 将业务错误码映射到 HTTP 状态码
// 注意：业务 code 与 HTTP status 是两个维度，不要混淆
func httpStatusFromCode(code int) int {
	switch code {
	case errorx.CodeUnauthorized, errorx.CodeTokenInvalid,
		errorx.CodeTokenExpired, errorx.CodeTokenMissing:
		return http.StatusUnauthorized // 401
	case errorx.CodeForbidden, errorx.CodeNotGroupOwner:
		return http.StatusForbidden // 403
	case errorx.CodeNotFound, errorx.CodeUserNotFound,
		errorx.CodeMsgNotFound, errorx.CodeGroupNotFound:
		return http.StatusNotFound // 404
	case errorx.CodeTooManyReqs:
		return http.StatusTooManyRequests // 429
	case errorx.CodeInvalidParam:
		return http.StatusBadRequest // 400
	default:
		return http.StatusOK // 业务错误统一用 200，靠 code 区分
	}
}
