package helpers

type CustomEror struct {
	Message string
	Code    int
}

func (e *CustomEror) Error() string {
	return e.Message
}

func System(msg string) *CustomEror {
	return &CustomEror{Message: msg, Code: 500}
}

func NotFound(msg string) *CustomEror {
	return &CustomEror{Message: msg, Code: 404}
}

func BadRequest(msg string) *CustomEror {
	return &CustomEror{Message: msg, Code: 400}
}
