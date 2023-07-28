package repository

type ErrObjectNotFound struct {
}

func (o ErrObjectNotFound) Error() string {
	return "object not found"
}
