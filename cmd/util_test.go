package cmd

type exitMemory struct {
	code int
}

func (e *exitMemory) Exit(i int) {
	e.code = i
}
