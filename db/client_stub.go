package db

type BQClientStub struct {
	records []*QueryResult
	err     error
}

func NewBQClientStub(results []*QueryResult, err error) BQClientStub {
	var queryError error
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
func (c *BQClientStub) SendQuery(query string) ([]*QueryResult, error) {
	return c.records, c.err
}
