package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

type Builder interface {
	BuildPartA(int) Builder
	BuildPartB(bool) Builder
	BuildPartC(string) Builder
	Build() Product
}

type Product struct {
	PartA int
	PartB bool
	PartC string
}

type ConcreteBuilder struct {
	partA int
	partB bool
	partC string
}

func (cb *ConcreteBuilder) BuildPartA(a int) Builder {
	cb.partA = a
	return cb
}

func (cb *ConcreteBuilder) BuildPartB(b bool) Builder {
	cb.partB = b
	return cb
}

func (cb *ConcreteBuilder) BuildPartC(c string) Builder {
	cb.partC = c
	return cb
}

func (cb *ConcreteBuilder) Build() Product {
	return Product{
		PartA: cb.partA,
		PartB: cb.partB,
		PartC: cb.partC,
	}
}

/*
Строитель — это порождающий паттерн проектирования, который позволяет создавать сложные объекты пошагово.
Строитель даёт возможность использовать один и тот же код строительства для получения разных представлений объектов.

Паттерн Строитель предлагает вынести конструирование объекта за пределы его собственного класса,
поручив это дело отдельным объектам, называемым строителями.

Паттерн предлагает разбить процесс конструирования объекта на отдельные шаги.
Чтобы создать объект, нам нужно поочерёдно вызывать методы строителя. Причём не нужно запускать все шаги,
а только те, что нужны для производства объекта определённой конфигурации.

Применимость
- Когда мы хотим избавиться от «телескопического конструктора».

Допустим, у вас есть один конструктор с десятью опциональными параметрами.
Его неудобно вызывать, поэтому мы можем создать ещё десять конструкторов с меньшим количеством параметров.
Но всё, что они делают — это переадресуют вызов к базовому конструктору, подавая какие-то значения
по умолчанию в параметры, которые пропущены в них самих.
Паттерн Строитель позволяет собирать объекты пошагово, вызывая только те шаги, которые нам нужны.
А значит, больше не нужно пытаться «запихнуть» в конструктор все возможные опции продукта.

- Когда нам нужно собирать сложные составные объекты
Строитель конструирует объекты пошагово, а не за один проход. Более того, шаги строительства можно выполнять рекурсивно.
Заметим, что Строитель не позволяет посторонним объектам иметь доступ к конструируемому объекту,
пока тот не будет полностью готов. Это предохраняет клиентский код от получения незаконченных «битых» объектов.

Преимущества
- Позволяет создавать продукты пошагово.
- Позволяет использовать один и тот же код для создания различных продуктов.
- Изолирует сложный код сборки продукта от его основной бизнес-логики.

Недостатки
Усложняет код программы из-за введения дополнительных классов.

*/
