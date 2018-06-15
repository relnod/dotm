package dotfiles

// TODO: make public
func isExcluded(dir string) bool {
	for _, exclude := range excludes {
		if dir == exclude {
			return true
		}
	}

	return false
}
