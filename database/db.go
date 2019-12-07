package database

import (
	"database/sql"

	"db-forum/models"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var clearDB = `DELETE FROM users; DELETE FROM forum; DELETE FROM thread; DELETE FROM post; DELETE FROM voice;`

type DB struct {
	pg *sql.DB

	CreateUserStmt        *sql.Stmt
	GetUserStmt           *sql.Stmt
	GetUserByUsernameStmt *sql.Stmt
	UpdateUserStmt        *sql.Stmt

	CreateForumStmt             *sql.Stmt
	GetForumStmt                *sql.Stmt
	GetForumThreadsStmt         *sql.Stmt
	GetForumThreadsWithTimeStmt *sql.Stmt

	CreateThreadStmt    *sql.Stmt
	GetThreadStmt       *sql.Stmt
	GetThreadByIDStmt   *sql.Stmt
	GetThreadBySlugStmt *sql.Stmt

	CreatePostStmt  *sql.Stmt
	GetPostByIDStmt *sql.Stmt

	GetPrevVoteThreadStmt *sql.Stmt
	CreatVoteThreadStmt   *sql.Stmt
	UpdateVoteThreadStmt  *sql.Stmt
	BigInsert             *sql.Stmt
}

var (
	db           *DB
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)

func GetStatus() *models.Status {
	var status models.Status
	db.pg.QueryRow(`SELECT count(*) FROM users;`).Scan(&status.User)
	db.pg.QueryRow(`SELECT count(*) FROM thread;`).Scan(&status.Thread)
	db.pg.QueryRow(`SELECT count(*) FROM post;`).Scan(&status.Post)
	db.pg.QueryRow(`SELECT count(*) FROM forum;`).Scan(&status.Forum)
	return &status
}

func InitDB(DSN string) error {
	var err error
	var newDB DB
	newDB.pg, err = sql.Open("postgres", DSN)
	if err != nil {
		return errors.Wrap(err, "can't open database")
	}
	if err = newDB.pg.Ping(); err != nil {
		return errors.Wrap(err, "can't connect to database")
	}
	db = &newDB
	if err = prepareQueries(); err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	if _, err = db.pg.Exec(clearDB); err != nil {
		return errors.Wrap(err, "can't clear db")
	}
	return nil
}

func prepareQueries() error {
	var err error
	prepare := make(map[string]**sql.Stmt)
	prepare[createUser] = &db.CreateUserStmt
	prepare[getUser] = &db.GetUserStmt
	prepare[updateUser] = &db.UpdateUserStmt
	prepare[getUserByUsername] = &db.GetUserByUsernameStmt

	prepare[createForum] = &db.CreateForumStmt
	prepare[getForum] = &db.GetForumStmt
	prepare[getForumThreadsWithTime] = &db.GetForumThreadsWithTimeStmt
	prepare[getForumThreads] = &db.GetForumThreadsStmt

	prepare[createThread] = &db.CreateThreadStmt
	prepare[getThread] = &db.GetThreadStmt
	prepare[getThreadByID] = &db.GetThreadByIDStmt
	prepare[getThreadBySlug] = &db.GetThreadBySlugStmt

	prepare[createPost] = &db.CreatePostStmt
	prepare[getPostByID] = &db.GetPostByIDStmt

	prepare[bigInsert] = &db.BigInsert
	prepare[updateVoteThread] = &db.UpdateVoteThreadStmt
	prepare[createVoteThread] = &db.CreatVoteThreadStmt
	for query, value := range prepare {
		*value, err = db.pg.Prepare(query)
		if err != nil {
			return errors.Wrap(err, "can't prepare query "+query)
		}
	}
	return err
}

func ClearTable() {
	db.pg.Exec(clearDB)
}
