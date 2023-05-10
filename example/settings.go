package example

type Case1 struct {
	Array []string
	Name  int
}

type Config struct {
	TestArray []string
	Case1     *Case1
	Case3     []Case1
}
