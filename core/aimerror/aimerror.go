package aimerror

import "strings"

// 错误包含所有发生的错误
type Errors []error

// GetErrors获取已经发生的所有错误并返回一个错误片段(Error type)
func (errs Errors) GetErrors() []error {
	return errs
}

//	检查是否是错误
func (errs Errors) IsError() bool {
	return len(errs) > 0
}

//	Add为给定的错误片添加一个错误
func (errs Errors) Add(newErrors ...error) Errors {
	for _, err := range newErrors {
		if err == nil {
			continue
		}

		if errors, ok := err.(Errors); ok {
			errs = errs.Add(errors...)
		} else {
			ok = true
			for _, e := range errs {
				if err == e {
					ok = false
				}
			}
			if ok {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

//	Error获取已发生的所有错误的片段，并将其作为格式化字符串返回
func (errs Errors) Error() string {
	var errors = []string{}
	for _, e := range errs {
		errors = append(errors, e.Error())
	}
	return strings.Join(errors, "\r\n")
}
