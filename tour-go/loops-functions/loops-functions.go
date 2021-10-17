package main

import (
	"fmt"
	"unsafe"
	"math/rand"
	"math"
	"time"
)

func IsNaN(f float64) (is bool) {
	return f != f
}

func IsInf(f float64, sign int) bool {
	return sign >= 0 && f > MaxFloat64 || sign <= 0 && f < -MaxFloat64
}

func NaN() float64 { return Float64frombits(uvnan) }

func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }

const (
	uvnan    = 0x7FF8000000000001
	uvinf    = 0x7FF0000000000000
	uvneginf = 0xFFF0000000000000
	uvone    = 0x3FF0000000000000
	mask     = 0x7FF
	shift    = 64 - 11 - 1
	bias     = 1023
	signMask = 1 << 63
	fracMask = 1<<shift - 1
	MaxFloat64             = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
)

func Float64frombits(b uint64) float64 { return *(*float64)(unsafe.Pointer(&b)) }

func MathSqrt(x float64) float64 {
	// special cases
	switch {
	case x == 0 || IsNaN(x) || IsInf(x, 1):
		return x
	case x < 0:
		return NaN()
	}
	ix := Float64bits(x)
	// normalize x
	exp := int((ix >> shift) & mask)
	if exp == 0 { // subnormal x
		for ix&(1<<shift) == 0 {
			ix <<= 1
			exp--
		}
		exp++
	}
	exp -= bias // unbias exponent
	ix &^= mask << shift
	ix |= 1 << shift
	if exp&1 == 1 { // odd exp, double x to make it even
		ix <<= 1
	}
	exp >>= 1 // exp = exp/2, exponent of square root
	// generate sqrt(x) bit by bit
	ix <<= 1
	var q, s uint64               // q = sqrt(x)
	r := uint64(1 << (shift + 1)) // r = moving bit from MSB to LSB
	for r != 0 {
		t := s + r
		if t <= ix {
			s = t + r
			ix -= t
			q += r
		}
		ix <<= 1
		r >>= 1
	}
	// final rounding
	if ix != 0 { // remainder, result not exact
		q += q & 1 // round according to extra bit
	}
	ix = q>>1 + uint64(exp-1+bias)<<shift // significand + biased exponent)
	return Float64frombits(ix)
}

func mod(x float64) float64 {
	if x < 0{
		return x * -1
	}
	return x
}

func Sqrt(x float64) float64 {
	closeAtLeast := 0.00000001
	z := 1.0
	for mod(z*z-x) >= closeAtLeast {
		z -= (z * z - x) / (2 * z)
	}
	return z
}

const numberOfElements = 10000000

func mysqrt(x []float64) {
	start := time.Now()
	for i := 0; i < numberOfElements; i++ {
		Sqrt(x[i])
	}
	elapsed := time.Since(start)
	fmt.Println("MySqrt time:", elapsed)
}

func softMathSqrt(x []float64) {
	start := time.Now()
	for i := 0; i < numberOfElements; i++ {

		MathSqrt(x[i])
	}
	elapsed := time.Since(start)
	fmt.Println("SoftSqrt time:", elapsed)
}

func archMathSqrt(x []float64) {
	start := time.Now()
	for i := 0; i < numberOfElements; i++ {
		math.Sqrt(x[i])
	}
	elapsed := time.Since(start)
	fmt.Println("math.Sqrt time:", elapsed)
}

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	x := make([]float64, numberOfElements)
	for i:= 0; i < numberOfElements; i++ {
		x[i] = (float64)(rand.Intn(100000))
	}

	mysqrt(x)
	softMathSqrt(x)
	archMathSqrt(x)
	
	xn := 2000000.0
	fmt.Println(Sqrt(xn))
	fmt.Println(MathSqrt(xn))
	fmt.Println(math.Sqrt(xn))
}