package example

type Case1 struct {
	Array []string
}
type Case2Struct struct {
	Name   string
	Value  int
	Value2 int
}
type Case2 struct {
	ArrayStruct []Case2Struct
}

type Config struct {
	Case1 *Case1
	Case2 *Case2
}
