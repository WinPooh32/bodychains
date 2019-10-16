package nbody

import (
	"math"

	"github.com/robaho/fixed"
)

// Number is decimal alias for fixed.Fixed
type Number = fixed.Fixed

// G is Gravitational constant
var (
	//G = fixed.NewI(6674, -8) // 6674E-8
	G = fixed.NewI(6674, 0)

	Half = fixed.NewI(5, 1) // 5E-1
	Cube = fixed.NewI(3, 0) // 3
	Eps  = fixed.NewI(1, 1) // 1E-4
)

// Vector represents (x,y,z) 3d vector
type Vector struct {
	X, Y, Z Number
}

// Body is simple particle of Universe
type Body struct {
	Coord    Vector
	Velocity Vector
	Mass     Number
}

func dist(l, r Body) Number {
	// tx = l.x - r.x
	tx := l.Coord.X.Sub(r.Coord.X)
	ty := l.Coord.Y.Sub(r.Coord.Y)
	tz := l.Coord.Z.Sub(r.Coord.Z)

	// tx = tx*tx
	tx = tx.Mul(tx)
	ty = ty.Mul(ty)
	tz = tz.Mul(tz)

	return fixed.NewF(math.Sqrt(tx.Add(ty).Add(tz).Float()))
}

func force(l, r Body) Vector {
	var f Vector
	d := dist(l, r)

	if d.Cmp(Eps) == 1 && d.Int() < 1000 {
		//module = G * ((l.mass * r.mass) / d^3)

		mass := l.Mass.Mul(r.Mass).Div(d)   // = (l.mass * r.mass)/d
		module := mass.Div(d.Mul(d)).Mul(G) // = (mass / d^2) * G

		// f.x = (l.x - r.x) * module
		f.X = l.Coord.X.Sub(r.Coord.X).Mul(module)
		f.Y = l.Coord.Y.Sub(r.Coord.Y).Mul(module)
		f.Z = l.Coord.Z.Sub(r.Coord.Z).Mul(module)
	}

	return f
}
