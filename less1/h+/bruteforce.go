package main

func bruteforce(orders []string) int {
	type item struct {
		right    int
		left     int
		inverter bool
	}

	items := make([]item, 0, len(orders))
	for _, order := range orders {
		right, left := countRigthLeft(order)
		items = append(items, item{right, left, len(order)%2 != 0})
	}

	var maximum int
	HeapPermutations(items, func(items []item) {
		count := 0
		right := true
		for _, it := range items {
			if right {
				count += it.right
			} else {
				count += it.left
			}
			if it.inverter {
				right = !right
			}
		}
		maximum = max(maximum, count)
	})

	return maximum
}

func bruteforce2(orders []string) int {
	type item struct {
		vRight   int
		vLeft    int
		right    int
		left     int
		inverter bool
	}

	items := make([]item, 0, len(orders))

	countItem := func(order string) item {
		it := item{inverter: len(order)%2 != 0}
		for i, c := range []byte(order) {
			if c == 'S' {
				it.right += (i + 1) & 1
				it.left += i & 1
				it.vRight += (i + 1) & 1
				it.vLeft += i & 1
			} else { // D
				it.right += i & 1
				it.left += (i + 1) & 1
			}
		}
		return it
	}

	for _, order := range orders {
		items = append(items, countItem(order))
	}

	var maximum, vMaximum int
	GeneratePermutations(items, func(items []item) {
		count, vCount := 0, 0
		right := true
		for _, it := range items {
			if right {
				count += it.right
				vCount += it.vRight
			} else {
				count += it.left
				vCount += it.vLeft
			}
			if it.inverter {
				right = !right
			}
		}
		if count > maximum {
			maximum = count
			vMaximum = vCount
		}
	})

	return vMaximum
}

// GeneratePermutations генерирует все возможные перестановки слайса, вызывая handler для каждой.
//
// Важные предупреждения:
//   - Исходный слайс модифицируется в процессе работы
//   - Callback получает временный буфер, который:
//   - Не должен изменяться в обработчике
//   - Актуален только во время вызова callback
//   - После завершения функции исходный слайс будет в неопределенном состоянии
//
// Производительность:
//   - Не создает копий данных (работает напрямую с переданным слайсом)
//   - Использует рекурсивный алгоритм с backtracking
//   - Глубина рекурсии = длина слайса
//
// Пример:
//
//	arr := []int{1,2,3}
//	GeneratePermutations(arr, func(p []int) {
//	    // p использует тот же буфер что и arr
//	    fmt.Println(p) // OK: чтение данных
//	})
func GeneratePermutations[T any](arr []T, handler func([]T)) {
	if len(arr) == 0 || handler == nil {
		return
	}

	var generate func(int)
	generate = func(index int) {
		if index == len(arr) {
			handler(arr)
			return
		}

		for i := index; i < len(arr); i++ {
			arr[index], arr[i] = arr[i], arr[index]
			generate(index + 1)
			arr[index], arr[i] = arr[i], arr[index] // backtrack
		}
	}

	generate(0)
}

// HeapPermutations генерирует перестановки используя алгоритм Хипа с максимальной производительностью.
//
// Особенности:
//   - Работает непосредственно с переданным слайсом
//   - Не восстанавливает исходный порядок элементов
//   - Callback получает временный буфер:
//   - Данные становятся невалидными после завершения callback
//   - Модификация буфера может привести к некорректной работе
//
// Преимущества:
//   - На 40% меньше операций свопа чем у рекурсивной версии
//   - Нет накладных расходов на вызовы функций
//   - Постоянная память: O(n) для служебного массива
//
// Пример:
//
//	data := []string{"a","b","c"}
//	HeapPermutations(arr, func(p []string) {
//	    // p использует тот же буфер что и arr
//	    process(p)
//	})
//	// data теперь содержит последнюю сгенерированную перестановку
func HeapPermutations[T any](arr []T, handler func([]T)) {
	if len(arr) == 0 || handler == nil {
		return
	}

	n := len(arr)
	c := make([]int, n)

	handler(arr)

	for i := 0; i < n; {
		if c[i] < i {
			if i%2 == 0 {
				arr[0], arr[i] = arr[i], arr[0]
			} else {
				arr[c[i]], arr[i] = arr[i], arr[c[i]]
			}

			handler(arr)
			c[i]++
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}
}
