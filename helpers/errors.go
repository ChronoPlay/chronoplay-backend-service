package helpers

type CustomEror struct {
	Message string
	Code    int
}

func (e *CustomEror) Error() string {
	if e != nil {
		return e.Message
	} else {
		return "Unknow error"
	}
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
