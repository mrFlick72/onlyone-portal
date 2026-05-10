package plan

import "testing"

func TestStatusValidity(t *testing.T) {
	cases := map[TodoStatus]bool{
		StatusTodo:        true,
		StatusInProgress:  true,
		StatusDone:        true,
		StatusAborted:     true,
		TodoStatus(""):    false,
		TodoStatus("FOO"): false,
	}
	for s, want := range cases {
		if got := s.IsValid(); got != want {
			t.Errorf("IsValid(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestAllowedTransitions(t *testing.T) {
	type tr struct {
		from, to TodoStatus
		want     bool
	}
	cases := []tr{
		{StatusTodo, StatusInProgress, true},
		{StatusTodo, StatusAborted, true},
		{StatusTodo, StatusDone, false},
		{StatusTodo, StatusTodo, false},
		{StatusInProgress, StatusDone, true},
		{StatusInProgress, StatusAborted, true},
		{StatusInProgress, StatusTodo, false},
		{StatusInProgress, StatusInProgress, false},
		{StatusDone, StatusTodo, false},
		{StatusDone, StatusInProgress, false},
		{StatusDone, StatusAborted, false},
		{StatusAborted, StatusTodo, false},
		{StatusAborted, StatusInProgress, false},
		{StatusAborted, StatusDone, false},
	}
	for _, c := range cases {
		if got := c.from.CanTransitionTo(c.to); got != c.want {
			t.Errorf("%s -> %s: got %v, want %v", c.from, c.to, got, c.want)
		}
	}
}
