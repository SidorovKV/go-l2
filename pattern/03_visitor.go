package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

type Visitor interface {
	VisitA(elementA ElementA) string
	VisitB(elementB ElementB) string
}

type Element interface {
	Accept(visitor Visitor)
}

type ElementA struct {
	info string
}

func (a *ElementA) Accept(v Visitor) {
	v.VisitA(*a)
}

type ElementB struct {
	info string
}

func (b *ElementB) Accept(v Visitor) {
	v.VisitB(*b)
}

type ConcreteVisitor struct{}

func (c *ConcreteVisitor) VisitA(a ElementA) string {
	return a.info
}

func (c *ConcreteVisitor) VisitB(b ElementB) string {
	return b.info
}

/*
Посетитель — это поведенческий паттерн проектирования, который позволяет добавлять в программу новые операции,
не изменяя классы объектов, над которыми эти операции могут выполняться.

Паттерн Посетитель предлагает разместить новое поведение в отдельном классе,
вместо того чтобы множить его сразу в нескольких классах. Объекты, с которыми должно было быть связано поведение,
не будут выполнять его самостоятельно. Вместо этого вы будете передавать эти объекты в методы посетителя.

Код поведения, скорее всего, должен отличаться для объектов разных классов,
поэтому и методов у посетителя должно быть несколько. Названия и принцип действия этих методов будет схож,
но основное отличие будет в типе принимаемого в параметрах объекта.

Применимость
- Когда нам нужно выполнить какую-то операцию над всеми элементами сложной структуры объектов.

Посетитель позволяет применять одну и ту же операцию к объектам различных классов.

- Когда над объектами сложной структуры объектов надо выполнять некоторые не связанные между собой операции,
но мы не хотим «засорять» классы такими операциями.

Посетитель позволяет извлечь родственные операции из классов, составляющих структуру объектов,
поместив их в один класс-посетитель. Если структура объектов является общей для нескольких приложений,
то паттерн позволит в каждое приложение включить только нужные операции.

- Когда новое поведение имеет смысл только для некоторых классов из существующей иерархии.

Посетитель позволяет определить поведение только для этих классов, оставив его пустым для всех остальных.

Преимущества
- Упрощает добавление операций, работающих со сложными структурами объектов.
- Объединяет родственные операции в одном классе.
- Посетитель может накапливать состояние при обходе структуры элементов.

Недостатки
- Паттерн не оправдан, если иерархия элементов часто меняется.
- Может привести к нарушению инкапсуляции элементов.
*/
