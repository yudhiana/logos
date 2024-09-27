package ward

func NewObservable() *ConcreteSubject {
	return &ConcreteSubject{
		Observers: make(map[*ConcreteObserver]struct{}),
	}
}

type ConcreteSubject struct {
	Observers map[*ConcreteObserver]struct{}
	Message   interface{}
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

func (c *ConcreteSubject) Notify() {
	for obs := range c.Observers {
		obs.Update(c.Message)
	}
}

func (c *ConcreteSubject) SetData(data interface{}) {
	c.Message = data
	c.Notify()
}
