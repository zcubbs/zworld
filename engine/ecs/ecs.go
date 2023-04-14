package ecs

import "reflect"

type Id uint32

type Component interface {
	ComponentSet(interface{})
}

type BasicStorage struct {
	list map[Id]interface{}
}

func NewBasicStorage() *BasicStorage {
	return &BasicStorage{
		list: make(map[Id]interface{}),
	}
}

func (s *BasicStorage) Read(id Id) (interface{}, bool) {
	val, ok := s.list[id]
	return val, ok
}

func (s *BasicStorage) Write(id Id, val interface{}) {
	s.list[id] = val
}

type Engine struct {
	reg       map[string]*BasicStorage
	idCounter Id
}

func NewEngine() *Engine {
	return &Engine{
		reg:       make(map[string]*BasicStorage),
		idCounter: 0,
	}
}

func (e *Engine) NewId() Id {
	id := e.idCounter
	e.idCounter++
	return id
}

func name(t interface{}) string {
	name := reflect.TypeOf(t).String()
	if name[0] == '*' {
		return name[1:]
	}

	return name
}

func GetStorage(e *Engine, t interface{}) *BasicStorage {
	name := name(t)
	storage, ok := e.reg[name]
	if !ok {
		e.reg[name] = NewBasicStorage()
		storage, _ = e.reg[name]
	}
	return storage
}

func Read(e *Engine, id Id, val Component) bool {
	storage := GetStorage(e, val)
	newVal, ok := storage.Read(id)
	if ok {
		val.ComponentSet(newVal)
	}

	return ok
}

func Write(e *Engine, id Id, val interface{}) {
	s := GetStorage(e, val)
	s.Write(id, val)
}

func Each(e *Engine, t interface{}, f func(id Id, a interface{})) {
	s := GetStorage(e, t)
	for id, a := range s.list {
		f(id, a)
	}
}
