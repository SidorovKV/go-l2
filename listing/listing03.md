Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
<nil>
false

Хоть значение у переменной err и nil, но запись об os.PathError(который реализует интерфейс error)
попала в itables. Таким образом переменная по типу != nil.
Пустой интерфейс не создаёт запись в itables.
```
