package pattern

import "fmt"

type Handler interface {
	Handle(order Order)
	SetNext(handler Handler)
}

type Order struct {
	volume    int
	weight    int
	flammable bool
	stinky    bool
}

type VolumeHandler struct {
	maxVolume int
	next      Handler
}

func (vh *VolumeHandler) Handle(order Order) {
	if order.volume <= vh.maxVolume {
		vh.next.Handle(order)
	} else {
		fmt.Println("Deny delivery")
	}
}

func (vh *VolumeHandler) SetNext(handler Handler) {
	vh.next = handler
}

type WeightHandler struct {
	maxWeight int
	next      Handler
}

func (wh *WeightHandler) Handle(order Order) {
	if order.weight <= wh.maxWeight {
		wh.next.Handle(order)
	} else {
		fmt.Println("Deny delivery")
	}
}

func (wh *WeightHandler) SetNext(handler Handler) {
	wh.next = handler
}

type FlammableHandler struct {
	isFlammableAllowed bool
	next               Handler
}

func (fh *FlammableHandler) Handle(order Order) {
	if order.flammable {
		if fh.isFlammableAllowed {
			fh.next.Handle(order)
		} else {
			fmt.Println("Deny delivery")
		}
	} else {
		fh.next.Handle(order)
	}
}

func (fh *FlammableHandler) SetNext(handler Handler) {
	fh.next = handler
}

type StinkyHandler struct {
	isStinkyAllowed bool
	next            Handler
}

func (sh *StinkyHandler) Handle(order Order) {
	if order.stinky {
		if sh.isStinkyAllowed {
			fmt.Println("To delivery")
		} else {
			fmt.Println("Deny delivery")
		}
	} else {
		fmt.Println("To delivery")
	}
}

func (sh *StinkyHandler) SetNext(handler Handler) {
	sh.next = handler
}

func ChainOfResp() {
	volumeHandler := &VolumeHandler{maxVolume: 10}
	weightHandler := &WeightHandler{maxWeight: 10}
	flammableHandler := &FlammableHandler{isFlammableAllowed: true}
	stinkyHandler := &StinkyHandler{isStinkyAllowed: false}

	volumeHandler.SetNext(weightHandler)
	weightHandler.SetNext(flammableHandler)
	flammableHandler.SetNext(stinkyHandler)

	volumeHandler.Handle(Order{volume: 5, weight: 5, flammable: true, stinky: false})
	volumeHandler.Handle(Order{volume: 15, weight: 5, flammable: true, stinky: false})
	volumeHandler.Handle(Order{volume: 5, weight: 15, flammable: true, stinky: true})
	volumeHandler.Handle(Order{volume: 5, weight: 5, flammable: false, stinky: true})
	volumeHandler.Handle(Order{volume: 5, weight: 5, flammable: true, stinky: false})
}

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Цепочка обязанностей — это поведенческий паттерн проектирования, который позволяет передавать запросы
последовательно по цепочке обработчиков. Каждый последующий обработчик решает, может ли он обработать
запрос сам и стоит ли передавать запрос дальше по цепи.

Применимость
- Когда программа должна обрабатывать разнообразные запросы несколькими способами, но заранее неизвестно,
какие конкретно запросы будут приходить и какие обработчики для них понадобятся.

С помощью Цепочки обязанностей мы можем связать потенциальных обработчиков в одну цепь и при получении запроса
поочерёдно спрашивать каждого из них, не хочет ли он обработать запрос.

- Когда важно, чтобы обработчики выполнялись один за другим в строгом порядке.

Цепочка обязанностей позволяет запускать обработчиков последовательно один за другим в том порядке,
в котором они находятся в цепочке.

- Когда набор объектов, способных обработать запрос, должен задаваться динамически.

В любой момент мы можем вмешаться в существующую цепочку и переназначить связи так, чтобы убрать или добавить новое звено.

Преимущества
- Уменьшает зависимость между клиентом и обработчиками.
- Реализует принцип единственной обязанности.
- Реализует принцип открытости/закрытости.

Недостатки
- Запрос может остаться никем не обработанным.
*/
