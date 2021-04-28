package driver

import (
	"gondola/orm/index"
	"gondola/orm/query"
	"reflect"
)

type JoinType int

const (
	InnerJoin JoinType = iota
	OuterJoin
	LeftJoin
	RightJoin
)

type Join interface {
	Model() Model
	Type() JoinType
	Query() query.Q
}

type Model interface {
	Type() reflect.Type
	Table() string
	Fields() *Fields
	Indexes() []*index.Index
	Map(qname string) (string, reflect.Type, error)
	Skip() bool
	Join() Join
}
