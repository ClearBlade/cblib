package listutil

import d "github.com/clearblade/cblib/diff"

func CompareLists[T any](after []T, before []T, compare func(T, T) bool) *d.UnsafeDiff[T] {
	diff := d.UnsafeDiff[T]{after, before, nil, nil, compare}
	d.Diff(&diff)
	return &diff
}

func CompareListsAndFilter[T any](after []T, before []T, compare func(T, T) bool, filter func(T) bool) *d.UnsafeDiff[T] {
	diff := d.UnsafeDiff[T]{after, before, nil, nil, compare}
	d.Diff(&diff)
	diff.Added = FilterSlice(diff.Added, filter)
	diff.Removed = FilterSlice(diff.Removed, filter)
	return &diff
}

// FilterSlice returns the items of the slice `s` for which `predicate` returns true.
func FilterSlice[T any](s []T, predicate func(T) bool) []T {
	filtered := make([]T, 0, len(s))

	for _, item := range s {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
