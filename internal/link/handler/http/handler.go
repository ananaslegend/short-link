package http

type Link struct {
	linkGetter   LinkGetter
	linkInserter LinkInserter
}

func New(linkGetter LinkGetter, linkInserter LinkInserter) *Link {
	return &Link{
		linkGetter:   linkGetter,
		linkInserter: linkInserter,
	}
}
