package responseerror

type ResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"statuscode"`
}

func CreateError(err string, code int) ResponseError {
	var responseError = ResponseError{
		Message: err,
		Code:    code,
	}
	return responseError
}
