package rv

import (
	"fmt"
	"github.com/jfcg/bob"
	"testing"
)

const (
	N   = 99
	tol = 1e-9
)

func Test1(t *testing.T) {
	var b bob.Bob
	r := Rvg{&b}

	m, v := .0, .0
	for i := N; i > 0; i-- {
		d := r.Uni()
		m += d
		v += d * d
	}
	m /= N
	v = (v - N*m*m) / (N - 1)
	fmt.Println(m, v)
}

func Test2(t *testing.T) {
	x := "şevkşıv"
	a, c := FrequencyS(x)
	fmt.Println("alphabet:", string(a))
	e, s := Entropy(c)
	fmt.Println(c, e, s)
}
