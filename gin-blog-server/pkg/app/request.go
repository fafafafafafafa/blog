package app

import (
	"go-gin-example/pkg/logging"

	"github.com/astaxie/beego/validation"
)

func MarkError(errors []*validation.Error) {

	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
