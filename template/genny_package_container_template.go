package PackageName

import "github.com/cheekybits/genny/generic"

type GenericType generic.Type

type GenericNameContainer struct {
	s []GenericType
}

func NewGenericNameContainer() *GenericNameContainer {
	return &GenericNameContainer{s: []GenericType{}}
}
func (c *GenericNameContainer) Put(val GenericType) {
	c.s = append(c.s, val)
}
func (c *GenericNameContainer) Get() GenericType {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}
