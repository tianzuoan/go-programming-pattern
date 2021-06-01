package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// check is implement the relative interface
var _ Shape = (*Circle)(nil)
var _ Shape = (*Rectangle)(nil)
var _ ColorfulVisitor = JsonVisitor
var _ ColorfulVisitor = XmlVisitor

func k8sVisitorDemo() {
	//visitorDemo()
	//infoVisitorDemo()
	infoVisitorWithDecorateDemo()
}

func visitorDemo() {
	shapes := []Shape{Circle{6}, Rectangle{3, 7}}
	//c := Circle{10}
	//r := Rectangle{100, 200}
	//shapes := []Shape{c, r}
	for _, s := range shapes {
		s.Accept(JsonVisitor)
		s.Accept(XmlVisitor)
	}
}

func infoVisitorDemo() {
	info := Info{}
	var v Visitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	//Visitor 们一层套一层；
	//我用 loadFile 假装从文件中读取数据；
	//最后执行 v.Visit(loadfile) ，这样，我们上面的代码就全部开始激活工作了。
	v.Visit(loadFile)

	//上面的代码有以下几种功效：
	//解耦了数据和程序；
	//使用了修饰器模式；
	//还做出了 Pipeline 的模式。
}

func infoVisitorWithDecorateDemo() {
	info := Info{}
	var v Visitor = &info

	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}

	nameVisitorFunc := func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	}

	OtherThingsVisitorFunc := func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	}

	//v = LogVisitor{v}
	//v = NameVisitor{v}
	//v = OtherThingsVisitor{v}

	//这种对于装饰方式方式是便利顺序调用
	//如果变量v采用嵌套赋值（上面注释打开），那么装饰方法nameVisitorFunc，OtherThingsVisitorFunc就没有必要了，采用的将是嵌套调用
	v = NewDecoratedVisitor(v, nameVisitorFunc, OtherThingsVisitorFunc)

	err := v.Visit(loadFile)
	if err != nil {
		panic(err)
	}
}

type ColorfulVisitor func(shape Shape)

type Shape interface {
	Accept(colorfulVisitor ColorfulVisitor)
}

type Circle struct {
	Radius int64
}

type Rectangle struct {
	Weight, Height int64
}

func (c Circle) Accept(colorfulVisitor ColorfulVisitor) {
	colorfulVisitor(c)
}

func (r Rectangle) Accept(colorfulVisitor ColorfulVisitor) {
	colorfulVisitor(r)
}

func JsonVisitor(shape Shape) {
	contentBytes, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println("json marsha content:", string(contentBytes))
}

func XmlVisitor(shape Shape) {
	contentBytes, err := xml.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println("xml marsha content:", string(contentBytes))
}

type VisitorFunc func(*Info, error) error

type Visitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

type LogVisitor struct {
	visitor Visitor
}

type NameVisitor struct {
	visitor Visitor
}

type OtherThingsVisitor struct {
	visitor Visitor
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

type DecoratedVisitor struct {
	visitor    Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}
	return DecoratedVisitor{v, fn}
}

// Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		if err := fn(info, nil); err != nil {
			return err
		}
		for _, visitFunc := range v.decorators {
			if err := visitFunc(info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}
