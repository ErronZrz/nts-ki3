package addr

type Generator interface {
	TotalNum() int
	HasNext() bool
	NextHost() string
}
