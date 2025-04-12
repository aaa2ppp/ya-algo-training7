package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`5 3
Cappuccino 25
Car 5
Food 4
Apartment 1
Shopping 7
`)},
			`4 9
Apartment
Car
Food
Shopping
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1 1
event1 100`)},
			`0 0`,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func emptyIfNil[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	}
	return slice
}

type test struct {
	name           string
	maxWeightDif   int
	events         []event
	wantTotalDays  int
	wantEventNames []string
	debug          bool
}

type generator func() test

func test_solve(t *testing.T, tt test, solve solveFunc) {
	t.Run(tt.name, func(t *testing.T) {

		defer func(v bool) { debugEnable = v }(debugEnable)
		debugEnable = tt.debug

		var once sync.Once
		showTest := func() {
			var sb strings.Builder
			fmt.Fprintln(&sb, len(tt.events), tt.maxWeightDif)
			for _, event := range tt.events {
				fmt.Fprintln(&sb, event.name, event.weight)
			}
			t.Log(sb.String())
		}

		gotTotalDays, gotEventNames := solve(tt.maxWeightDif, slices.Clone(tt.events))
		if gotTotalDays != tt.wantTotalDays {
			once.Do(showTest)
			t.Errorf("solve() gotTotalDays = %v, want %v", gotTotalDays, tt.wantTotalDays)
		}
		if !reflect.DeepEqual(emptyIfNil(gotEventNames), emptyIfNil(tt.wantEventNames)) {
			once.Do(showTest)
			t.Errorf("solve() gotEventNames = %v, want %v", gotEventNames, tt.wantEventNames)
		}
	})
}

func Test_solve(t *testing.T) {
	tests := []test{
		{
			"1",
			3,
			[]event{
				{"Cappuccino", 25},
				{"Car", 5},
				{"Food", 4},
				{"Apartment", 1},
				{"Shopping", 7},
			},
			9,
			[]string{
				"Apartment",
				"Car",
				"Food",
				"Shopping",
			},
			true,
		},
		{
			"2",
			1,
			[]event{
				{"event1", 100},
			},
			0,
			[]string{},
			true,
		},
	}

	for _, tt := range tests {
		test_solve(t, tt, solve)
	}
}

func randomWord(rand *rand.Rand, maxLen int) string {
	var sb strings.Builder
	n := rand.Intn(maxLen) + 1
	for i := 0; i < n; i++ {
		sb.WriteByte('a' + byte(rand.Intn(26)))
	}
	return sb.String()
}

func oneEventOneDay(rand *rand.Rand, n int) test {
	events := make([]event, n)
	names := make([]string, n)
	for i := 1; i <= n; i++ {
		name := randomWord(rand, 40)
		events[i-1] = event{name, i}
		names[i-1] = name
	}
	slices.Sort(names)
	return test{
		"oneEventOneDay",
		n,
		events,
		n,
		names,
		false,
	}
}

func Test_solve_oneEventOneDay(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	for i := 0; i < 10; i++ {
		test_solve(t, oneEventOneDay(rand, 10), solve)
	}
}

var (
	bench_totalDays  int
	bench_eventNames []string
)

func bigValues() test {
	// Ограничения:
	// 	1 <= N <= 1000
	// 	1 <= D <= 1000
	// 	1 <= mi <= 10000
	// 	1 <= |namei| <= 40

	rand := rand.New(rand.NewSource(1))
	n := 1000
	d := 1000
	maxM := 10000
	events := make([]event, n)
	for i := 0; i < n; i++ {
		m := rand.Intn(maxM) + 1
		name := randomWord(rand, 40)
		events[i] = event{name, m}
	}
	return test{
		"bigValues",
		d,
		events,
		0,
		nil,
		false,
	}
}

func Benchmark_solve(b *testing.B) {
	tt := bigValues()
	buf := make([]event, len(tt.events))
	for i := 0; i < b.N; i++ {
		copy(buf, tt.events)
		bench_totalDays, bench_eventNames = solve(tt.maxWeightDif, buf)
	}
}

// Генераторы тестовых данных от DeepSeek -- или заставь дурака богу молиться... :(
// Но лучше чем ничего...
func genAllInOneDay() test {
	events := []event{
		{"a", 1}, {"b", 2}, {"c", 3},
	}
	D := 6
	names := []string{"a", "b", "c"}
	return test{"genAllInOneDay", D, events, 3, names, true}
}

func genNoCancellations() test {
	events := []event{{"x", 2}, {"y", 3}}
	return test{"genNoCancellations", 1, events, 0, []string{}, true}
}

func genOptimalReuse() test {
	mi := []int{5, 1, 1, 1, 1, 1}
	names := []string{"a", "b", "c", "d", "e", "f"}
	rand := rand.New(rand.NewSource(1))
	perm := rand.Perm(6)
	events := make([]event, 6)
	for i, p := range perm {
		events[i] = event{names[p], mi[p]}
	}
	sort.Strings(names)
	return test{"genOptimalReuse", 1, events, 10, names, true}
}

func genNoReuse() test {
	events := []event{
		{"m1", 3}, {"m2", 3}, {"m3", 3}, {"m4", 3},
	}
	names := []string{"m1", "m2", "m3", "m4"}
	return test{"genNoReuse", 3, events, 4, names, true}
}

func genMinimal() test {
	return test{"genMinimal", 1, []event{{"single", 1}}, 1, []string{"single"}, true}
}

func genLargeInput() test {
	n := 1000
	events := make([]event, n)
	names := make([]string, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("event%d", i+1)
		events[i] = event{name, 1}
		names[i] = name
	}
	sort.Strings(names)
	return test{"genLargeInput", 1000, events, 1000, names, false}
}

// Тестовые функции
func Test_solve_fromDeepSeek(t *testing.T) {
	generators := []generator{
		genAllInOneDay,
		genNoCancellations,
		genOptimalReuse,
		genNoReuse,
		genMinimal,
		genLargeInput,
	}

	for _, generate := range generators {
		test_solve(t, generate(), solve)
	}
}
