package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	String() string
}

type Circle struct {
	radius float64
}
func (c *Circle) Area() float64{
	return math.Pi*(c.radius)*(c.radius)
}
func (c *Circle) String() string{
	return fmt.Sprintf("Circle(r = %.2f)", c.radius)
}

func main(){
	var s1 Shape = &Circle{radius: 5}
	fmt.Printf("%s\n",s1.String())
	fmt.Printf("Area: %.2f\n",s1.Area())
}