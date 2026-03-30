package slice

import (
	"slices"
	"testing"
)

func TestSliceContains(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}

	if !s.Contains(3) {
		t.Fatal("expected slice to contain 3")
	}

	if s.Contains(6) {
		t.Fatal("expected slice to not contain 6")
	}
}

func TestSliceFilter(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}
	got := s.Filter(func(x int) bool { return x%2 == 0 })
	want := Slice[int]{2, 4}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected filtered slice: got %v, want %v", got, want)
	}
}

func TestSliceFirst(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	first := s.First()
	if first == nil {
		t.Fatal("expected non-nil pointer for non-empty slice")
	}
	if *first != 1 {
		t.Fatalf("unexpected first value: got %d, want 1", *first)
	}

	empty := Slice[int]{}
	if empty.First() != nil {
		t.Fatal("expected nil pointer for empty slice")
	}
}

func TestSliceLast(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	last := s.Last()
	if last == nil {
		t.Fatal("expected non-nil pointer for non-empty slice")
	}
	if *last != 3 {
		t.Fatalf("unexpected last value: got %d, want 3", *last)
	}

	empty := Slice[int]{}
	if empty.Last() != nil {
		t.Fatal("expected nil pointer for empty slice")
	}
}

func TestSliceIsEmpty(t *testing.T) {
	if !(Slice[int]{}).IsEmpty() {
		t.Fatal("expected empty slice to be empty")
	}

	if (Slice[int]{1}).IsEmpty() {
		t.Fatal("expected non-empty slice to not be empty")
	}
}

func TestSliceMap(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	got := s.Map(func(x int) int { return x * 2 })
	want := Slice[int]{2, 4, 6}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected mapped slice: got %v, want %v", got, want)
	}
}

func TestSliceForEach(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	total := 0

	s.ForEach(func(x int) {
		total += x
	})

	if total != 6 {
		t.Fatalf("unexpected total: got %d, want 6", total)
	}
}

func TestSliceReduce(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}
	got := s.Reduce(func(acc, x int) int { return acc + x }, 0)

	if got != 15 {
		t.Fatalf("unexpected reduced value: got %d, want 15", got)
	}
}

func TestSliceReverse(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4, 5}
	got := s.Reverse()
	want := Slice[int]{5, 4, 3, 2, 1}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected reversed slice: got %v, want %v", got, want)
	}
}

func TestSliceUnique(t *testing.T) {
	s := Slice[int]{1, 1, 2, 2, 3, 3}
	got := s.Unique()
	want := Slice[int]{1, 2, 3}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected unique slice: got %v, want %v", got, want)
	}
}

func TestSliceLenAndCap(t *testing.T) {
	s := make(Slice[int], 2, 5)

	if s.Len() != 2 {
		t.Fatalf("unexpected len: got %d, want 2", s.Len())
	}

	if s.Cap() != 5 {
		t.Fatalf("unexpected cap: got %d, want 5", s.Cap())
	}
}

func TestSliceAppend(t *testing.T) {
	s := Slice[int]{1, 2}
	got := s.Append(3, 4)
	want := Slice[int]{1, 2, 3, 4}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected appended slice: got %v, want %v", got, want)
	}

	if s.Len() != 2 {
		t.Fatalf("append should not change original slice length: got %d, want 2", s.Len())
	}
}

func TestToMap(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}

	users := Slice[user]{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	got := ToMap(users, func(u user) int { return u.ID })

	if len(got) != 2 {
		t.Fatalf("unexpected map size: got %d, want 2", len(got))
	}

	if got[1].Name != "Alice" {
		t.Fatalf("unexpected user for key 1: got %v", got[1])
	}

	if got[2].Name != "Bob" {
		t.Fatalf("unexpected user for key 2: got %v", got[2])
	}
}
