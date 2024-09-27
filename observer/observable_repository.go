package observer

type Observable interface {
	Register(Observer)
	Unregister(Observer)
	TriggerEvent(string, interface{})
}
