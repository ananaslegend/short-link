package usecases

import (
	"context"
	"fmt"
)

type LinkSaver interface {
	AddLink(link, alias string) error
}

func AddLink(c context.Context, ls LinkSaver, link, alias string) error {
	const op = "usecases.link.AddLink"

	if err := ls.AddLink(link, alias); err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return err
	}

	return nil
}
