package shapes

import "math"

type Rectangle struct {
	W, H float64
}

type Circle struct {
	R float64
}

// по трём сторонам
type Triangle struct {
	A, B, C float64
}

type Ring struct {
	Outer, Inner float64
}

// прямоугольник
func (s Rectangle) Perimeter() float64 {
	return (s.H + s.W) * 2
}

func (s Rectangle) Area() float64 {
	return s.H * s.W
}

// круг
func (s Circle) Perimeter() float64 {
	return 2 * math.Pi * s.R
}

func (s Circle) Area() float64 {
	return math.Pi * s.R * s.R
}

// треугольник
func (s Triangle) Perimeter() float64 {
	return s.A + s.B + s.C
}

func (s Triangle) Area() float64 {
	p := s.Perimeter() / 2
	result := p * (p - s.A) * (p - s.B) * (p - s.C)
	if result > 0 {
		return math.Sqrt(result)
	} else {
		return 0
	}

}

//кольцо
func (s Ring) Perimeter() float64 {
	return 2 * math.Pi * (s.Outer + s.Inner)
}

func (s Ring) Area() float64 {
	return math.Pi*s.Outer*s.Outer - math.Pi*s.Inner*s.Inner
}

func (s *Rectangle) Scale(k float64) {
	s.H *= k
	s.W *= k
}
