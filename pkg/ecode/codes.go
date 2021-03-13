package ecode

func newErrCode(reason, toast string, detail interface{}, code int, message string) *StatusError {
	if detail == nil {
		detail = map[string]interface{}{}
	}
	return &StatusError{
		Code:      code,
		Message:   message,
		ErrorInfo: &ErrorInfo{Reason: reason, Toast: toast, Detail: detail},
	}
}

// // OK 正常返回业务错误信息
// // HTTP Mapping: 200 业务错误信息
// func OK(reason, toast string, detail interface{}) *StatusError {
// 	return newErrCode(reason, toast, detail, 0, "Ok")
// }

// Cancelled The operation was cancelled, typically by the caller.
// HTTP Mapping: 499 Client Closed Request
func Cancelled(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 1, "Canceled")
}

// Unknown errs.
// HTTP Mapping: 500 Internal Server Error
func Unknown(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 2, "Unknown")
}

// InvalidArgument The client specified an invalid argument.
// HTTP Mapping: 400 Bad Request
func InvalidArgument(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 3, "InvalidArgument")
}

// DeadlineExceeded The deadline expired before the operation could complete.
// HTTP Mapping: 504 Gateway Timeout
func DeadlineExceeded(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 4, "DeadlineExceeded")
}

// NotFound Some requested entity (e.g., file or directory) was not found.
// HTTP Mapping: 404 Not Found
func NotFound(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 5, "NotFound")
}

// AlreadyExists The entity that a client attempted to create (e.g., file or directory) already exists.
// HTTP Mapping: 409 Conflict
func AlreadyExists(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 6, "AlreadyExists")
}

// PermissionDenied The caller does not have permission to execute the specified operation.
// HTTP Mapping: 403 Forbidden
func PermissionDenied(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 7, "PermissionDenied")
}

// ResourceExhausted Some resource has been exhausted, perhaps a per-user quota, or
// perhaps the entire file system is out of space.
// HTTP Mapping: 429 Too Many Requests
func ResourceExhausted(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 8, "ResourceExhausted")
}

// FailedPrecondition The operation was rejected because the system is not in a state
// required for the operation's execution.
// HTTP Mapping: 400 Bad Request
func FailedPrecondition(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 9, "FailedPrecondition")
}

// Aborted The operation was aborted, typically due to a concurrency issue such as
// a sequencer check failure or transaction abort.
// HTTP Mapping: 409 Conflict
func Aborted(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 10, "Aborted")
}

// OutOfRange The operation was attempted past the valid range.  E.g., seeking or
// reading past end-of-file.
// HTTP Mapping: 400 Bad Request
func OutOfRange(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 11, "OutOfRange")
}

// Unimplemented The operation is not implemented or is not supported/enabled in this service.
// HTTP Mapping: 501 Not Implemented
func Unimplemented(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 12, "Unimplemented")
}

// Internal This means that some invariants expected by the
// underlying system have been broken.  This errs code is reserved
// for serious errs.
//
// HTTP Mapping: 500 Internal Server Error
func Internal(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 13, "Internal")
}

// Unavailable The service is currently unavailable.
// HTTP Mapping: 503 Service Unavailable
func Unavailable(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 14, "Unavailable")
}

// DataLoss Unrecoverable data loss or corruption.
// HTTP Mapping: 500 Internal Server Error
func DataLoss(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 15, "DataLoss")
}

// Unauthorized The request does not have valid authentication credentials for the operation.
// HTTP Mapping: 401 Unauthorized
func Unauthorized(reason, toast string, detail interface{}) *StatusError {
	return newErrCode(reason, toast, detail, 16, "Unauthenticated")
}

// gRPC2HTTPCode http code
var gRPC2HTTPCode = map[int]int{
	0:  200,
	1:  499,
	2:  500,
	3:  400,
	4:  504,
	5:  404,
	6:  409,
	7:  403,
	8:  429,
	9:  400,
	10: 409,
	11: 400,
	12: 501,
	13: 500,
	14: 503,
	15: 500,
	16: 401,
}
