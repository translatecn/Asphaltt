package hqu

import "testing"

func TestFibN(t *testing.T) {
	testFibn(0, 1, t)
	testFibn(1, 1, t)
	testFibn(2, 1, t)
	testFibn(3, 2, t)
	testFibn(4, 3, t)
	testFibn(5, 5, t)
	testFibn(6, 8, t)
	testFibn(7, 13, t)
	testFibn(8, 21, t)
	testFibn(9, 34, t)
	testFibn(10, 55, t)
	testFibn(11, 89, t)
	testFibn(12, 144, t)
	testFibn(13, 233, t)
}

func testFibn(n, expected int, t *testing.T) {
	testFibnWithAlgo(n, expected, "iteration", t)
	testFibnWithAlgo(n, expected, "recursion", t)
	testFibnWithAlgo(n, expected, "polynomial", t)
}

func testFibnWithAlgo(n, expected int, algo string, t *testing.T) {
	var res int
	if algo == "iteration" {
		res = FiboN(n)
	} else if algo == "recursion" {
		res = FiboNRecursive(n)
	} else if algo == "polynomial" {
		res = FiboNPolynomial(n)
	}
	if res != expected {
		t.Logf("the %dth fibonacci number, algo:%s, result:%d, expected:%d", n, algo, res, expected)
		t.Fail()
	}
}

func TestPow(t *testing.T) {
	testPow(2.33, 29.47295521, 4, t)
}

func testPow(x, expected float64, n int, t *testing.T) {
	res := pow(float64(x), n)
	if res-expected >= 0.0001 {
		t.Logf("%f^%d, result:%f, expected:%f", x, n, res, expected)
		t.Fail()
	}
}
