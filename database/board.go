package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/kr/pretty"
	"github.com/nanoteck137/beldum/tools/utils"
	"github.com/nanoteck137/beldum/types"
)

type Board struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Name string `db:"name"`

	ProjectId string `db:"project_id"`

	OrderNumber int64 `db:"order_number"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func BoardQuery() *goqu.SelectDataset {
	query := dialect.From("boards").
		Select(
			"boards.rowid",

			"boards.id",
			"boards.name",

			"boards.project_id",

			"boards.order_number",

			"boards.created",
			"boards.updated",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllBoards(ctx context.Context) ([]Board, error) {
	query := BoardQuery()

	var items []Board
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetBoardById(ctx context.Context, id string) (Board, error) {
	query := BoardQuery().
		Where(goqu.I("boards.id").Eq(id))

	var item Board
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Board{}, ErrItemNotFound
		}

		return Board{}, err
	}

	return item, nil
}

func (db *Database) GetBoardsByProject(ctx context.Context, projectId string) ([]Board, error) {
	query := BoardQuery().
		Where(goqu.I("boards.project_id").Eq(projectId))

	var items []Board
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreateBoardParams struct {
	Id   string
	Name string

	ProjectId string

	OrderNumber int64

	Created int64
	Updated int64
}

func (db *Database) CreateBoard(ctx context.Context, params CreateBoardParams) (Board, error) {
	pretty.Println(params)

	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateBoardId()
	}

	query := dialect.Insert("boards").
		Rows(goqu.Record{
			"id":   id,
			"name": params.Name,

			"project_id": params.ProjectId,

			"order_number": params.OrderNumber,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"boards.id",
			"boards.name",

			"boards.project_id",

			"boards.order_number",

			"boards.created",
			"boards.updated",
		).
		Prepared(true)

	var item Board
	err := db.Get(&item, query)
	if err != nil {
		return Board{}, err
	}

	return item, nil
}

func (db *Database) DeleteBoard(ctx context.Context, id string) error {
	query := dialect.Delete("boards").
		Prepared(true).
		Where(goqu.I("boards.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type BoardChanges struct {
	Name types.Change[string]

	ProjectId types.Change[string]

	OrderNumber types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateBoard(ctx context.Context, id string, changes BoardChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "project_id", changes.ProjectId)

	addToRecord(record, "order_number", changes.OrderNumber)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("boards").
		Set(record).
		Where(goqu.I("boards.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
