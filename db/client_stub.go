package db

import "github.com/tatamiya/gcp-cost-notification/utils"

type BQClientStub struct {
	records []*QueryResult
	err     *utils.CustomError
}

func NewBQClientStub(results []*QueryResult, err error) BQClientStub {
	var queryError *utils.CustomError
	if err == nil {
		queryError = nil
	} else {
		queryError = NewQueryError("Failed", err)
	}
	return BQClientStub{
		records: results,
		err:     queryError,
	}
}
func (c *BQClientStub) SendQuery(query string) ([]*QueryResult, *utils.CustomError) {
	return c.records, c.err
}
