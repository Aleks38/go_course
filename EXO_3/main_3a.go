package main

import (
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceTo(other Point) float64 {
	return math.Sqrt(
		math.Pow(other.X-p.X, 2)+
		math.Pow(other.Y-p.Y, 2),
	)
}

type Rectangle struct {
	MIN Point
	MAX Point
}
func (r Rectangle) Width() float64 {
	return r.MAX.X - r.MIN.X
}

func (r Rectangle) Height() float64 {
	return r.MAX.Y - r.MIN.Y
}

func (r Rectangle) Area() float64 {
	return r.Width() * r.Height()
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width() + r.Height())
}

func (r *Rectangle) Move(dx, dy float64) {
	r.MIN = Point{r.MIN.X + dx, r.MIN.Y + dy}
	r.MAX = Point{r.MAX.X + dx, r.MAX.Y + dy}
}

func (r Rectangle) Describe() string {
	return fmt.Sprintf("Rectangle: MIN=(%.2f, %.2f), MAX=(%.2f, %.2f), Width=%.2f, Height=%.2f, Area=%.2f, Perimeter=%.2f",
		r.MIN.X, r.MIN.Y, r.MAX.X, r.MAX.Y, r.Width(), r.Height(), r.Area(), r.Perimeter())
}

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

func (c Circle) Describe() string {
	return fmt.Sprintf("Cercle: Centre=(%.2f, %.2f), Rayon=%.2f, Aire=%.2f, Circonference=%.2f",
		c.Center.X, c.Center.Y, c.Radius, c.Area(), c.Circumference())
}

func main() {
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 4, Y: 6}

	dist := p1.DistanceTo(p2)
	fmt.Printf("Distance entre p1(%v, %v) et p2(%v, %v) = %.2f\n", p1.X, p1.Y, p2.X, p2.Y, dist)

	rect := Rectangle{
		MIN: Point{X: 0, Y: 0},
		MAX: Point{X: 5, Y: 3},
	}

	fmt.Printf("Rectangle MIN=%v MAX=%v\n", rect.MIN, rect.MAX)
	fmt.Printf("Largeur=%.2f Hauteur=%.2f Surface=%.2f Perimetre=%.2f\n",
		rect.Width(), rect.Height(), rect.Area(), rect.Perimeter())

	rect.Move(2, 3)

	fmt.Printf("Après Move(2,3) -> MIN=%v MAX=%v\n", rect.MIN, rect.MAX)

	circ := Circle{
		Center: Point{X: 2, Y: 3},
		Radius: 5,
	}

	fmt.Printf("Area of circle: %.2f\n", circ.Area())
	fmt.Printf("Circle: %.2f\n", circ.Circumference())

	fmt.Printf("Scaling circle by factor of 2...\n")
	circ.Scale(2)
	fmt.Printf("New area of circle: %.2f\n", circ.Area())
	fmt.Printf("New circumference of circle: %.2f\n", circ.Circumference())
}
