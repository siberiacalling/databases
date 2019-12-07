package database

import (
	"db-forum/models"

	"database/sql"

	"github.com/pkg/errors"
)

var createThread = `INSERT INTO thread (title, author, forum, message, created, slug) VALUES ($1, $2, $3, $4, $5, $6) RETURNING slug, id;`
var updateForumCount = `UPDATE forum SET threads = threads + 1 WHERE slug = $1;`
var getThreadByID = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE id = $1;`
var getThreadBySlug = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE slug = $1;`
var createVoteThread = `INSERT INTO voice (nickname, vote, thread_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`
var getThread = `SELECT id, title, author, forum, message, votes, created, slug FROM thread WHERE id = $1 OR slug = $2;`
var updateVoteByID = `UPDATE voice SET prev_vote = vote, vote = $1 WHERE thread_id = $2 AND nickname = $3 RETURNING (vote - prev_vote);`
var updateThread = `UPDATE thread SET title = coalesce(coalesce(nullif($2, ''), title)),
			message = coalesce(coalesce(nullif($3, ''), message))
			WHERE id = $1 RETURNING title, message;`
var updateVoteThread = `UPDATE thread SET votes = votes + $1 WHERE id = $2 RETURNING votes;`

func CreateThread(thread *models.Thread) (*models.Thread, error) {

}

func GetThreadByID(id string) (*models.Thread, error) {
	var thread models.Thread
	if err := db.pg.QueryRow(getThreadByID, id).Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from thread")
	}
	return &thread, nil
}

func GetThreadByIDint32(id int32) (*models.Thread, error) {
	var thread models.Thread
	if err := db.pg.QueryRow(getThreadByID, id).Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from thread")
	}
	return &thread, nil
}

func GetThreadBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	if err := db.GetThreadBySlugStmt.QueryRow(slug).Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from thread")
	}
	return &thread, nil
}

func GetThread(id string, slug string) (*models.Thread, error) {
	var thread models.Thread
	if err := db.GetThreadStmt.QueryRow(id, slug).Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created, &thread.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from thread")
	}
	return &thread, nil
}

func VoteThread(vote *models.Vote) (newVote int32, err error) {
	tx, err := db.pg.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "can't start tx")
	}
	var diff int32
	if err := db.pg.QueryRow(updateVoteByID, vote.Voice, vote.ThreadId, vote.Nickname).Scan(&diff); err != nil {
		if err != sql.ErrNoRows {
			tx.Rollback()
			return 0, errors.Wrap(err, "can't update voice")
		}
		if _, err := db.CreatVoteThreadStmt.Exec(vote.Nickname, vote.Voice, vote.ThreadId); err != nil {
			tx.Rollback()
			return 0, errors.Wrap(err, "can't insert into voice")
		}
		diff = vote.Voice
	}
	if err := db.UpdateVoteThreadStmt.QueryRow(diff, vote.ThreadId).Scan(&newVote); err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "can't update thread")
	}
	tx.Commit()
	return newVote, nil
}

func UpdateThread(thread *models.Thread) (*models.Thread, error) {
	newThread := *thread
	updateThreadStmt, err := db.pg.Prepare(updateThread)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare query")
	}
	if err := updateThreadStmt.QueryRow(thread.ID, thread.Title, thread.Message).Scan(&newThread.Title, &newThread.Message); err != nil {
		return nil, errors.Wrap(err, "can't update thread")
	}
	return &newThread, nil
}
