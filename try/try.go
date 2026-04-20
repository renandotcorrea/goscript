// Package try provides error handling utilities that panic on errors.
// These are useful for scripts and command-line tools where error handling can be simplified.
package try

// Try panics with the provided error if it is not nil.
func Try(err error) {
	if err != nil {
		panic(err)
	}
}

// Try1 returns the provided value if err is nil, otherwise it panics with err.
func Try1[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

// Try2 returns the provided values if err is nil, otherwise it panics with err.
// Example usage:
//
//	value1, value2 := try.Try2(someFunc())
func Try2[T1 any, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return value1, value2
}

// Try3 returns the provided values if err is nil, otherwise it panics with err.
// Example usage:
//
//	value1, value2, value3 := try.Try3(someFunc())
func Try3[T1 any, T2 any, T3 any](value1 T1, value2 T2, value3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}

	return value1, value2, value3
}

// Try4 returns the provided values if err is nil, otherwise it panics with err.
// Example usage:
//
//	value1, value2, value3, value4 := try.Try4(someFunc())
func Try4[T1 any, T2 any, T3 any, T4 any](value1 T1, value2 T2, value3 T3, value4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}

	return value1, value2, value3, value4
}

// Handle catches panics from Try functions and passes them to the provided handler.
// Non-error panics are re-panicked.
// It should be deferred at the beginning of a function to handle any panics from Try calls.
// Example usage:
//
//	defer try.Handle(func(err error) {
//		fmt.Println("Error:", err)
//	})
//	try.Try(someFunc())
func Handle(handler func(err error)) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			handler(err)
		} else {
			panic(r)
		}
	}
}
