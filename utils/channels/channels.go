package channels

// TODO: Rework this
func OK(done chan bool) bool {
	select {
	case ok := <-done:
		if ok {
			return ok
		}

	}
	return false
}
