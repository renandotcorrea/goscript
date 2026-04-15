// Package try provides error handling utilities that panic on errors.
// These are useful for scripts and command-line tools where error handling can be simplified.
package try

// Try panics if the provided error is not nil.
func Try(err error) {
	if err != nil {
		panic(err)
	}
}

// Try1 returns the provided value if err is nil, otherwise it panics with the error
func Try1[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

// Try2 returns the provided values if err is nil, otherwise it panics with the error.
// Example usage:
//
//	value1, value2 := try.Try2(func() (T1, T2, error) {
//		// Your code that returns T1, T2 and an error
//	}())
func Try2[T1 any, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return value1, value2
}

// Try3 returns the provided values if err is nil, otherwise it panics with the error
// Example usage:
//
//	value1, value2, value3 := try.Try3(func() (T1, T2, T3, error) {
//		// Your code that returns T1, T2, T3 and an error
//	}())
func Try3[T1 any, T2 any, T3 any](value1 T1, value2 T2, value3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}

	return value1, value2, value3
}

// Try4 returns the provided values if err is nil, otherwise it panics with the error.
// Example usage:
//
//	value1, value2, value3, value4 := try.Try4(func() (T1, T2, T3, T4, error) {
//		// Your code that returns T1, T2, T3, T4 and an error
//	}())
func Try4[T1 any, T2 any, T3 any, T4 any](value1 T1, value2 T2, value3 T3, value4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}

	return value1, value2, value3, value4
}

// Handle executes the provided handler function if a panic occurs during the execution of a Try block.
// It should be deferred at the beginning of a function to catch any panics that occur within the Try block.
// Example usage:
//
//	defer try.Handle(func(err error) {
//		println("Error:", err.Error())
//	})
//
//	try.Try(http.Get("http://httpbin.org/get").JSON(&dest))
func Handle(handler func(err error)) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			handler(err)
		} else {
			panic(r)
		}
	}
}
