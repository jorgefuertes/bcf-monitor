package monitor

type Monitorizable interface {
	Check() error
	IsUp() bool
	Down() bool
	Up() bool
}