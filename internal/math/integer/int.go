/*
Package math provides basic math functions for int (+int 32, uints) so we don't
have to do type casting everywhere.
*/

package integer

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
