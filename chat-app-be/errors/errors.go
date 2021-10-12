package errors

type UnauthorizedError struct {
	Err error
}

func (e *UnauthorizedError) Error() string {
	if e.Err == nil {
		return "Unauthorized"
	}

	return e.Err.Error()
}

type NotFountError struct {
	Err error
}

func (e *NotFountError) Error() string {
	if e.Err == nil {
		return "Not found"
	}

	return e.Err.Error()
}

type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	if e.Err == nil {
		return "Bad Request"
	}

	return e.Err.Error()
}
