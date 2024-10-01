package observer

func Recover() {
	if r := recover(); r != nil {
		NewObserver("panic", Panic).Notify(r)
	}
}
