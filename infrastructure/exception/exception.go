package exception

import "net/http"

type Exception struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []Causes `json:"causes,omitempty"`
}

// user sent a wrong password: there will be a password field and the reason
type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// podemos usar bibliotecas que podemos acabar utilizando que recebe error por parâmetro
// para não criar uma variável e retornar, criamos um método chamado error.
// quando precisar passar error por parâmetro, passaremos esse método.
func (e *Exception) Error() string {
	return e.Message
}

// constructor
func NewException(message string, code int, err string, causes []Causes) *Exception {
	return &Exception{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

func BadRequestException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func UnauthorizedRequestException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

func BadRequestValidationException(message string, causes []Causes) *Exception {
	return &Exception{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func InternalServerException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func NotFoundException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}

func ForbiddenException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

func TooManyRequestsException(message string) *Exception {
	return &Exception{
		Message: message,
		Err:     "too_many_requests",
		Code:    http.StatusTooManyRequests,
	}
}
