package object

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	BOOLEAN_OBJ            = "BOOLEAN"
	NULL_OBJ               = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
