// Package slice provides a generic Slice[T] type with functional programming utilities.
// It supports operations like Filter, Map, Reduce, and other sequence transformations.
package slice

import (
	"reflect"
)

// Slice is a generic slice type that wraps a Go slice and provides functional methods.
type Slice[T any] []T

// Contains checks if the slice contains the specified value.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	fmt.Println(slice.Contains(3)) // Output: true
//	fmt.Println(slice.Contains(6)) // Output: false
func (s Slice[T]) Contains(value T) bool {
	for _, v := range s {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}

	return false
}

// Filter returns a new slice containing only the elements that satisfy the provided predicate function.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	evenNumbers := slice.Filter(func(x int) bool {
//		return x%2 == 0
//	})
//	fmt.Println(evenNumbers) // Output: [2, 4]
func (s Slice[T]) Filter(predicate func(T) bool) Slice[T] {
	return filter(s, predicate)
}

func filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// First returns the first element of the slice, or nil if the slice is empty.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	fmt.Println(slice.First()) // Output: 1
//
//	emptySlice := Slice[int]{}
//	fmt.Println(emptySlice.First()) // Output: nil
func (s Slice[T]) First() *T {
	if s.IsEmpty() {
		return nil
	}
	return &s[0]
}

// Last returns the last element of the slice, or nil if the slice is empty.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	fmt.Println(slice.Last()) // Output: 5
//
//	emptySlice := Slice[int]{}
//	fmt.Println(emptySlice.Last()) // Output: nil
func (s Slice[T]) Last() *T {
	if s.IsEmpty() {
		return nil
	}
	return &s[s.Len()-1]
}

// IsEmpty returns true if the slice is empty, false otherwise.
func (s Slice[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Map returns a new slice containing the results of applying the provided transform function to each element of the original slice.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	mapped := slice.Map(func(x int) int {
//		return x * 2
//	})
//	fmt.Println(mapped) // Output: [2, 4, 6, 8, 10]
func (s Slice[T]) Map(transform func(T) T) Slice[T] {
	result := make(Slice[T], len(s))
	for i, v := range s {
		result[i] = transform(v)
	}
	return result
}

// ForEach executes the provided action function for each element in the slice.
func (s Slice[T]) ForEach(action func(T)) {
	for _, v := range s {
		action(v)
	}
}

// Reduce reduces the slice to a single value by applying the provided reducer function to each element of the slice, starting with the initial value.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	sum := slice.Reduce(func(acc, x int) int {
//		return acc + x
//	}, 0)
//	fmt.Println(sum) // Output: 15
func (s Slice[T]) Reduce(reducer func(T, T) T, initial T) T {
	return reduce(s, reducer, initial)
}

func reduce[T any](s []T, reducer func(T, T) T, initial T) T {
	result := initial
	for _, v := range s {
		result = reducer(result, v)
	}
	return result
}

// Reverse returns a new slice with the elements in reverse order.
// Example usage:
//
//	slice := Slice[int]{1, 2, 3, 4, 5}
//	reversed := slice.Reverse()
//	fmt.Println(reversed) // Output: [5, 4, 3, 2, 1]
func (s Slice[T]) Reverse() Slice[T] {
	return reverse(s)
}

func reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Unique returns a new slice containing only the unique elements from the original slice.
// Example usage:
//
//	slice := Slice[int]{1, 2, 2, 3, 4, 4, 5}
//	unique := slice.Unique()
//	fmt.Println(unique) // Output: [1, 2, 3, 4, 5]
func (s Slice[T]) Unique() Slice[T] {
	return unique(s)
}

func unique[T any](s []T) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		exists := false
		for _, current := range result {
			if reflect.DeepEqual(current, v) {
				exists = true
				break
			}
		}

		if !exists {
			result = append(result, v)
		}
	}

	return result
}

// Len returns the length of the slice.
func (s Slice[T]) Len() int {
	return len(s)
}

// Cap returns the capacity of the slice.
func (s Slice[T]) Cap() int {
	return cap(s)
}

// Append appends the specified values to the slice and returns the resulting slice.
func (s Slice[T]) Append(values ...T) Slice[T] {
	return append(s, values...)
}

// Chunk splits the current slice into chunks of the provided size.
// If n is less than or equal to zero, it returns an empty []Slice[T].
func (s Slice[T]) Chunk(n int) []Slice[T] {
	if n <= 0 {
		return []Slice[T]{}
	}

	result := make([]Slice[T], 0, (s.Len()+n-1)/n)
	for i := 0; i < s.Len(); i += n {
		end := i + n
		if end > s.Len() {
			end = s.Len()
		}

		result = append(result, Slice[T](s[i:end]))
	}

	return result
}

// FlatMap maps each input item to a slice and flattens all results in order.
func FlatMap[T any, U any](s []T, fn func(T) []U) Slice[U] {
	result := make(Slice[U], 0)
	for _, v := range s {
		result = append(result, fn(v)...)
	}

	return result
}

// ToMap converts a slice of values into a map using the provided key selector function.
// Example usage:
//
//	type User struct {
//		ID   int
//		Name string
//	}
//
//	users := Slice[User]{
//		{ID: 1, Name: "Alice"},
//		{ID: 2, Name: "Bob"},
//	}
//
//	userMap := users.ToMap(func(u User) int {
//		return u.ID
//	})
//
//	fmt.Println(userMap) // Output: map[1:{1 Alice} 2:{2 Bob}]
func ToMap[K comparable, V any](s []V, keySelector func(V) K) map[K]V {
	result := make(map[K]V)
	for _, v := range s {
		key := keySelector(v)
		result[key] = v
	}

	return result
}
