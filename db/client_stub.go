package db

type BQClientStub struct {
	records []*QueryResult
	err     error
}

func NewBQClientStub(results []*QueryResult, err error) BQClientStub {
	return BQClientStub{
		records: results,
		err:     NewQueryError("Failed", err),
	}
}
func (c *BQClientStub) SendQuery(query string) ([]*QueryResult, error) {
	return c.records, c.err
}
