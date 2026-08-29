package expslices

import (
	"unicode"
	"unicode/utf8"
)

// упр 4.3 Донован Перепишите функцию reverse так, чтобы вместо среза она использовала указатель на массив.

func Reverse(s *[5]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Упражнение 4.4. Напишите версию функции rotate, которая работает в один проход.
// сдвинем на n влево
func Rotate(s []int, n int) {
	if n <= 0 || len(s) == 0 {
		return
	}
	n = n % len(s)

	if n == 0 {
		return
	}

	temp := make([]int, n)
	copy(temp, s[:n])
	copy(s, s[n:])
	copy(s[(len(s)-n):], temp)
}

// Упражнение 4.5. Напишите функцию, которая без выделения дополнительной памяти
// удаляет все смежные дубликаты в срезе []string.

func NoDuplicate(s []string) []string {
	if len(s) <= 1 {
		return s
	}
	j := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[j-1] {
			s[j] = s[i]
			j++
		}
	}
	return s[:j]
}

//Упражнение 4.6. Напишите функцию, которая без выделения дополнительной памяти
//преобразует последовательности смежных пробельных символов Unicode (см. Unicode. IsSpace)
//в срезе []byte в кодировке UTF-8 в один пробел ASCII.

func NoDuplicateSpaces(b []byte) []byte {
	j := 0
	wasSpace := false
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if unicode.IsSpace(r) {
			if !wasSpace {
				b[j] = ' '
				j++
				wasSpace = true
			}

		} else {
			if i != j {
				copy(b[j:j+size], b[i:i+size])
			}
			j += size
			wasSpace = false
		}

		i += size
	}
	return b[:j]
}

//Упражнение 4.7. Перепишите функцию reverse так, чтобы она без выделения
//дополнительной памяти обращала последовательность символов среза []byte, который
// представляет строку в кодировке UTF-8. Сможете ли вы обойтись без выделения новой памяти?

func ReverseBytes(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ReverseNew(b []byte) {
	for i := 0; i < len(b); {
		_, size := utf8.DecodeRune(b[i:])
		ReverseBytes(b[i : i+size])
		i += size
	}

	ReverseBytes(b)
}
