package database

import (
	"db-forum/models"

	"database/sql"

	"github.com/pkg/errors"
)

var createForum = `INSERT INTO forum (title, author, slug) VALUES ($1, $2, $3);`
var getForum = `SELECT title, author, slug, posts, threads FROM forum WHERE slug = $1 LIMIT 1;`
var getForumThreads = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE forum = $1;`
var getForumThreadsWithTime = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE forum = $1 AND created <= $2;`
var getThreadInfo = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE forum = $1`

func GetForumThreads(forum string, since string, order string, limit int) (*[]models.Thread, error) {
	threads := make([]models.Thread, 0)
	var rows *sql.Rows
	var err error
	query := getThreadInfo
	if since != "" {
		switch order {
		case "DESC":
			query += " AND created <= $2 ORDER BY created DESC LIMIT $3;"
		case "ASC":
			query += " AND created >= $2 ORDER BY created LIMIT $3;"
		}
		rows, err = db.pg.Query(query, forum, since, limit)
	} else {
		switch order {
		case "DESC":
			query += " ORDER BY created DESC LIMIT $2;"
		case "ASC":
			query += " ORDER BY created LIMIT $2;"
		}
		rows, err = db.pg.Query(query, forum, limit)
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't select from thread")
	}
	for rows.Next() {
		var thread models.Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug); err != nil {
			return nil, errors.Wrap(err, "can't scan rows")
		}
		threads = append(threads, thread)
	}
	return &threads, nil
}

func GetForum(slug string) (*models.Forum, error) {
	var forum models.Forum
	if err := db.GetForumStmt.QueryRow(slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from db")
	}
	return &forum, nil
}

func CreateForum(forum *models.Forum) (*models.Forum, error) {
	_, err := db.CreateForumStmt.Exec(forum.Title, forum.User, forum.Slug)
	if err != nil {
		f, err := GetForum(forum.Slug)
		if err != nil {
			if err == ErrNotFound {
				return nil, errors.New("can't insert into db")
			}
			return nil, errors.Wrap(err, "can't get from forum")
		}
		return f, ErrDuplicate
	}
	return forum, nil
}

