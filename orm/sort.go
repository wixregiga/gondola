package orm

import (
	"gondola/orm/driver"
)

type Sort int

const (
	// ASC indicates that the results of the given query should be
	// returned in ascending order for the given field.
	ASC = Sort(driver.ASC)
	// DESC indicates that the results of the given query should be
	// returned in descending order for the given field.
	DESC = Sort(driver.DESC)
)
