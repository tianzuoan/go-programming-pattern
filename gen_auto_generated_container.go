// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package main

type Int64Container struct {
	s []int64
}

func NewInt64Container() *Int64Container {
	return &Int64Container{s: []int64{}}
}
func (c *Int64Container) Put(val int64) {
	c.s = append(c.s, val)
}
func (c *Int64Container) Get() int64 {
	r := c.s[0]
	c.s = c.s[1:]
	return r
}