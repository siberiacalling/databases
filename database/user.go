package database

import (
	"database/sql"
	"db-forum/models"

	"github.com/pkg/errors"
)

var createUser = `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;`
var getUserByUsername = `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 LIMIT 1;`
var getUser = `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2;`
var updateUser = `UPDATE users SET fullname = coalesce(coalesce(nullif($2, ''), fullname)), 
			email = coalesce(coalesce(nullif($3, ''), email)), 
			about = coalesce(coalesce(nullif($4, ''), about)) WHERE nickname = $1 RETURNING fullname, email, about;`
var getForumUsers = `SELECT nickname, fullname, about, email FROM users WHERE ( nickname IN (SELECT author FROM post WHERE forum = $1) 
					OR nickname IN (SELECT author FROM thread WHERE forum = $1) ) `

func CreateUser(user *models.User) (*[]models.User, error) {
	var users []models.User
	tx, err := db.pg.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare tx")
	}
	res, err := db.CreateUserStmt.Exec(user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "can't insert into users")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "can't get affected rows")
	}
	if ra == 0 {
		usr, err := GetUser(user.Nickname, user.Email)
		if err != nil {
			if err == ErrNotFound {
				tx.Rollback()
				return nil, errors.New("can't insert into db")
			}
			tx.Rollback()
			return nil, errors.Wrap(err, "can't get from users")
		}
		tx.Rollback()
		return usr, ErrDuplicate
	}
	users = append(users, *user)
	tx.Commit()
	return &users, nil
}

func GetUserByUsername(nickname string) (*models.User, error) {
	var user models.User
	if err := db.GetUserByUsernameStmt.QueryRow(nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
	}
	return &user, nil
}

func GetUser(nickname string, email string) (*[]models.User, error) {
	var users []models.User
	rows, err := db.GetUserStmt.Query(nickname, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from users")
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
			return nil, errors.Wrap(err, "can't scan row")
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}
	return &users, nil
}

func UpdateUser(user *models.User) (*[]models.User, error) {
	var users []models.User
	var newUser models.User
	err := db.UpdateUserStmt.QueryRow(user.Nickname, user.Fullname, user.Email, user.About).Scan(&newUser.Fullname, &newUser.Email, &newUser.About)
	if err != nil {
		usr, err := GetUser(user.Nickname, user.Email)
		if err != nil {
			if err == ErrNotFound {
				return nil, ErrNotFound
			}
			return nil, errors.Wrap(err, "can't get from users")
		}
		return usr, ErrDuplicate
	}
	newUser.Nickname = user.Nickname
	users = append(users, newUser)
	return &users, nil
}

func GetForumUsers(slug string, limit string, since string, desc string) ([]models.User, error) {
	query := getForumUsers
	users := make([]models.User, 0)
	var rows *sql.Rows
	var err error
	if since != "" {
		if desc == "true" {
			query += "AND nickname < $2 ORDER BY nickname DESC LIMIT $3;"
		} else {
			query += "AND nickname > $2 ORDER BY nickname LIMIT $3;"
		}
		rows, err = db.pg.Query(query, slug, since, limit)
	} else {
		if desc == "true" {
			query += "ORDER BY nickname DESC LIMIT $2;"
		} else {
			query += "ORDER BY nickname LIMIT $2;"
		}
		rows, err = db.pg.Query(query, slug, limit)
	}
	if err != nil {
		return users, errors.Wrap(err, "can't select users from forum")
	}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
			return users, errors.Wrap(err, "can't scan user")
		}
		users = append(users, user)
	}
	return users, nil
}
