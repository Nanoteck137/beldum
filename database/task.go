package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/beldum/tools/utils"
	"github.com/nanoteck137/beldum/types"
)

type Task struct {
	RowId int `db:"rowid"`

	Id    string `db:"id"`
	Title string `db:"title"`

	ProjectId string `db:"project_id"`

	BoardId   string `db:"board_id"`
	BoardName string `db:"board_name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Tags sql.NullString `db:"tags"`
}

func TaskQuery() *goqu.SelectDataset {
	tags := dialect.From("tasks_tags").
		Select(
			goqu.I("tasks_tags.task_id").As("task_id"),
			goqu.Func("group_concat", goqu.I("tags.slug"), ",").As("tags"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(
				goqu.I("tasks_tags.project_id").Eq(goqu.I("tags.project_id")),
				goqu.I("tasks_tags.tag_slug").Eq(goqu.I("tags.slug")),
			),
		).
		GroupBy(goqu.I("tasks_tags.task_id"))

	query := dialect.From("tasks").
		Select(
			"tasks.rowid",

			"tasks.id",
			"tasks.title",

			"tasks.project_id",
			"tasks.board_id",

			"tasks.created",
			"tasks.updated",

			goqu.I("boards.name").As("board_name"),

			goqu.I("tags.tags").As("tags"),
		).
		Prepared(true).
		Join(
			goqu.I("boards"),
			goqu.On(goqu.I("tasks.board_id").Eq(goqu.I("boards.id"))),
		).
		LeftJoin(
			tags.As("tags"),
			goqu.On(goqu.I("tasks.id").Eq(goqu.I("tags.task_id"))),
		).
		Order(goqu.I("tasks.title").Asc())

	return query
}

func (db *Database) GetAllTasks(ctx context.Context) ([]Task, error) {
	query := TaskQuery()

	var items []Task
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTaskById(ctx context.Context, id string) (Task, error) {
	query := TaskQuery().
		Where(goqu.I("tasks.id").Eq(id))

	var item Task
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrItemNotFound
		}

		return Task{}, err
	}

	return item, nil
}

func (db *Database) GetTasksByBoard(ctx context.Context, boardId string) ([]Task, error) {
	query := TaskQuery().
		Where(goqu.I("tasks.board_id").Eq(boardId))

	var items []Task
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTasksByProject(ctx context.Context, projectId string) ([]Task, error) {
	query := TaskQuery().
		Where(goqu.I("tasks.project_id").Eq(projectId))

	var items []Task
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type CreateTaskParams struct {
	Id    string
	Title string

	ProjectId string
	BoardId   string

	Created int64
	Updated int64
}

func (db *Database) CreateTask(ctx context.Context, params CreateTaskParams) (Task, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateTaskId()
	}

	query := dialect.Insert("tasks").
		Rows(goqu.Record{
			"id":    id,
			"title": params.Title,

			"project_id": params.ProjectId,
			"board_id":   params.BoardId,

			"created": created,
			"updated": updated,
		}).
		Returning(
			"tasks.id",
			"tasks.title",

			"tasks.project_id",
			"tasks.board_id",

			"tasks.created",
			"tasks.updated",
		).
		Prepared(true)

	var item Task
	err := db.Get(&item, query)
	if err != nil {
		return Task{}, err
	}

	return item, nil
}

func (db *Database) DeleteTask(ctx context.Context, id string) error {
	query := dialect.Delete("tasks").
		Prepared(true).
		Where(goqu.I("tasks.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type TaskChanges struct {
	Title types.Change[string]

	ProjectId types.Change[string]
	BoardId   types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateTask(ctx context.Context, id string, changes TaskChanges) error {
	record := goqu.Record{}

	addToRecord(record, "title", changes.Title)

	addToRecord(record, "project_id", changes.ProjectId)
	addToRecord(record, "board_id", changes.BoardId)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("tasks").
		Set(record).
		Where(goqu.I("tasks.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) AddTaskTag(ctx context.Context, taskId, projectId, tagSlug string) error {
	ds := dialect.Insert("tasks_tags").
		Prepared(true).
		Rows(goqu.Record{
			"task_id":    taskId,
			"project_id": projectId,
			"tag_slug":   tagSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

// TODO(patrik): Generalize
func (db *Database) RemoveAllTaskTags(ctx context.Context, taskId string) error {
	query := dialect.Delete("tasks_tags").
		Prepared(true).
		Where(goqu.I("tasks_tags.task_id").Eq(taskId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
