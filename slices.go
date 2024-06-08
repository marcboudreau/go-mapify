package mapify

// FromSlice creates a map using the provided slice of E elements and the
// key function to determine the map key for each of the elements in the slice.
// If the key function returns the same key for multiple elements, the previous
// element stored with the duplicated key will be overwritten. To create a map
// that can handle duplicate keys, see FromSliceWithDuplicates.
func FromSlice[E any, K comparable](s []E, key func(e E) K) map[K]E {
	m := make(map[K]E)

	for _, e := range s {
		k := key(e)

		m[k] = e
	}

	return m
}

// FromSliceWithDuplicates creates a map using the provided slice of E elements
// and the key function to determine the map key for each of the elements in the
// slice. The slice elements are stored slices in the map, so if the key
// function returns the same key for multiple slice elements, they will all be
// stored in the same slice in the map.
func FromSliceWithDuplicates[E any, K comparable](s []E, walk func(e E) K) map[K][]E {
	m := make(map[K][]E)

	for _, e := range s {
		k := walk(e)

		m[k] = append(m[k], e)
	}

	return m
}
