package ward

type Observable interface {
	Register(Observer)
	Unregister(Observer)
	Notify()
}
