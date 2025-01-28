package excelize

import (
	"strconv"
	"testing"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache()
	key := "Sheet1!A1"
	value := formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	}
	cache.Add(key, value)

	result, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Expected to find key %s in cache, but it was not found", key)
	}
	if result.value.Number != value.value.Number || result.value.Type != value.value.Type {
		t.Fatalf("Expected value %v, got %v", value.value, result.value)
	}
}

func TestCacheLRUEviction(t *testing.T) {
	cache := NewCache()
	cache.SetLimit(2)

	cache.Add("Sheet1!A1", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	})
	cache.Add("Sheet1!A2", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 100},
		err:   nil,
	})

	// Access the first item to make it recently used
	cache.Get("Sheet1!A1")

	// Add a third item, which should evict the least recently used item ("Sheet1!A2")
	cache.Add("Sheet1!A3", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 200},
		err:   nil,
	})

	_, ok := cache.Get("Sheet1!A2")
	if ok {
		t.Fatalf("Expected key Sheet1!A2 to be evicted, but it was found in the cache")
	}
	_, ok = cache.Get("Sheet1!A1")
	if !ok {
		t.Fatalf("Expected key Sheet1!A1 to be in the cache, but it was not found")
	}
	_, ok = cache.Get("Sheet1!A3")
	if !ok {
		t.Fatalf("Expected key Sheet1!A3 to be in the cache, but it was not found")
	}
}

func TestCacheDisable(t *testing.T) {
	cache := NewCache()
	key := "Sheet1!A1"
	value := formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	}
	cache.Add(key, value)

	cache.DisableCache()

	_, ok := cache.Get(key)
	if ok {
		t.Fatalf("Expected cache to be disabled, but key %s was found in the cache", key)
	}

	cache.Add("Sheet1!A2", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 100},
		err:   nil,
	})
	_, ok = cache.Get("Sheet1!A2")
	if ok {
		t.Fatalf("Expected cache to be disabled, but key Sheet1!A2 was found in the cache")
	}
}

func TestCacheInvalidate(t *testing.T) {
	cache := NewCache()

	key := "Sheet1!A1"
	value := formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	}
	cache.Add(key, value)

	cache.Invalidate()

	_, ok := cache.Get(key)
	if ok {
		t.Fatalf("Expected cache to be invalidated, but key %s was found in the cache", key)
	}
}

func TestCacheSetLimit(t *testing.T) {
	cache := NewCache()
	cache.SetLimit(1)

	cache.Add("Sheet1!A1", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	})
	cache.Add("Sheet1!A2", formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 100},
		err:   nil,
	})

	_, ok := cache.Get("Sheet1!A1")
	if ok {
		t.Fatalf("Expected key Sheet1!A1 to be evicted, but it was found in the cache")
	}
	_, ok = cache.Get("Sheet1!A2")
	if !ok {
		t.Fatalf("Expected key Sheet1!A2 to be in the cache, but it was not found")
	}
}

func TestCacheDisabledAddGet(t *testing.T) {
	cache := NewCache()
	cache.DisableCache()
	key := "Sheet1!A1"
	value := formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	}
	cache.Add(key, value)
	_, ok := cache.Get(key)
	if ok {
		t.Fatalf("Expected cache to be disabled, but key %s was found in the cache", key)
	}
}

func TestCacheDisabledInvalidate(t *testing.T) {
	cache := NewCache()
	cache.DisableCache()
	key := "Sheet1!A1"
	value := formulaResult{
		value: formulaArg{Type: ArgNumber, Number: 42},
		err:   nil,
	}
	cache.Add(key, value)
	cache.Invalidate()

	_, ok := cache.Get(key)
	if ok {
		t.Fatalf("Expected cache to be disabled, but key %s was found in the cache", key)
	}
}

// BenchmarkWithCache benchmarks creating a file, adding cells, and calculating values with the cache enabled.
func BenchmarkWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := NewFile()

		// Add 1000 rows and 5 columns of data with values and formulas
		for j := 1; j <= 1000; j++ {
			// Set values in columns A, B, C, D, E
			f.SetCellValue("Sheet1", "A"+strconv.Itoa(j), j)                                           // Column A: Simple value
			f.SetCellValue("Sheet1", "B"+strconv.Itoa(j), j*2)                                         // Column B: Another value
			f.SetCellFormula("Sheet1", "C"+strconv.Itoa(j), "=A"+strconv.Itoa(j)+"+B"+strconv.Itoa(j)) // Column C: Sum of A and B
			f.SetCellFormula("Sheet1", "D"+strconv.Itoa(j), "=C"+strconv.Itoa(j)+"*2")                 // Column D: Double the value of C
			f.SetCellFormula("Sheet1", "E"+strconv.Itoa(j), "=D"+strconv.Itoa(j)+"-A"+strconv.Itoa(j)) // Column E: Subtract A from D
		}

		// Calculate all formulas multiple times to simulate repeated calculations
		for k := 0; k < 10; k++ { // Run calculations 10 times
			for j := 1; j <= 1000; j++ {
				// Calculate formulas in columns C, D, E
				_, _ = f.CalcCellValue("Sheet1", "C"+strconv.Itoa(j)) // Sum of A and B
				_, _ = f.CalcCellValue("Sheet1", "D"+strconv.Itoa(j)) // Double the value of C
				_, _ = f.CalcCellValue("Sheet1", "E"+strconv.Itoa(j)) // Subtract A from D
			}
		}
	}
}

// BenchmarkWithoutCache benchmarks creating a file, adding cells, and calculating values with the cache disabled.
func BenchmarkWithoutCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := NewFile()

		// Disable the cache
		f.DisableCache()

		// Add 1000 rows and 5 columns of data with values and formulas
		for j := 1; j <= 1000; j++ {
			// Set values in columns A, B, C, D, E
			f.SetCellValue("Sheet1", "A"+strconv.Itoa(j), j)                                           // Column A: Simple value
			f.SetCellValue("Sheet1", "B"+strconv.Itoa(j), j*2)                                         // Column B: Another value
			f.SetCellFormula("Sheet1", "C"+strconv.Itoa(j), "=A"+strconv.Itoa(j)+"+B"+strconv.Itoa(j)) // Column C: Sum of A and B
			f.SetCellFormula("Sheet1", "D"+strconv.Itoa(j), "=C"+strconv.Itoa(j)+"*2")                 // Column D: Double the value of C
			f.SetCellFormula("Sheet1", "E"+strconv.Itoa(j), "=D"+strconv.Itoa(j)+"-A"+strconv.Itoa(j)) // Column E: Subtract A from D
		}

		// Calculate all formulas multiple times to simulate repeated calculations
		for k := 0; k < 10; k++ { // Run calculations 10 times
			for j := 1; j <= 1000; j++ {
				// Calculate formulas in columns C, D, E
				_, _ = f.CalcCellValue("Sheet1", "C"+strconv.Itoa(j)) // Sum of A and B
				_, _ = f.CalcCellValue("Sheet1", "D"+strconv.Itoa(j)) // Double the value of C
				_, _ = f.CalcCellValue("Sheet1", "E"+strconv.Itoa(j)) // Subtract A from D
			}
		}
	}
}

func TestCacheHits(t *testing.T) {
	f := NewFile()
	f.SetCacheLimit(10)
	f.SetCellValue("Sheet1", "A1", 50)
	f.SetCellFormula("Sheet1", "A2", "=A1*2")
	f.SetCellFormula("Sheet1", "A3", "=A1+A2")
	f.SetCellFormula("Sheet1", "A4", "=A1+A2+A3+10")

	// Initial calculation
	_, _ = f.CalcCellValue("Sheet1", "A1")
	_, _ = f.CalcCellValue("Sheet1", "A2") // will hit 1 cache when accessing A1
	_, _ = f.CalcCellValue("Sheet1", "A3") // will hit 2 cache when accessing A1 and A2
	_, _ = f.CalcCellValue("Sheet1", "A4") // will hit 4 cache when accessing A1, A2, A3

	if f.cache.hits != 7 {
		t.Fatalf("Expected 7 cache hits, got %d", f.cache.hits)
	}
	if f.cache.misses != 3 {
		t.Fatalf("Expected 3 cache misses, got %d", f.cache.misses)
	}
}
