package service

import (
	"context"
	"errors"
	"github.com/ananaslegend/short-link/internal/save/repository"
	"github.com/ananaslegend/short-link/pkg/clog"
	"github.com/ananaslegend/short-link/pkg/shortner"
	"github.com/google/uuid"
	"strings"
)

type InsertLinkRepo interface {
	InsertLink(ctx context.Context, link, alias string) error
}

type Service struct {
	repo InsertLinkRepo
}

func New(lp InsertLinkRepo) *Service {
	return &Service{
		repo: lp,
	}
}

func (s Service) AddLink(ctx context.Context, link, alias string) (string, error) {
	if len(alias) == 0 {
		alias = shortner.MakeShorter(uuid.New().ID())
		ctx = clog.WithBool(ctx, "auto_alias", true)
	}

	ctx = clog.WithString(ctx, "alias", alias)

	if !strings.Contains(link, "http") {
		link = "http://" + link
	}

	ctx = clog.WithString(ctx, "link", link)

	if err := s.repo.InsertLink(ctx, link, alias); err != nil {
		clog.Ctx(ctx).Error("insert link", clog.ErrorMsg(err))

		if errors.Is(err, repository.ErrAliasAlreadyExists) {
			return "", ErrAliasAlreadyExists
		}

		return "", err
	}

	return alias, nil
}
