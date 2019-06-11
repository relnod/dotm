package dotm

// HelpError annotates a help message to an error.
type HelpError struct {
	err  error
	help string
}

// Error returns the error message.
func (h *HelpError) Error() string {
	return h.err.Error()
}

// Help returns the annotated help message.
func (h *HelpError) Help() string {
	return h.help
}

// Unwrap returns the underlying error.
func (h *HelpError) Unwrap() error {
	return h.err
}
