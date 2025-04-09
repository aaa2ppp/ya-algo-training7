package main

type stack[T any] []T

func (s stack[T]) empty() bool {
	return len(s) == 0
}

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s stack[T]) top() T {
	n := len(s)
	return s[n-1]
}

func (s *stack[T]) pop() T {
	old := *s
	n := len(old)
	v := old[n-1]
	*s = old[:n-1]
	return v
}

type queue[T any] struct {
	input  stack[T]
	output stack[T]
}

func (q *queue[T]) empty() bool {
	return q.input.empty() && q.output.empty()
}

func (q *queue[T]) len() int {
	return len(q.input) + len(q.output)
}

func (q *queue[T]) push(v T) {
	q.input.push(v)
}

func (q *queue[T]) front() T {
	q.pour()
	return q.output.top()
}

func (q *queue[T]) pop() T {
	q.pour()
	return q.output.pop()
}

func (q *queue[T]) pour() {
	if len(q.output) == 0 {
		for !q.input.empty() {
			q.output.push(q.input.pop())
		}
	}
}
