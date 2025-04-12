package main

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"slices"
	"strings"
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
			args{strings.NewReader(`3 7
4 1 2
3 1 2
2 1 2
`)},
			`3 3
1 2 3 
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`3 7
4 1 3
3 1 2
2 1 1
`)},
			`2 2
2 3 
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1 0
1 0 0`)},
			`0 0`,
			true,
		},
		{
			"3.1",
			args{strings.NewReader(`1 0
1 0 1`)},
			`1 0
1`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1 1
1 0 0`)},
			`1 0
1`,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)

			// if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
			// 	t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			// }

			// XXX: Now check only first line. To full check use Test_solve
			if gotOut := out.String(); lines(trimLines(gotOut))[0] != lines(tt.wantOut)[0] {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func lines(s string) []string {
	s = strings.TrimRight(s, "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return lines
}

func trimLines(text string) string {
	lines := lines(text)
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func Test_solve(t *testing.T) {
	type args struct {
		maxVolume int
		items     []item
	}
	tests := []struct {
		name      string
		args      args
		wantCost  int
		wantItems []int
		debug     bool
	}{
		{
			"1",
			args{
				7,
				[]item{
					{1, 4, 1, 2},
					{2, 3, 1, 2},
					{3, 2, 1, 2},
				},
			},
			3, []int{1, 2, 3},
			true,
		},
		{
			"2",
			args{
				7,
				[]item{
					{1, 4, 1, 3},
					{2, 3, 1, 2},
					{3, 2, 1, 1},
				},
			},
			2, []int{2, 3},
			true,
		},
		{
			"3",
			args{
				0,
				[]item{
					{1, 1, 0, 0},
				},
			},
			0, []int{},
			true,
		},
		{
			"3.1",
			args{
				0,
				[]item{
					{1, 1, 0, 1},
				},
			},
			0, []int{1},
			true,
		},
		{
			"4",
			args{
				1,
				[]item{
					{1, 1, 0, 0},
				},
			},
			0, []int{1},
			true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			gotCost, gotItems := solve(tt.args.maxVolume, slices.Clone(tt.args.items))
			if gotCost != tt.wantCost {
				t.Errorf("solve() cost = %v, want %v", gotCost, tt.wantCost)
			}

			checkItems := func(gotItems []int, wantCost int) bool {
				totalCost := 0
				volume := 0
				for _, num := range gotItems {
					it := tt.args.items[num-1]
					totalCost += it.cost
					volume += it.volume
				}

				if totalCost != wantCost {
					t.Errorf("total cost of items = %v, want %v", totalCost, wantCost)
					return false
				}

				pressure := max(0, volume-tt.args.maxVolume)

				for _, num := range gotItems {
					it := tt.args.items[num-1]
					if it.pressure < pressure {
						t.Errorf("pressure of item #%v = %v, pressure %v", it.id, it.pressure, pressure)
						return false
					}
				}

				return true
			}

			if !reflect.DeepEqual(gotItems, tt.wantItems) && !checkItems(gotItems, tt.wantCost) {
				t.Errorf("solve() gotItems = %v, wantItems %v", gotItems, tt.wantItems)
			}
		})
	}
}

type test struct {
	name      string
	maxVolume int
	items     []item
	wantCost  int
	wantItems []int
	debug     bool
}

type generator func() test

func bigValues() test {
	// Ограничения:
	// 	1 <= N <= 100
	// 	0 <= S <= 10^9
	// 	1 <= vi <= 1000
	// 	0 <= ci <= 10^6
	// 	0 <= pi <= 10^9

	rand := rand.New(rand.NewSource(1))
	n := 100
	items := make([]item, n)
	totalV := 0
	for i := 0; i < n; i++ {
		v := rand.Intn(1000) + 1
		c := rand.Intn(1e6)
		p := rand.Intn(1000) // не имеет смысла делать очень большим
		items[i] = item{i + 1, v, c, p}
		totalV += v
	}
	return test{
		"bigValues",
		totalV * 2 / 3, // не имет смысла делать очень большим, все что > totalV обрезаю
		items,
		0,
		nil,
		false,
	}
}

var (
	bench_cost  int
	bench_items []int
)

func Benchmark_solve(b *testing.B) {
	tt := bigValues()
	buf := make([]item, len(tt.items))
	for i := 0; i < b.N; i++ {
		copy(buf, tt.items)
		bench_cost, bench_items = solve(tt.maxVolume, buf)
	}
}

// в этот раз показал ему код решения...
func TestSolve_fromDeepSeek(t *testing.T) {
	// Корректные случаи:
	// Тест 1: Все предметы помещаются без превышения базового объёма.
	// Тест 2: Предметы превышают базовый объём, но давление допустимо для всех.
	// Тест 3: Один из предметов не может быть выбран из-за давления.
	// Тест 4: Проверка случая с S=0.
	// Тест 8: Выбор оптимального подмножества с учётом давления и стоимости.

	// Граничные случаи:
	// Тест 5: Один допустимый предмет.
	// Тест 6: Предмет не может быть выбран из-за недостаточного давления.
	// Тест 7: Пустой список предметов.

	// Сортировка:
	// Для упрощения проверки ID сортируются, так как порядок вывода может зависеть от реализации алгоритма.

	tests := []struct {
		name      string
		maxVolume int
		items     []item
		wantCost  int
		wantIDs   []int
	}{
		{
			name:      "all items fit without pressure",
			maxVolume: 10,
			items: []item{
				{id: 1, volume: 3, cost: 10, pressure: 0},
				{id: 2, volume: 2, cost: 20, pressure: 5},
			},
			wantCost: 30,
			wantIDs:  []int{1, 2},
		},
		{
			name:      "items fit with pressure within limits",
			maxVolume: 5,
			items: []item{
				{id: 1, volume: 3, cost: 10, pressure: 2},
				{id: 2, volume: 3, cost: 20, pressure: 3},
			},
			wantCost: 30,
			wantIDs:  []int{1, 2},
		},
		{
			name:      "select subset due to pressure",
			maxVolume: 5,
			items: []item{
				{id: 1, volume: 4, cost: 10, pressure: 1},
				{id: 2, volume: 3, cost: 20, pressure: 2},
			},
			wantCost: 20,
			wantIDs:  []int{2},
		},
		{
			name:      "S=0 with valid items",
			maxVolume: 0,
			items: []item{
				{id: 1, volume: 2, cost: 10, pressure: 2},
				{id: 2, volume: 3, cost: 20, pressure: 3},
			},
			wantCost: 20,
			wantIDs:  []int{2},
		},
		{
			name:      "single valid item",
			maxVolume: 5,
			items: []item{
				{id: 1, volume: 3, cost: 10, pressure: 0},
			},
			wantCost: 10,
			wantIDs:  []int{1},
		},
		{
			name:      "single invalid item due to pressure",
			maxVolume: 5,
			items: []item{
				{id: 1, volume: 6, cost: 10, pressure: 0},
			},
			wantCost: 0,
			wantIDs:  []int{},
		},
		{
			name:      "no items",
			maxVolume: 10,
			items:     []item{},
			wantCost:  0,
			wantIDs:   []int{},
		},
		{
			name:      "choose between items with different pressure and cost",
			maxVolume: 5,
			items: []item{
				{id: 1, volume: 4, cost: 20, pressure: 2},
				{id: 2, volume: 3, cost: 25, pressure: 1},
				{id: 3, volume: 3, cost: 15, pressure: 3},
			},
			wantCost: 40,
			wantIDs:  []int{3, 2}, // После сортировки и выбора
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCost, gotIDs := solve(tt.maxVolume, tt.items)
			if gotCost != tt.wantCost {
				t.Errorf("solve() gotCost = %v, want %v", gotCost, tt.wantCost)
			}

			// Сортируем ID для сравнения, если порядок не важен
			slices.Sort(gotIDs)
			slices.Sort(tt.wantIDs)
			if !slices.Equal(gotIDs, tt.wantIDs) {
				t.Errorf("solve() gotIDs = %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func BenchmarkSolve_fromDeepSeek(b *testing.B) {
	// Особенности реализации:

	// Генерация реалистичных данных:
	// Для smallItems, mediumItems и largeItems генерируются элементы с линейно растущими характеристиками
	// Каждый следующий элемент имеет:
	// Объём: 100 + i*20
	// Стоимость: 500 + i*100
	// Давление: 50 + i*15

	// Изоляция тестовых данных:
	// Перед каждым запуском создаётся копия исходных данных
	// Решает проблему с модификацией слайса в функции solve (сортировка элементов)

	// Раздельные бенчмарки:
	// Small: 10 элементов (базовый случай)
	// Medium: 50 элементов (проверка средней нагрузки)
	// Large: 100 элементов (максимальное разрешённое количество)
	// AllFit: специальный случай, когда все элементы влезают без перегрузки

	// Типичные сценарии:
	// Комбинация элементов с разной "рентабельностью"
	// Элементы с возрастающим допустимым давлением
	// Реалистичные значения объёмов и стоимости

	// Генерация тестовых данных
	smallItems := generateItems(10)
	mediumItems := generateItems(50)
	largeItems := generateItems(100)

	// Бенчмарки для разных размеров входных данных
	b.Run("Small", func(b *testing.B) { benchmarkSolve(b, 1000, smallItems) })
	b.Run("Medium", func(b *testing.B) { benchmarkSolve(b, 5000, mediumItems) })
	b.Run("Large", func(b *testing.B) { benchmarkSolve(b, 20000, largeItems) })

	// Специальный случай: все элементы влезают без давления
	b.Run("AllFit", func(b *testing.B) {
		items := []item{
			{id: 1, volume: 50, cost: 100, pressure: 100},
			{id: 2, volume: 70, cost: 150, pressure: 150},
			{id: 3, volume: 30, cost: 200, pressure: 200},
		}
		benchmarkSolve(b, 500, items)
	})
}

func benchmarkSolve(b *testing.B, maxVolume int, items []item) {
	// Подготовка копии элементов для каждого запуска
	origItems := make([]item, len(items))
	copy(origItems, items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создаем новую копию для каждой итерации
		itemsCopy := make([]item, len(origItems))
		copy(itemsCopy, origItems)
		solve(maxVolume, itemsCopy)
	}
}

func generateItems(n int) []item {
	items := make([]item, n)
	for i := 0; i < n; i++ {
		volume := 100 + i*20
		cost := 500 + i*100
		pressure := 50 + i*15
		items[i] = item{
			id:       i + 1,
			volume:   volume,
			cost:     cost,
			pressure: pressure,
		}
	}
	return items
}
