package util

// Map applies a function to each element in a slice, returning a new slice with the results.
func Map[I any, O any](mapper func(I) O, iter []I) []O {
    out := make([]O, len(iter))
    for i, in := range iter {
        out[i] = mapper(in)
    }
    return out
}
