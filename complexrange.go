package mkf

func newComplexRange(base [2]rune, excludes [][2]rune) *complexRange {
	if base[0] > base[1] {
		return nil
	}

	for _, v := range excludes {
		if v[0] > v[1] {
			return nil
		}
	}

	return &complexRange{
		base:     base,
		excludes: excludes,
	}
}
