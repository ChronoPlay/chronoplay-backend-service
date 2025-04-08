package helpers

type CustomEror struct {
	Code    uint32
	Message string
}

func (c *CustomEror) Custom(code uint32, message string) *CustomEror {
	return &CustomEror{
		Code:    code,
		Message: message,
	}
}

func (c *CustomEror) BadRequest(message string) *CustomEror {
	return &CustomEror{
		Code:    400,
		Message: message,
	}
}

func (c *CustomEror) System(message string) *CustomEror {
	return &CustomEror{
		Code:    500,
		Message: message,
	}
}
