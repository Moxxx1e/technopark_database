package requestReader

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	. "github.com/technopark_database/internal/consts"
	"github.com/technopark_database/internal/helpers/errors"
)

type RequestReader struct {
	cntx      echo.Context
	validator *validator.Validate
}

func NewRequestReader(cntx echo.Context) *RequestReader {
	return &RequestReader{
		cntx:      cntx,
		validator: validator.New(),
	}
}

func (rr *RequestReader) Read(request interface{}) *errors.Error {
	if err := rr.cntx.Bind(request); err != nil {
		return errors.New(CodeInternalServerError, err)
	}

	//if err := rr.validator.Struct(request); err != nil {
	//	return errors.Get(CodeBadRequest)
	//}
	return nil
}
