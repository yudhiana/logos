package ward

import "github.com/yudhiana99/ward/observer"

func Recover() {
	if r := recover(); r != nil {
		observer.NewObserver("panic", observer.Panic).Notify(r)
	}
}
