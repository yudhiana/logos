package observer

type ObserverAction func(data interface{})

type Observer interface {
	Notify(data interface{})
}
