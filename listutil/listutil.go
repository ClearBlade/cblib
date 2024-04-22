package listutil

import d "github.com/clearblade/cblib/diff"

func CompareLists(after []interface{}, before []interface{}, compare func(interface{}, interface{}) bool) *d.UnsafeDiff {
	diff := d.UnsafeDiff{after, before, nil, nil, compare}
	d.Diff(&diff)
	return &diff
}
