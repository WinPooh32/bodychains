package nbody

import (
	"crypto/sha256"

	"github.com/robaho/fixed"
)

type Universe []Body

func univForce(b Body, univ Universe) Vector {
	var sum Vector

	for i := range univ {
		f := force(b, univ[i])

		sum.X = sum.X.Sub(f.X)
		sum.Y = sum.Y.Sub(f.Y)
		sum.Z = sum.Z.Sub(f.Z)
	}

	return sum
}

func StepVelocity(dt Number, subUniv, fullUniv Universe) {
	size := uint64(len(subUniv))
	for i := uint64(0); i < size; i++ {
		// velocity = velocity + ( force / mass ) * dt;

		b := &(subUniv)[i]

		if b.Mass.IsZero() {
			continue
		}

		dm := dt.Div(b.Mass) // = dt / mass
		uforce := univForce(*b, fullUniv)

		// = uforce * dm
		dx := uforce.X.Mul(dm)
		dy := uforce.Y.Mul(dm)
		dz := uforce.Z.Mul(dm)

		// b.velocity.x += dx
		b.Velocity.X = b.Velocity.X.Add(dx)
		b.Velocity.Y = b.Velocity.Y.Add(dy)
		b.Velocity.Z = b.Velocity.Z.Add(dz)
	}
}

func ApplyVelocity(dt Number, univ Universe) {
	size := uint64(len(univ))
	for i := uint64(0); i < size; i++ {
		b := &(univ)[i]

		// position = position + velocity * dt;
		b.Coord.X = b.Coord.X.Add(b.Velocity.X.Mul(dt))
		b.Coord.Y = b.Coord.Y.Add(b.Velocity.Y.Mul(dt))
	}
}

func HashUniverse(univ Universe, index int64) []byte {
	avg := stateAvg(univ)

	indexBinary, _ := (fixed.NewI(index, 0)).MarshalBinary()

	binary, _ := avg.MarshalBinary()
	binary = append(binary, indexBinary...)

	hash := sha256.Sum256(binary)
	return hash[:]
}

func stateAvg(univ Universe) Number {
	total := len(univ) * 3 // len * fields_count
	avg := Number{}
	invLen := fixed.NewF(1.0 / float64(total))

	for i := range univ {
		x := univ[i].Coord.X.Mul(invLen)
		y := univ[i].Coord.Y.Mul(invLen)
		z := univ[i].Coord.Z.Mul(invLen)
		avg = avg.Add(x).Add(y).Add(z)
	}

	return avg
}

func startContiniousSimulation(ticks uint64, dt Number, univ Universe) {
	// for t := uint64(0); t < ticks; t++ {
	// 	begin := time.Now()

	// 	StepVelocity(dt, univ)
	// 	ApplyVelocity(dt, univ)

	// 	avg := stateAvg(univ).Add(fixed.NewI(int64(t), 0))
	// 	binary, _ := avg.MarshalBinary()
	// 	hash := sha256.Sum256(binary)

	// 	fmt.Printf("%15fms, hash: %x\n",
	// 		float64(time.Now().Sub(begin).Nanoseconds())*float64(0.000001),
	// 		hash[:],
	// 	)
	// }
}
