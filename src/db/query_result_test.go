package db

import (
	"fmt"
)

func ExampleQueryResult_String() {
	sampleQueryResult := &QueryResult{
		Service: "Cloud SQL", Monthly: 1000.0, Yesterday: 400.0,
	}
	fmt.Println(sampleQueryResult.String())
	// Output: {Service: Cloud SQL, Monthly: 1000.000000, Yesterday: 400.000000}
}
