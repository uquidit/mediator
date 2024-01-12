package securechangeapi

import (
	"fmt"
)

type my_trigger struct {
	id        int
	name      string
	path      string
	arguments string
}

func (mt my_trigger) GetLabel() string {
	return fmt.Sprintf("%s: mediator '%s', arguments '%s' (#%d)", mt.name, mt.path, mt.arguments, mt.id)
}
func (mt my_trigger) GetValue() int {
	return mt.id
}
