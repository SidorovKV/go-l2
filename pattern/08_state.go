package pattern

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

type LightBulbState interface {
	isLighted()
}

type LightBulb struct {
	state LightBulbState
}

type Lighted struct{}

func (l *Lighted) isLighted() {
	fmt.Println("Lighted")
}

type Off struct{}

func (o *Off) isLighted() {
	fmt.Println("Light is off")
}

func (l *LightBulb) On() {
	l.state = &Lighted{}
}

func (l *LightBulb) Off() {
	l.state = &Off{}
}

func StateExample() {
	lb := LightBulb{&Off{}}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		if rnd.Intn(2) == 0 {
			lb.Off()
		} else {
			lb.On()
		}
		lb.state.isLighted()
		time.Sleep(300 * time.Millisecond)
		fmt.Println()
	}
}

/*
Состояние — это поведенческий паттерн проектирования, который позволяет объектам менять поведение в зависимости
от своего состояния. Извне создаётся впечатление, что изменился класс объекта.

Паттерн Состояние предлагает создать отдельные классы для каждого состояния, в котором может пребывать объект,
а затем вынести туда поведения, соответствующие этим состояниям.

Вместо того, чтобы хранить код всех состояний, первоначальный объект, называемый контекстом,
будет содержать ссылку на один из объектов-состояний и делегировать ему работу, зависящую от состояния.

Благодаря тому, что объекты состояний будут иметь общий интерфейс, контекст сможет делегировать работу состоянию, н
е привязываясь к его классу. Поведение контекста можно будет изменить в любой момент,
подключив к нему другой объект-состояние.

Очень важным нюансом, отличающим этот паттерн от Стратегии, является то, что и контекст, и сами конкретные состояния
могут знать друг о друге и инициировать переходы от одного состояния к другому.

Применимость
- Когда у нас есть объект, поведение которого кардинально меняется в зависимости от внутреннего состояния,
причём типов состояний много, и их код часто меняется.

Паттерн предлагает выделить в собственные классы все поля и методы, связанные с определёнными состояниями.
Первоначальный объект будет постоянно ссылаться на один из объектов-состояний, делегируя ему часть своей работы.
Для изменения состояния в контекст достаточно будет подставить другой объект-состояние.

- Когда код класса содержит множество больших, похожих друг на друга, условных операторов,
которые выбирают поведения в зависимости от текущих значений полей класса.

- Паттерн предлагает переместить каждую ветку такого условного оператора в собственный класс.
Тут же можно поселить и все поля, связанные с данным состоянием.

- Когда мы сознательно используем табличную машину состояний, построенную на условных операторах,
но вынуждены мириться с дублированием кода для похожих состояний и переходов.

Паттерн Состояние позволяет реализовать иерархическую машину состояний, базирующуюся на наследовании.
Мы можете отнаследовать похожие состояния от одного родительского класса и вынести туда весь дублирующий код.

Преимущества
- Избавляет от множества больших условных операторов машины состояний.
- Концентрирует в одном месте код, связанный с определённым состоянием.
- Упрощает код контекста.

Недостатки
- Может неоправданно усложнить код, если состояний мало и они редко меняются.
*/
