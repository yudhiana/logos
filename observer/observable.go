package observer

func NewObservable() Observable {
	return &ConcreteSubject{
		Observers: make(map[*ConcreteObserver]struct{}),
	}
}

type ConcreteSubject struct {
	Observers map[*ConcreteObserver]struct{}
}

func (c *ConcreteSubject) Register(o Observer) {
	obs, ok := o.(*ConcreteObserver)
	if !ok {
		return
	}
	c.Observers[obs] = struct{}{}
}

func (c *ConcreteSubject) Unregister(o Observer) {
	obs, ok := o.(*ConcreteObserver)
	if !ok {
		return
	}
	delete(c.Observers, obs)
}

func (c *ConcreteSubject) TriggerEvent(event string, data interface{}) {
	for obs := range c.Observers {
		obs.Notify(data)
	}
}
