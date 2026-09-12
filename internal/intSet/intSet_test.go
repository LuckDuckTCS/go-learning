package intset

import (
	"maps"
	"math/rand"
	"testing"
)

// тест на равносильность intSet и map[int]bool по длине
func TestIntSetAgainstMap(t *testing.T) {
	rng := rand.New(rand.NewSource(0))

	var s IntSet
	m := make(map[int]bool)

	for i := 0; i < 1000; i++ {
		x := rng.Intn(500)

		switch rng.Intn(3) {
		case 0:
			// add в оба
			m[x] = true
			s.Add(x)
		case 1:
			delete(m, x)
			s.Remove(x)
		case 2:
			if s.Has(x) != m[x] {
				t.Fatalf("Has(%d) = %v, want %v", x, s.Has(x), m[x])
			}
		}
		if s.Len() != len(m) {
			t.Fatalf("Len %v, want %v", s.Len(), len(m))
		}
	}
	for x := 0; x < 500; x++ {
		if s.Has(x) != m[x] {
			t.Errorf("Has(%d) = %v, want %v", x, s.Has(x), m[x])
		}
	}
}

func TestIntSetBinaryOps(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	maxValue := 500

	for i := 0; i < 200; i++ {

		count1 := rng.Intn(200)
		count2 := rng.Intn(200)
		opName := ""
		// случайное наполнение двух множеств
		sVals := randomSet(rng, maxValue, count1)
		tVals := randomSet(rng, maxValue, count2)
		s := setFromMap(sVals)
		tm := setFromMap(tVals)
		// случайный выбор одной из четырёх операций
		switch rng.Intn(4) {
		case 0:
			// union: s.UnionWith(t) против объединения карт
			sVals = unionMaps(sVals, tVals)
			s.UnionWith(tm)
			opName = "union"
		case 1:
			sVals = intersectMaps(sVals, tVals)
			s.IntersectWith(tm)
			opName = "intersect"
		case 2: // difference
			sVals = differenceMaps(sVals, tVals)
			s.DifferenceWith(tm)
			opName = "difference"
		case 3: // symmetric difference
			sVals = symmetricDiffMaps(sVals, tVals)
			s.SymmetricDifference(tm)
			opName = "symmetricDiff"
		}
		if s.Len() != len(sVals) {
			t.Fatalf("iter %d: Len %v, want %v, op %s", i, s.Len(), len(sVals), opName)
		}
		for x := 0; x < maxValue; x++ {
			if s.Has(x) != sVals[x] {
				t.Fatalf("iter %d: Has(%d) = %v, want %v, op %s", i, x, s.Has(x), sVals[x], opName)
			}
		}

	}

}

func randomSet(rng *rand.Rand, maxValue, count int) map[int]bool {
	m := make(map[int]bool)
	for i := 0; i < count; i++ {
		x := rng.Intn(maxValue)
		m[x] = true
	}
	return m
}

func setFromMap(m map[int]bool) *IntSet {
	var result IntSet
	for key := range m {
		result.Add(key)
	}
	return &result
}

func unionMaps(a, b map[int]bool) map[int]bool {
	m := maps.Clone(b)
	for key := range a {
		m[key] = true
	}
	return m
}

func intersectMaps(a, b map[int]bool) map[int]bool {
	m := make(map[int]bool)
	for key := range a {
		if b[key] {
			m[key] = true
		}
	}
	return m
}

func differenceMaps(a, b map[int]bool) map[int]bool {
	m := maps.Clone(a)
	for key := range b {
		delete(m, key)
	}
	return m
}

func symmetricDiffMaps(a, b map[int]bool) map[int]bool {
	return unionMaps(differenceMaps(a, b), differenceMaps(b, a))
}
