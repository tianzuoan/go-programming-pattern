package main

import (
	"fmt"
)

func goGenerationDemo() {
	generateByShellScript()
}
func generateByShellScript() {
	generateUint32Example()
	generateStringExample()
}

//go:generate genny -in=./template/genny_package_template.go -out=gen_auto_generated.go gen "KeyType=string,int ValueType=string,int PackageName=main"
//go:generate genny -in=./template/genny_package_container_template.go -out=gen_auto_generated_container.go gen "GenericType=int64 PackageName=main GenericName=Int64"
func generateByThirdPackage() {

}

//go:generate bash ./gen.sh ./template/container.tmp.template main uint32 container
func generateUint32Example() {
	var u uint32 = 42
	c := NewUint32Container()
	c.Put(u)
	v := c.Get()
	fmt.Printf("generateExample: %d (%T)\n", v, v)
}

//go:generate bash ./gen.sh ./template/container.tmp.template main string container
func generateStringExample() {
	var s string = "Hello"
	c := NewStringContainer()
	c.Put(s)
	v := c.Get()
	fmt.Printf("generateExample: %s (%T)\n", v, v)
}
