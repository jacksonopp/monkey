package object

type ReturnValue struct {
	Value Object
}

func (r ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
