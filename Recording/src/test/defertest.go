package main

import "fmt"

type Person struct {
	name string
}

func (p *Person)show() {
	fmt.Println(p.name)
}

func main() {
	p := &Person{"wukai"}
	defer p.show()

	p = &Person{"qq"}


}
