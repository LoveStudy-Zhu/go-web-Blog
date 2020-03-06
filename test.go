package main

import "fmt"

type Student struct {
	a int64
	b int64
}

func main()  {
	student := make([] *Student,2)
	var stu *Student
	stu=&Student{1,2}
	student[0] =  stu
	fmt.Println(student[0].a)
}