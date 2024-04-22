package listutil

import d "github.com/clearblade/cblib/diff"

func CompareLists(after []interface{}, before []interface{}, compare func(interface{}, interface{}) bool) *d.UnsafeDiff {
	diff := d.UnsafeDiff{after, before, nil, nil, compare}
	d.Diff(&diff)
	return &diff
}

func CompareListsAndFilter(after []interface{}, before []interface{}, compare func(interface{}, interface{}) bool, filter func(interface{}) bool) *d.UnsafeDiff {
	diff := d.UnsafeDiff{after, before, nil, nil, compare}
	d.Diff(&diff)
	diff.Added = FilterSlice(diff.Added, filter)
	diff.Removed = FilterSlice(diff.Removed, filter)
	return &diff
}

// FilterSlice returns the items of the slice `s` for which `predicate` returns true.
func FilterSlice(s []interface{}, predicate func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0, len(s))

	for _, item := range s {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
