package example

type Case1 struct {
	Array []string
	Name  int
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
	Test1 []string
	//Field1 struct {
	//	TEST string
	//}
	//Case1 *Case1
	//Case2 *Case2
}
