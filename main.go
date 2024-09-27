package main

import ward "github.com/yudhiana99/ward/observer"

func main() {

	// Create subject
	Observable := ward.NewObservable()

	// Create observers
	observer := ward.NewObserver()

	observer2 := ward.NewObserver()

	// Register observers
	Observable.Register(observer)
	Observable.Register(observer2)

	// Set data
	Observable.SetData("next processing data")
}
