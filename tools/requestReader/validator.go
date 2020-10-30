package requestReader

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func (rr *RequestReader) Read(request interface{}) error {
	if err := rr.cntx.Bind(request); err != nil {
		return err
	}

	if err := rr.validator.Struct(request); err != nil {
		return err
	}
	return nil
}
