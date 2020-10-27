package resourcetree

// ReadUsing reads a new *Tree using the given TreeReader.
func ReadUsing(r TreeReader) (*Tree, error) {

	var tree Tree
	var err error

	tree.Collections, err = r.ReadCollections()
	if err != nil {
		return nil, err
	}

	// NOTE: add other readers below...

	return &tree, nil
}

// WriteUsing writes a *Tree using the given TreeWriter.
func WriteUsing(w TreeWriter, tree *Tree) error {

	var err error

	err = w.WriteCollections(tree.Collections)
	if err != nil {
		return err
	}

	// NOTE: add other writers below...

	return nil
}
