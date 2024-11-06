package core

type IEquals interface {
	Equals(interface{}) bool
}

type ICopy interface {
	Copy() interface{}
}
