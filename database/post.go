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

func CreatePost(post *models.Post) (*models.Post, error) {
	newPost := *post
	if err := db.CreatePostStmt.QueryRow(post.Parent, post.Author, post.Message, post.Forum, post.Thread).Scan(&newPost.ID, &newPost.Created); err != nil {
		return nil, errors.Wrap(err, "can't insert into post")
	}
	return &newPost, nil
}

func CreatePosts(posts *[]models.Post, threadSlug string) (*[]models.Post, error) {

}

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

func GetPostsFlat(thread int32, limit string, since string, desc string) (*[]models.Post, error) {

}

func GetPostsTree(thread int32, limit string, since string, desc string) (*[]models.Post, error) {

}

func GetPostsParentTree(thread int32, limit string, since string, desc string) (*[]models.Post, error) {

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
