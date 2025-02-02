package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/nanoteck137/beldum/core"
	"github.com/nanoteck137/beldum/database"
	"github.com/nanoteck137/beldum/tools/utils"
	"github.com/nanoteck137/beldum/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id    string `json:"id"`
	Title string `json:"name"`

	BoardId   string `json:"boardId"`
	BoardName string `json:"boardName"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type Board struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Items []Task `json:"items"`
}

type GetProjects struct {
	Projects []Project `json:"projects"`
}

type GetProjectById struct {
	Project
}

type CreateProject struct {
	Id string `json:"id"`
}

type CreateProjectBody struct {
	Name string `json:"name"`
}

func (b *CreateProjectBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateProjectBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type GetProjectBoards struct {
	Boards []Board `json:"boards"`
}

type ShallowBoard struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Order int64  `json:"order"`
}

type GetAllProjectBoards struct {
	Boards       []ShallowBoard `json:"boards"`
	HiddenBoards []ShallowBoard `json:"hiddenBoards"`
}

type GetProjectTasks struct {
	Tasks []Task `json:"tasks"`
}

type CreateTask struct {
	Id string `json:"id"`
}

type CreateTaskBody struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`

	BoardId string `json:"boardId"`
}

func TransformTags(arr []string) []string {
	for i, v := range arr {
		arr[i] = utils.Slug(v)
	}

	return arr
}

func (b *CreateTaskBody) Transform() {
	b.Title = transform.String(b.Title)
	b.Tags = TransformTags(b.Tags)
}

func (b CreateTaskBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Title, validate.Required),
	)
}

type CreateBoard struct {
	Id string `json:"id"`
}

type CreateBoardBody struct {
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`

	ProjectId string `json:"projectId"`
}

func (b *CreateBoardBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateBoardBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type EditBoardBody struct {
	Name  *string `json:"name,omitempty"`
	Order *int64  `json:"order,omitempty"`
}

func (b *EditBoardBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditBoardBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

func InstallTaskHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "CreateProject",
			Method:       http.MethodPost,
			Path:         "/projects",
			ResponseType: CreateProject{},
			BodyType:     CreateProjectBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateProjectBody](c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().CreateProject(ctx, database.CreateProjectParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				_, err = app.DB().CreateBoard(ctx, database.CreateBoardParams{
					Name:      "Backlog",
					ProjectId: project.Id,
					OrderNumber: sql.NullInt64{
						Int64: 0,
						Valid: true,
					},
				})
				if err != nil {
					return nil, err
				}

				_, err = app.DB().CreateBoard(ctx, database.CreateBoardParams{
					Name:      "Work in progress",
					ProjectId: project.Id,
					OrderNumber: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
				})
				if err != nil {
					return nil, err
				}

				_, err = app.DB().CreateBoard(ctx, database.CreateBoardParams{
					Name:      "Done",
					ProjectId: project.Id,
					OrderNumber: sql.NullInt64{
						Int64: 2,
						Valid: true,
					},
				})
				if err != nil {
					return nil, err
				}

				return CreateProject{
					Id: project.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetUserProjects",
			Method:       http.MethodGet,
			Path:         "/projects",
			ResponseType: GetProjects{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				projects, err := app.DB().GetProjectsByUser(ctx, user.Id)
				if err != nil {
					return nil, err
				}

				res := GetProjects{
					Projects: make([]Project, len(projects)),
				}

				for i, project := range projects {
					res.Projects[i] = Project{
						Id:   project.Id,
						Name: project.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetProjectById",
			Method:       http.MethodGet,
			Path:         "/projects/:projectId",
			ResponseType: GetProjectById{},
			Errors:       []pyrin.ErrorType{ErrTypeProjectNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				projectId := c.Param("projectId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, projectId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ProjectNotFound()
					}

					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				return GetProjectById{
					Project: Project{
						Id:   project.Id,
						Name: project.Name,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateBoard",
			Method:       http.MethodPost,
			Path:         "/boards",
			ResponseType: CreateBoard{},
			BodyType:     CreateBoardBody{},
			Errors:       []pyrin.ErrorType{ErrTypeProjectNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				body, err := pyrin.Body[CreateBoardBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, body.ProjectId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ProjectNotFound()
					}

					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				orderNumber := sql.NullInt64{}

				if !body.Hidden {
					boards, err := app.DB().GetBoardsByProject(ctx, project.Id, false)
					if err != nil {
						return nil, err
					}

					if len(boards) > 0 {
						lastBoard := boards[len(boards)-1]
						orderNumber = sql.NullInt64{
							Int64: lastBoard.OrderNumber.Int64 + 1,
							Valid: true,
						}
					} else {
						orderNumber = sql.NullInt64{
							Int64: 0,
							Valid: true,
						}
					}
				}

				board, err := app.DB().CreateBoard(ctx, database.CreateBoardParams{
					Name:        body.Name,
					ProjectId:   project.Id,
					OrderNumber: orderNumber,
				})
				if err != nil {
					return nil, err
				}

				return CreateBoard{
					Id: board.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetAllProjectBoards",
			Method:       http.MethodGet,
			Path:         "/projects/:projectId/boards/all",
			ResponseType: GetAllProjectBoards{},
			Errors:       []pyrin.ErrorType{ErrTypeProjectNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				projectId := c.Param("projectId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, projectId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ProjectNotFound()
					}

					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				boards, err := app.DB().GetBoardsByProject(ctx, project.Id, false)
				if err != nil {
					return nil, err
				}

				hiddenBoards, err := app.DB().GetBoardsByProject(ctx, project.Id, true)
				if err != nil {
					return nil, err
				}

				res := GetAllProjectBoards{
					Boards:       make([]ShallowBoard, len(boards)),
					HiddenBoards: make([]ShallowBoard, len(hiddenBoards)),
				}

				for i, board := range boards {
					res.Boards[i] = ShallowBoard{
						Id:    board.Id,
						Name:  board.Name,
						Order: board.OrderNumber.Int64,
					}
				}

				for i, board := range hiddenBoards {
					res.HiddenBoards[i] = ShallowBoard{
						Id:    board.Id,
						Name:  board.Name,
						Order: board.OrderNumber.Int64,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetProjectBoards",
			Method:       http.MethodGet,
			Path:         "/projects/:projectId/boards",
			ResponseType: GetProjectBoards{},
			Errors:       []pyrin.ErrorType{ErrTypeProjectNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				projectId := c.Param("projectId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, projectId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ProjectNotFound()
					}

					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				boards, err := app.DB().GetBoardsByProject(ctx, project.Id, false)
				if err != nil {
					return nil, err
				}

				res := GetProjectBoards{
					Boards: make([]Board, len(boards)),
				}

				for i, board := range boards {
					dbItems, err := app.DB().GetTasksByBoard(ctx, board.Id)
					if err != nil {
						return nil, err
					}

					items := make([]Task, len(dbItems))

					for i, item := range dbItems {
						items[i] = Task{
							Id:        item.Id,
							Title:     item.Title,
							BoardId:   item.BoardId,
							BoardName: item.BoardName,
							Tags:      utils.SplitString(item.Tags.String),
							Created:   item.Created,
							Updated:   item.Updated,
						}
					}

					res.Boards[i] = Board{
						Id:    board.Id,
						Name:  board.Name,
						Items: items,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetProjectTasks",
			Method:       http.MethodGet,
			Path:         "/projects/:projectId/tasks",
			ResponseType: GetProjectTasks{},
			Errors:       []pyrin.ErrorType{ErrTypeProjectNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				projectId := c.Param("projectId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, projectId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ProjectNotFound()
					}

					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				tasks, err := app.DB().GetTasksByProject(ctx, project.Id)
				if err != nil {
					return nil, err
				}

				res := GetProjectTasks{
					Tasks: make([]Task, len(tasks)),
				}

				for i, task := range tasks {
					res.Tasks[i] = Task{
						Id:        task.Id,
						Title:     task.Title,
						BoardId:   task.BoardId,
						BoardName: task.BoardName,
						Tags:      utils.SplitString(task.Tags.String),
						Created:   task.Created,
						Updated:   task.Updated,
					}
				}

				return res, nil
			},
		},

		// TODO(patrik): Move
		pyrin.ApiHandler{
			Name:     "EditBoard",
			Method:   http.MethodPatch,
			Path:     "/boards/:boardId",
			BodyType: EditBoardBody{},
			Errors:   []pyrin.ErrorType{ErrTypeBoardNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				boardId := c.Param("boardId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[EditBoardBody](c)
				if err != nil {
					return nil, err
				}

				board, err := app.DB().GetBoardById(ctx, boardId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, BoardNotFound()
					}

					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, board.ProjectId)
				if err != nil {
					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, BoardNotFound()
				}

				changes := database.BoardChanges{}
				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != board.Name,
					}
				}

				if body.Order != nil {
					changes.OrderNumber = types.Change[sql.NullInt64]{
						Value: sql.NullInt64{
							Int64: *body.Order,
							Valid: *body.Order != 0,
						},
						Changed: true,
					}
				}

				// if body.Hidden != nil {
				// 	if *body.Hidden {
				// 	} else {
				// 		boards, err := app.DB().GetBoardsByProject(ctx, project.Id, false)
				// 		if err != nil {
				// 			return nil, err
				// 		}
				//
				// 		if len(boards) > 0 {
				// 			lastBoard := boards[len(boards)-1]
				// 			orderNumber = lastBoard.OrderNumber.Int64 + 1
				// 		}
				// 	}
				// }

				// var orderNumber int64 = 0
				//
				// if body.Order != nil {
				// 	orderNumber = int64(*body.Order)
				// } else {
				// 	boards, err := app.DB().GetBoardsByProject(ctx, project.Id, false)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				//
				// 	if len(boards) > 0 {
				// 		lastBoard := boards[len(boards)-1]
				// 		orderNumber = lastBoard.OrderNumber.Int64 + 1
				// 	}
				// }
				//
				// changes.OrderNumber = types.Change[sql.NullInt64]{
				// 	Value: sql.NullInt64{
				// 		Int64: orderNumber,
				// 		Valid: true,
				// 	},
				// 	Changed: true,
				// }

				err = app.DB().UpdateBoard(ctx, board.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// TODO(patrik): Move
		pyrin.ApiHandler{
			Name:         "CreateTask",
			Method:       http.MethodPost,
			Path:         "/tasks",
			ResponseType: CreateTask{},
			BodyType:     CreateTaskBody{},
			Errors:       []pyrin.ErrorType{ErrTypeBoardNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateTaskBody](c)
				if err != nil {
					return nil, err
				}

				board, err := app.DB().GetBoardById(ctx, body.BoardId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, BoardNotFound()
					}

					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, board.ProjectId)
				if err != nil {
					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, ProjectNotFound()
				}

				task, err := app.DB().CreateTask(ctx, database.CreateTaskParams{
					Title:     body.Title,
					ProjectId: project.Id,
					BoardId:   board.Id,
				})
				if err != nil {
					return nil, err
				}

				for _, tag := range body.Tags {
					err := app.DB().CreateTag(ctx, project.Id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}

					err = app.DB().AddTaskTag(ctx, task.Id, project.Id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				return CreateTask{
					Id: task.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteTask",
			Method: http.MethodDelete,
			Path:   "/tasks/:taskId",
			Errors: []pyrin.ErrorType{ErrTypeTaskNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taskId := c.Param("taskId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				task, err := app.DB().GetTaskById(ctx, taskId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaskNotFound()
					}

					return nil, err
				}

				project, err := app.DB().GetProjectById(ctx, task.ProjectId)
				if err != nil {
					return nil, err
				}

				if project.OwnerId != user.Id {
					return nil, TaskNotFound()
				}

				err = app.DB().DeleteTask(ctx, task.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// TODO(patrik): Move
		// TODO(patrik): Fix errors
		pyrin.ApiHandler{
			Name:   "MoveTask",
			Method: http.MethodPost,
			Path:   "/tasks/:taskId/move/:boardId",
			Errors: []pyrin.ErrorType{ErrTypeTaskNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				taskId := c.Param("taskId")
				boardId := c.Param("boardId")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				task, err := app.DB().GetTaskById(ctx, taskId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TaskNotFound()
					}

					return nil, err
				}

				checkBoard := func(boardId string) (database.Board, database.Project, error) {
					board, err := app.DB().GetBoardById(ctx, boardId)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							return database.Board{}, database.Project{}, BoardNotFound()
						}

						return database.Board{}, database.Project{}, err
					}

					project, err := app.DB().GetProjectById(ctx, board.ProjectId)
					if err != nil {
						return database.Board{}, database.Project{}, err
					}

					if project.OwnerId != user.Id {
						return database.Board{}, database.Project{}, ProjectNotFound()
					}

					return board, project, nil
				}

				dstBoard, dstProject, err := checkBoard(boardId)
				if err != nil {
					return nil, err
				}

				srcBoard, srcProject, err := checkBoard(task.BoardId)
				if err != nil {
					return nil, err
				}

				if dstProject.Id != srcProject.Id {
					return nil, ProjectNotFound()
				}

				if task.ProjectId != dstProject.Id {
					// TODO(patrik): Better error
					return nil, errors.New("Project not matching")
				}

				err = app.DB().UpdateTask(ctx, task.Id, database.TaskChanges{
					BoardId: types.Change[string]{
						Value:   dstBoard.Id,
						Changed: dstBoard.Id != srcBoard.Id,
					},
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
