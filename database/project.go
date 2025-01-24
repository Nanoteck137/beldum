package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/beldum/tools/utils"
	"github.com/nanoteck137/beldum/types"
)

type Project struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Name string `db:"name"`

	OwnerId string `db:"owner_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func ProjectQuery() *goqu.SelectDataset {
	query := dialect.From("projects").
		Select(
			"projects.rowid",

			"projects.id",
			"projects.name",

			"projects.owner_id",

			"projects.created",
			"projects.updated",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllProjects(ctx context.Context) ([]Project, error) {
	query := ProjectQuery()

	var items []Project
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetProjectById(ctx context.Context, id string) (Project, error) {
	query := ProjectQuery().
		Where(goqu.I("projects.id").Eq(id))

	var item Project
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Project{}, ErrItemNotFound
		}

		return Project{}, err
	}

	return item, nil
}

func (db *Database) GetProjectsByUser(ctx context.Context, userId string) ([]Project, error) {
	query := ProjectQuery().
		Where(goqu.I("projects.owner_id").Eq(userId))

	var items []Project
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreateProjectParams struct {
	Id   string
	Name string

	OwnerId string

	Created int64
	Updated int64
}

func (db *Database) CreateProject(ctx context.Context, params CreateProjectParams) (Project, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateProjectId()
	}

	query := dialect.Insert("projects").
		Rows(goqu.Record{
			"id":   id,
			"name": params.Name,

			"owner_id": params.OwnerId,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"projects.id",
			"projects.name",

			"projects.owner_id",

			"projects.created",
			"projects.updated",
		).
		Prepared(true)

	var item Project
	err := db.Get(&item, query)
	if err != nil {
		return Project{}, err
	}

	return item, nil
}

func (db *Database) DeleteProject(ctx context.Context, id string) error {
	query := dialect.Delete("projects").
		Prepared(true).
		Where(goqu.I("projects.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type ProjectChanges struct {
	Name types.Change[string]

	OwnerId types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateProject(ctx context.Context, id string, changes ProjectChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "owner_id", changes.OwnerId)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("projects").
		Set(record).
		Where(goqu.I("projects.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}
