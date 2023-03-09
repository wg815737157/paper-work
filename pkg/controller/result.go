package controller

type (
	BaseResult struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func NewSuccessResult(data interface{}) BaseResult {
	return NewResult(0, "ok", data)
}

func NewSuccessMsgResult(msg string, data interface{}) BaseResult {
	return NewResult(0, msg, data)
}

func NewResult(code int, message string, data interface{}) BaseResult {
	return BaseResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewFailedResult(code int, msg string) BaseResult {
	return BaseResult{
		Code:    code,
		Message: msg,
	}
}
