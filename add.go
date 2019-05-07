package dotm

// Add adds the given file to the profile under the given top level dir.
// If the file already exists under $HOME/path, A backup is created and then
// copied to the profile.
func Add(profile, dir, path string) error {
	c, err := LoadConfig()
	if err != nil {
		return err
	}
	p, err := c.Profile(profile)
	if err != nil {
		return err
	}

	return p.addFile(dir, path)
}
