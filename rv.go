package rv

import (
	"math"
	"unsafe"
)

//	Bob (probably a little sponge) represents any source that is Squeez()ed to output
//	"independent & uniformly distributed" uint64 values. Bob also Absorb()s uint64
//	values to possibly change its internal state.
type Bob interface {
	Absorb(uint64)
	Squeez() uint64
}

//	Random variable generator is a Bob
type Rvg struct {
	Bob
}

func (r Rvg) u12() float64 {
	i := r.Squeez() &^ (3 << 62) // set sign=exponent=0 so it'll be a double from [1,2)
	i |= (1<<10 - 1) << 52
	return *(*float64)(unsafe.Pointer(&i))
}

//	Returns a uniformly distributed float from (0,1)
func (r Rvg) Uni() float64 {
	x := .0
	for x == 0 {
		x = r.u12() - 1
	}
	return x
}

//	Returns a uniformly distributed float from [0,1)
func (r Rvg) Uni0() float64 {
	return r.u12() - 1
}

//	Returns a uniformly distributed float from (0,1]
func (r Rvg) Uni1() float64 {
	return 2 - r.u12()
}

//	Returns a uniformly distributed float from (-1,1)
func (r Rvg) Uni2() float64 {
	i := r.Squeez() &^ (1 << 62) // set exponent=0 so it'll be a double from +/-[1,2)
	i |= (1<<10 - 1) << 52

	f := *(*float64)(unsafe.Pointer(&i))
	if f > 0 {
		return f - 1
	}
	return f + 1
}

//	Returns two independent & normally distributed floats with zero mean and unit variance
func (r Rvg) Gauss() (float64, float64) {
	a, b, m := .0, .0, .0
	for m == 0 || m >= 1 {
		a = r.Uni2()
		b = r.Uni2()
		m = a*a + b*b
	}

	m = math.Sqrt(-2 * math.Log(m) / m)
	return m * a, m * b
}

//	Returns an exponentially distributed float with unit mean
func (r Rvg) Exp1() float64 {
	return -math.Log(r.Uni1())
}

//	Convert map to list
func map2ls(mp map[rune]uint32) ([]rune, []uint32) {
	i, al, cl := 0, make([]rune, len(mp)), make([]uint32, len(mp))
	for r, c := range mp {
		al[i] = r
		cl[i] = c
		i++
	}
	return al, cl
}

//	Returns alphabet & counts given a list of runes
func Frequency(ls []rune) (al []rune, cl []uint32) {
	mp := make(map[rune]uint32)
	for _, r := range ls {
		mp[r]++ // collect number of occurances
	}
	return map2ls(mp)
}

//	Returns unicode alphabet & counts given a utf8 string
func FrequencyS(s string) (al []rune, cl []uint32) {
	mp := make(map[rune]uint32)
	for _, r := range s {
		mp[r]++ // collect number of occurances
	}
	return map2ls(mp)
}

//	Returns entropy in bits & counts sum given a list of counts
func Entropy(cl []uint32) (float64, uint64) {
	a := uint64(0)
	for _, c := range cl {
		a += uint64(c)
	}

	e, s := .0, float64(a) // entropy, sum
	for _, c := range cl {
		if c > 0 {
			p := float64(c) / s
			e -= p * math.Log(p)
		}
	}
	return e * math.Log2E, a
}
