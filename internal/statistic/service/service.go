package service

type Statistic struct {
	redirectHandler RedirectHandler
}

func NewStatistic(redirectHandler RedirectHandler) *Statistic {
	return &Statistic{redirectHandler: redirectHandler}
}
