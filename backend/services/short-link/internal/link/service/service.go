package service

type Link struct {
	linkGetter                LinkGetter
	aliasedLinkInserter       AliasedLinkInserter
	aliasGenerator            AliasGenerator
	redirectStatisticProvider RedirectStatisticProvider
}

func New(
	lp LinkGetter,
	aliasedLinkInserter AliasedLinkInserter,
	aliasGenerator AliasGenerator,
	redirectStatisticProvider RedirectStatisticProvider,
) *Link {
	return &Link{
		linkGetter:                lp,
		aliasedLinkInserter:       aliasedLinkInserter,
		aliasGenerator:            aliasGenerator,
		redirectStatisticProvider: redirectStatisticProvider,
	}
}
