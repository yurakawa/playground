package main

import "fmt"

type ServiceImpl struct{}
type AnotherServiceImpl struct{}

func (s *ServiceImpl) Apply(id int) error {
	fmt.Println("ServiceImpl:", id)
	return nil
}
func (s *AnotherServiceImpl) DoSomething(id int) error {
	fmt.Println("AnotherServiceImpl:", id)
	return nil
}

// 上位階層が定義する抽象
type OrderService interface {
	Apply(id int) error
}
type AnotherService interface {
	DoSomething(id int) error
}

// 上位階層の利用者側の型
type Application struct {
	os OrderService
	as AnotherService
}

// 他言語のコンストラクタインジェクションに相当する実装
func NewApplication(os OrderService, as AnotherService) *Application {
	return &Application{os: os, as: as}
}

func (app *Application) Apply(id int) error {
	return app.os.Apply(id)
}
func (app *Application) DoSomething(id int) error {
	return app.as.DoSomething(id)
}

func main() {
	app := NewApplication(&ServiceImpl{}, &AnotherServiceImpl{})
	app.Apply(19)
	app.DoSomething(20)
}

/* Setterを用意しておいて DIする
func (app *Application) SetService(os OrderService) {
	app.os = os
}
func main() {
	app := &Application{}
	app.SetService(&ServiceImpl{})
	app.Apply(19)
}
*/

/*
func(app *Application) Apply(os OrderService, id int) error {
	return os.Apply(id)
}
func main() {
	app := &Application{}
	app.Apply(&ServiceImpl{}, 19)
}
*/
