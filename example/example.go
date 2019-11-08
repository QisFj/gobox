package main

import (
	"fmt"
	"github.com/qisfj/gobox"
)

func main() {
	var (
		t1  T1
		t2  T2
		t1s []T1
		t2s []T2
		err error
	)
	errCheck := func() {
		if err != nil {
			panic(err)
		}
	}
	box := gobox.New()
	for i := 0; i < 5; i++ {
		t1.Value = i + 1
		err = box.Set(&t1)
		errCheck()
		fmt.Println("T1: ", t1)
		t2.T1 = t1
		err = box.Set(&t2)
		errCheck()
		fmt.Println("T2: ", t2)
	}
	err = box.Set(&t1)
	errCheck()
	fmt.Print("T1: ", t1)
	t1.Value++
	err = box.Update(&t1, gobox.FieldFilter("ID", gobox.Eq(6)))
	errCheck()
	t1 = T1{}
	err = box.Get(&t1, gobox.FieldFilter("ID", gobox.Eq(6)))
	errCheck()
	fmt.Println("->", t1)

	err = box.Set(&t2)
	errCheck()
	fmt.Print("T2: ", t2)
	err = box.UpdateAttr(&t2, map[string]interface{}{"T1": t1}, gobox.FieldFilter("ID", gobox.Eq(6)))
	errCheck()
	t2 = T2{}
	err = box.Get(&t2, gobox.FieldFilter("ID", gobox.Eq(6)))
	errCheck()
	fmt.Println("-> ", t2)
	err = box.Get(&t1s, gobox.FieldFilter("ID", gobox.In(1, 3, 5)))
	errCheck()

	fmt.Println("T1s: ", t1s)
	err = box.Get(&t2s, gobox.FieldFilter("ID", gobox.In(2, 4, 6)))
	errCheck()
	fmt.Println("T2s: ", t2s)
}

type T1 struct {
	ID    int
	Value int
}

type T2 struct {
	ID int
	T1 T1
}
