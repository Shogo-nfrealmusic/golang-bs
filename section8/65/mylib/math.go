/*
mylib is my special library.
*/
package mylib

/*
Average returns the average of a slice of integers.
*/

func Average(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}
	return total / len(s)
}