package request

import "fmt"

type ContentTypeHeader struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

func (bh *ContentTypeHeader) Validate() error {
	if bh.Value != "application/json" {
		return fmt.Errorf("invalid Content-Type header: %s", bh.Value)
	}

	return nil
}
