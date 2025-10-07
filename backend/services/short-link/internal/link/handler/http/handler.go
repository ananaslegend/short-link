package http

type LinkHandler struct {
	linkGetter   LinkGetter
	linkInserter LinkInserter
}

func NewHandler(
	linkGetter LinkGetter,
	linkInserter LinkInserter,
) *LinkHandler {
	return &LinkHandler{
		linkGetter:   linkGetter,
		linkInserter: linkInserter,
	}
}
