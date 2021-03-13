package worker

// Worker schedule or job
type Worker interface {
	Start() error
	Stop() error
	Name() string
}
