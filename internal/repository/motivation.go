package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Motivation struct {
	Id      int    `db:"id"`
	Content string `db:"content"`
	Author  string `db:"author"`
}

func (r *Repository) GetRandomMotivation(ctx context.Context, dbpool *pgxpool.Pool) (m Motivation, err error) {
	row := dbpool.QueryRow(ctx, `select * from motivations order by random() limit 1`)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&m.Id, &m.Content, &m.Author)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}
