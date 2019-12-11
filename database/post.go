package database

import (
	"database/sql"
	"db-forum/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var createPost = `INSERT INTO post (parent, author, message, forum, thread) 
VALUES ($1, $2, $3, $4, $5) RETURNING id, created;`
var postInsert = `INSERT INTO post (parent, message, thread, author, forum) VALUES `
var getPath = `SELECT path FROM post WHERE id = $1 AND thread = $2;`
var updatePostPath = `UPDATE post SET root = $2, path = $3 WHERE id = $1;`
var updateForumPostsCount = `UPDATE forum SET posts = posts + $2 WHERE slug = $1; `
var getPostByID = `SELECT id, parent, author, message, is_edited, forum, thread, created 
FROM post WHERE id = $1;`
var updatePost = `UPDATE post SET message = coalesce(coalesce(nullif($2, ''), message)), is_edited = $3 WHERE id = $1 RETURNING message, author, is_edited, thread, created, forum;`
var bigInsert = `INSERT INTO post (parent, message, thread, author, forum) values ($1, $2, $3, $4, $5),($6, $7, $8, $9, $10),($11, $12, $13, $14, $15),($16, $17, $18, $19, $20),($21, $22, $23, $24, $25),($26, $27, $28, $29, $30),($31, $32, $33, $34, $35),($36, $37, $38, $39, $40),($41, $42, $43, $44, $45),($46, $47, $48, $49, $50),($51, $52, $53, $54, $55),($56, $57, $58, $59, $60),($61, $62, $63, $64, $65),($66, $67, $68, $69, $70),($71, $72, $73, $74, $75),($76, $77, $78, $79, $80),($81, $82, $83, $84, $85),($86, $87, $88, $89, $90),($91, $92, $93, $94, $95),($96, $97, $98, $99, $100),($101, $102, $103, $104, $105),($106, $107, $108, $109, $110),($111, $112, $113, $114, $115),($116, $117, $118, $119, $120),($121, $122, $123, $124, $125),($126, $127, $128, $129, $130),($131, $132, $133, $134, $135),($136, $137, $138, $139, $140),($141, $142, $143, $144, $145),($146, $147, $148, $149, $150),($151, $152, $153, $154, $155),($156, $157, $158, $159, $160),($161, $162, $163, $164, $165),($166, $167, $168, $169, $170),($171, $172, $173, $174, $175),($176, $177, $178, $179, $180),($181, $182, $183, $184, $185),($186, $187, $188, $189, $190),($191, $192, $193, $194, $195),($196, $197, $198, $199, $200),($201, $202, $203, $204, $205),($206, $207, $208, $209, $210),($211, $212, $213, $214, $215),($216, $217, $218, $219, $220),($221, $222, $223, $224, $225),($226, $227, $228, $229, $230),($231, $232, $233, $234, $235),($236, $237, $238, $239, $240),($241, $242, $243, $244, $245),($246, $247, $248, $249, $250),($251, $252, $253, $254, $255),($256, $257, $258, $259, $260),($261, $262, $263, $264, $265),($266, $267, $268, $269, $270),($271, $272, $273, $274, $275),($276, $277, $278, $279, $280),($281, $282, $283, $284, $285),($286, $287, $288, $289, $290),($291, $292, $293, $294, $295),($296, $297, $298, $299, $300),($301, $302, $303, $304, $305),($306, $307, $308, $309, $310),($311, $312, $313, $314, $315),($316, $317, $318, $319, $320),($321, $322, $323, $324, $325),($326, $327, $328, $329, $330),($331, $332, $333, $334, $335),($336, $337, $338, $339, $340),($341, $342, $343, $344, $345),($346, $347, $348, $349, $350),($351, $352, $353, $354, $355),($356, $357, $358, $359, $360),($361, $362, $363, $364, $365),($366, $367, $368, $369, $370),($371, $372, $373, $374, $375),($376, $377, $378, $379, $380),($381, $382, $383, $384, $385),($386, $387, $388, $389, $390),($391, $392, $393, $394, $395),($396, $397, $398, $399, $400),($401, $402, $403, $404, $405),($406, $407, $408, $409, $410),($411, $412, $413, $414, $415),($416, $417, $418, $419, $420),($421, $422, $423, $424, $425),($426, $427, $428, $429, $430),($431, $432, $433, $434, $435),($436, $437, $438, $439, $440),($441, $442, $443, $444, $445),($446, $447, $448, $449, $450),($451, $452, $453, $454, $455),($456, $457, $458, $459, $460),($461, $462, $463, $464, $465),($466, $467, $468, $469, $470),($471, $472, $473, $474, $475),($476, $477, $478, $479, $480),($481, $482, $483, $484, $485),($486, $487, $488, $489, $490),($491, $492, $493, $494, $495),($496, $497, $498, $499, $500) returning id, is_edited, created`
var postByThread = `SELECT id, parent, author, message, forum, thread, created FROM post WHERE thread = $1`
var parentQuery = `SELECT id, parent, author, message, forum, thread, created FROM post WHERE root IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 `

func GetPostByID(id int64) (*models.Post, error) {
	var post models.Post
	if err := db.GetPostByIDStmt.QueryRow(id).Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't select from post")
	}
	return &post, nil
}

func GetPostsTree(thread int32, limit string, since string, desc string) (*[]models.Post, error) {
	posts := make([]models.Post, 0)
	getPostTree := postByThread
	var rows *sql.Rows
	var err error
	if since != "" {
		if desc == "true" {
			getPostTree += ` AND path < (SELECT path FROM post WHERE id = $2 ) ORDER BY path DESC LIMIT $3;`
		} else {
			getPostTree += ` AND path > (SELECT path FROM post WHERE id = $2 ) ORDER BY path LIMIT $3;`
		}
		rows, err = db.pg.Query(getPostTree, thread, since, limit)
	} else {
		since = "0"
		if desc == "true" {
			getPostTree += ` ORDER BY path DESC LIMIT $2;`
		} else {
			getPostTree += ` ORDER BY path LIMIT $2;`
		}
		rows, err = db.pg.Query(getPostTree, thread, limit)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return &posts, nil
		}
		return nil, errors.Wrap(err, "can't select from posts")
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, errors.Wrap(err, "can't scan rows")
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func CreatePost(post *models.Post) (*models.Post, error) {
	newPost := *post
	if err := db.CreatePostStmt.QueryRow(post.Parent, post.Author, post.Message, post.Forum, post.Thread).Scan(&newPost.ID, &newPost.Created); err != nil {
		return nil, errors.Wrap(err, "can't insert into post")
	}
	return &newPost, nil
}

func CreatePosts(posts *[]models.Post, threadSlug string) (*[]models.Post, error) {
	tx, err := db.pg.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "can't start transaction")
	}
	resPosts := make([]models.Post, 0)
	var thread *models.Thread
	if govalidator.IsNumeric(threadSlug) {
		thread, err = GetThreadByID(threadSlug)
	} else {
		thread, err = GetThreadBySlug(threadSlug)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	queryEnd := " returning id, is_edited, created"
	var queryValues []string

	args := make([]interface{}, 0, len(*posts)*5)
	parents := make([]string, 0, len(*posts))

	if len(*posts) == 0 {
		return &resPosts, nil
	}

	if len(*posts) < 100 {
		for i, post := range *posts {
			author, err := GetUserByUsername(post.Author)
			if err != nil || author == nil {
				return nil, ErrNotFound
			}
			(*posts)[i].Author = author.Nickname
		}
	}

	if len(*posts) == 100 {
		for _, post := range *posts {
			if post.Parent != 0 {
				parents = append(parents, strconv.Itoa(int(post.Parent)))
			}
			args = append(args, post.Parent, post.Message, thread.ID, post.Author, thread.Forum)
		}
	} else {
		for _, post := range *posts {
			if post.Parent != 0 {
				parents = append(parents, strconv.Itoa(int(post.Parent)))
			}
			queryValues = append(queryValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5))
			args = append(args, post.Parent, post.Message, thread.ID, post.Author, thread.Forum)
		}
	}

	if len(parents) != 0 {
		rows, err := tx.Query(fmt.Sprint(`select thread from post where id in (`, strings.Join(parents, ","), ")"))
		hasP := false

		for rows.Next() {
			hasP = true

			var tId int32
			err = rows.Scan(&tId)
			if err != nil {
				log.Println(err)
			}

			if tId != thread.ID {
				return nil, ErrDuplicate
			}
		}

		if !hasP {
			return nil, ErrDuplicate
		}

	}

	query := postInsert
	query += strings.Join(queryValues, ",") + queryEnd
	var rows *sql.Rows
	if len(*posts) == 100 {
		rows, err = tx.Stmt(db.BigInsert).Query(args...)
	} else {
		rows, err = tx.Query(query, args...)
	}

	var par []string
	var nopar []string

	a := make(map[string]bool)
	for i, post := range *posts {
		if rows.Next() {
			err = rows.Scan(&((*posts)[i].ID), &((*posts)[i].IsEdited), &((*posts)[i].Created))
			if post.Parent != 0 {
				par = append(par, strconv.Itoa(int(post.ID)))
			} else {
				nopar = append(nopar, strconv.Itoa(int(post.ID)))
			}

			a["'"+post.Author+"'"] = true
			(*posts)[i].Forum = thread.Forum
			(*posts)[i].Thread = thread.ID
		}
	}
	rows.Close()

	auth := make([]string, 0, len(a))
	for key := range a {
		auth = append(auth, key)
	}

	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		log.Println("error on main query")
		log.Println(err)
		return nil, ErrDuplicate
	}
	tx.Exec(updateForumPostsCount, thread.Forum, len(*posts))

	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	for _, post := range *posts {
		var root int64
		sqlPath := make([]sql.NullInt64, 0)
		if post.Parent != 0 {
			if err = db.pg.QueryRow(getPath, post.Parent, thread.ID).Scan(pq.Array(&sqlPath)); err != nil {
				tx.Rollback()
				if err == sql.ErrNoRows {
					return nil, ErrDuplicate
				}
				return nil, errors.Wrap(err, "can't get path from parent post")
			}
			root = sqlPath[0].Int64
		} else {
			root = post.ID
		}
		sqlPath = append(sqlPath, sql.NullInt64{post.ID, true})
		updateStmt, err := db.pg.Prepare(updatePostPath)
		if err != nil {
			return nil, errors.Wrap(err, "can't prepare post path")
		}
		if _, err = updateStmt.Exec(post.ID, root, pq.Array(sqlPath)); err != nil {
			return nil, errors.Wrap(err, "can't update post path")
		}
	}

	return posts, nil
}

func GetPostsFlat(thread int32, limit string, since string, desc string) (*[]models.Post, error) {
	posts := make([]models.Post, 0)
	getPostsFlat := postByThread
	var rows *sql.Rows
	var err error
	if since != "" {
		if desc == "true" {
			getPostsFlat += " AND id < $2 ORDER BY id DESC LIMIT $3;"
		} else {
			getPostsFlat += " AND id > $2 ORDER BY id ASC LIMIT $3;"
		}
		rows, err = db.pg.Query(getPostsFlat, thread, since, limit)
	} else {
		if desc == "true" {
			getPostsFlat += " ORDER BY id DESC LIMIT $2;"
		} else {
			getPostsFlat += " ORDER BY id LIMIT $2;"
		}
		rows, err = db.pg.Query(getPostsFlat, thread, limit)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return &posts, nil
		}
		return nil, errors.Wrap(err, "can't select from posts")
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, errors.Wrap(err, "can't scan rows")
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func UpdatePost(post *models.Post) (*models.Post, error) {
	newPost := *post
	oldPost, err := GetPostByID(post.ID)
	if err != nil {
		return nil, err
	}
	if len(post.Message) == 0 {
		post.IsEdited = false
	} else {
		if oldPost != nil {
			if oldPost.Message == post.Message {
				post.IsEdited = false
			} else {
				post.IsEdited = true
			}
		}
	}

	if err := db.pg.QueryRow(updatePost, post.ID, post.Message, post.IsEdited).Scan(&newPost.Message, &newPost.Author, &newPost.IsEdited, &newPost.Thread, &newPost.Created, &newPost.Forum); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "can't update post")
	}
	return &newPost, nil
}

func GetPostsParentTree(thread int32, limit string, since string, desc string) (*[]models.Post, error) {
	posts := make([]models.Post, 0)
	getPostParentTree := parentQuery
	var rows *sql.Rows
	var err error
	if since != "" {
		if desc == "true" {
			getPostParentTree += ` AND root < (SELECT root FROM post WHERE id = $2 ) ORDER BY root DESC LIMIT $3)  ORDER BY root desc, path ;`
		} else {
			getPostParentTree += ` AND path > (SELECT path FROM post WHERE id = $2 ) ORDER BY id LIMIT $3) ORDER BY path;`
		}
		rows, err = db.pg.Query(getPostParentTree, thread, since, limit)
	} else {
		since = "0"
		if desc == "true" {
			getPostParentTree += `ORDER BY root DESC LIMIT $2) ORDER BY root DESC, path;`
		} else {
			getPostParentTree += `ORDER BY id LIMIT $2) ORDER BY path;`
		}
		rows, err = db.pg.Query(getPostParentTree, thread, limit)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return &posts, nil
		}
		return nil, errors.Wrap(err, "can't select from posts")
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, errors.Wrap(err, "can't scan rows")
		}
		posts = append(posts, post)
	}
	return &posts, nil
}


