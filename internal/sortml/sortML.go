package sortml

import "sort"

//Упражнение 7.8. Многие графические интерфейсы предоставляют таблицы с
//многоуровневой сортировкой с сохранением состояния: первичный ключ определяется
//по последнему щелчку на заголовке, вторичный — по предпоследнему и т.д.
//Определите реализацию sort.Interface для использования в такой таблице.

type DataT []string

type SortML struct {
	Table     []*DataT
	ColumnSet []int // срез приоритетов
}

func (x *SortML) Len() int { return len(x.Table) }

func (x *SortML) Less(i, j int) bool {

	for _, col := range x.ColumnSet {
		if (*x.Table[i])[col] != (*x.Table[j])[col] {
			return (*x.Table[i])[col] < (*x.Table[j])[col]
		}
	}
	return false
}

func (x *SortML) Swap(i, j int) { x.Table[i], x.Table[j] = x.Table[j], x.Table[i] }

func IsPalindrome(s sort.Interface) bool {
	size := s.Len()
	n := size / 2
	for i := range n {
		if s.Less(i, size-i-1) || s.Less(size-i-1, i) {
			return false
		}
	}
	return true
}
