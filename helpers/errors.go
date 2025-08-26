package helpers

type CustomError struct {
	Message string
	Code    int
}

func (e *CustomError) Error() string {
	if e != nil {
		return e.Message
	} else {
		return "Unknow error"
	}
}

func System(msg string) *CustomError {
	return &CustomError{Message: msg, Code: 500}
}

func NotFound(msg string) *CustomError {
	return &CustomError{Message: msg, Code: 404}
}

func BadRequest(msg string) *CustomError {
	return &CustomError{Message: msg, Code: 400}
}

func Unauthorized(msg string) *CustomError {
	return &CustomError{Message: msg, Code: 401}
}

func NoType(err error) *CustomError {
	if err == nil {
		return nil
	}

	if customErr, ok := err.(*CustomError); ok {
		return customErr
	}

	return &CustomError{Message: err.Error(), Code: 500}
}
