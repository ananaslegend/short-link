package service

type Link struct {
	linkGetter          LinkGetter
	aliasedLinkInserter AliasedLinkInserter
	aliasGenerator      AliasGenerator
}

func New(
	lp LinkGetter,
	aliasedLinkInserter AliasedLinkInserter,
	aliasGenerator AliasGenerator,
) *Link {
	return &Link{
		linkGetter:          lp,
		aliasedLinkInserter: aliasedLinkInserter,
		aliasGenerator:      aliasGenerator,
	}
}
