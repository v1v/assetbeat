package util

func Map[I any, O any](mapper func(I) O, iter []I) []O {
    out := make([]O, len(iter))
    for i, in := range iter {
        out[i] = mapper(in)
    }
    return out
}

func Filter[T any](predicate func(T) bool, iter []T) []T {
    out := make([]T, 0, len(iter))
    for _, in := range iter {
        if predicate(in) {
            out = append(out, in)
        }
    }
    return out
}