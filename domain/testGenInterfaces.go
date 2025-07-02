package domain

type Method struct {
	MethodName string
	Results    []string
	Params     []string
	CountIfs   []string
}

type Interface struct {
	InterfaceName string
}
