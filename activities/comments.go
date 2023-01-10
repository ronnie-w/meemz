package activities

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/meemz/authentication"
	"github.com/meemz/database"
)

type CommentStruct struct {
	CommentId string
	FileId    string
	Comment   string
}

type CommentResponse struct {
	Username           string
	ProfileImg         string
	CommentId          string
	CommentTime        string
	Comment            string
	CommentAction      bool
	CommentLikes       int
	CommentLikedByUser bool
	Replies            int
}

type Replies struct {
	Username           string
	ProfileImg         string
	CommentTime        string
	Comment            string
	CommentId          string
	CommentAction      bool
	CommentLikes       int
	CommentLikedByUser bool
}

type UpdateComplete struct {
	Updated bool
}

type Error struct {
	Err string
}

var db = database.Conn()

func TimeStamp(time_stamp string) string {
	t := time.Now().Format(time.RFC3339)
	t0, _ := time.Parse(time.RFC3339, t)
	t1, _ := time.Parse(time.RFC3339, time_stamp)

	raw_time := t0.Sub(t1).String()
	var parsed_time string

	time_parser := func(rt string) string {
		rt_float, _ := strconv.ParseFloat(rt, 32)

		strconv_float := func(num int) string {
			return strconv.Itoa(int(math.Round(rt_float / float64(num))))
		}

		if rt_float < 24 {
			return strconv.Itoa(int(rt_float)) + "h ago"
		} else if rt_float > 24 && rt_float < 731 {
			return strconv_float(24) + "d ago"
		} else if rt_float > 731 && rt_float < 8772 {
			return strconv_float(731) + " mon ago"
		} else {
			return strconv_float(8772) + "y ago"
		}
	}

	if !strings.Contains(raw_time, "m") {
		return strings.Trim(raw_time, "s") + "s ago"
	} else if !strings.Contains(raw_time, "h") {
		return strings.Split(raw_time, "m")[0] + " min ago"
	} else {
		parsed_time = time_parser(strings.Split(raw_time, "h")[0])
		return parsed_time
	}
}

func CommentForm(r *http.Request) *CommentStruct {
	r.ParseForm()
	details := new(CommentStruct)
	schema.NewDecoder().Decode(details, r.PostForm)

	return details
}

func CommentsUserDetails(uid string) (string, string) {
	var username string
	var profile_img string

	user_details := db.QueryRow("SELECT username, profileImg FROM users WHERE userId=?", uid)

	user_details.Scan(&username, &profile_img)

	return username, profile_img
}

func CommentLikes(comment_id string) int {
	var comment_likes int

	comment_likes_row := db.QueryRow("SELECT count(*) AS comment_likes FROM commentReplyLikes WHERE commentReplyId=?", comment_id)

	comment_likes_row.Scan(&comment_likes)

	return comment_likes
}

func CommentLikedByUser(comment_id, user_id string) bool {
	var comment_liked_by_user bool = false
	var comment_liked int

	commented := db.QueryRow("SELECT count(*) AS comment_liked FROM commentReplyLikes WHERE commentReplyId=? AND userId=?", comment_id, user_id)
	commented.Scan(&comment_liked)

	if comment_liked > 0 {
		comment_liked_by_user = true
		return comment_liked_by_user
	}

	return comment_liked_by_user
}

func FetchComments(rw http.ResponseWriter, r *http.Request) {
	var comment_response_arr []*CommentResponse

	user_id := authentication.FetchId(r)
	file_id := CommentForm(r).FileId

	comment_rows, err := db.Query("SELECT userId, commentId, commentTime, comment FROM comments WHERE fileId=?", file_id)
	if err != nil {
		log.Println(err)
	}

	for comment_rows.Next() {
		var comment_user_id string
		var username string
		var profile_img string
		var comment_id string
		var comment_time string
		var comment string
		var comment_action bool = false
		var comment_liked_by_user bool
		var comment_likes int
		var comment_replies int

		if err := comment_rows.Scan(&comment_user_id, &comment_id, &comment_time, &comment); err != nil {
			log.Println(err)
		}

		if comment_user_id == user_id {
			comment_action = true
		}

		username, profile_img = CommentsUserDetails(comment_user_id)
		comment_likes = CommentLikes(comment_id)
		comment_liked_by_user = CommentLikedByUser(comment_id, user_id)

		replies_num := db.QueryRow("SELECT count(*) AS comment_replies FROM replies WHERE commentId=?", comment_id)
		replies_num.Scan(&comment_replies)

		var comment_response *CommentResponse = &CommentResponse{
			Username:           username,
			ProfileImg:         profile_img,
			CommentId:          comment_id,
			CommentTime:        TimeStamp(comment_time),
			Comment:            authentication.TextParser(comment),
			CommentAction:      comment_action,
			CommentLikes:       comment_likes,
			CommentLikedByUser: comment_liked_by_user,
			Replies:            comment_replies,
		}

		comment_response_arr = append(comment_response_arr, comment_response)
	}

	defer comment_rows.Close()

	json.NewEncoder(rw).Encode(comment_response_arr)
}

func FetchReplies(rw http.ResponseWriter, r *http.Request) {
	var replies_response_arr []*Replies

	user_id := authentication.FetchId(r)
	comment_id := CommentForm(r).CommentId
	file_id := CommentForm(r).FileId

	replies_rows, err := db.Query("SELECT userId, replyId, replyTime, reply FROM replies WHERE fileId=? AND commentId=?", file_id, comment_id)
	if err != nil {
		log.Println(err)
	}

	for replies_rows.Next() {
		var reply_user_id string
		var reply_id string
		var username string
		var profile_img string
		var reply_time string
		var reply string
		var reply_action bool = false
		var reply_likes int
		var reply_liked_by_user bool

		if err := replies_rows.Scan(&reply_user_id, &reply_id, &reply_time, &reply); err != nil {
			log.Println(err)
		}

		if reply_user_id == user_id {
			reply_action = true
		}

		username, profile_img = CommentsUserDetails(reply_user_id)
		reply_likes = CommentLikes(reply_id)
		reply_liked_by_user = CommentLikedByUser(reply_id, user_id)

		var replies_response *Replies = &Replies{
			Username:           username,
			ProfileImg:         profile_img,
			CommentTime:        TimeStamp(reply_time),
			Comment:            authentication.TextParser(reply),
			CommentId:          reply_id,
			CommentAction:      reply_action,
			CommentLikes:       reply_likes,
			CommentLikedByUser: reply_liked_by_user,
		}

		replies_response_arr = append(replies_response_arr, replies_response)
	}

	defer replies_rows.Close()

	json.NewEncoder(rw).Encode(replies_response_arr)
}

func PostComment(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)

	if user_id != "" {
		comment := CommentForm(r).Comment
		file_id := CommentForm(r).FileId
		comment_id, _ := rand.Prime(rand.Reader, 70)

		comment_time := time.Now().Format(time.RFC3339)

		_ = db.QueryRow("INSERT INTO comments(userId, fileId, commentId, commentTime, comment) values(?, ?, ?, ?, ?)", user_id, file_id, comment_id.String(), comment_time, comment)

		json.NewEncoder(rw).Encode(UpdateComplete{true})

		notifications_row, _ := db.Query("SELECT userId FROM posts WHERE fileId=?", file_id)
		for notifications_row.Next() {
			var comment_to string
			notifications_row.Scan(&comment_to)
			Notify(r, "Commented on your post\n'"+comment+"'", []string{comment_to})
		}
	} else {
		json.NewEncoder(rw).Encode(Error{"You need to login to proceed"})
	}
}

func PostReply(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)

	if user_id != "" {
		reply := CommentForm(r).Comment
		file_id := CommentForm(r).FileId
		comment_id := CommentForm(r).CommentId
		reply_id, _ := rand.Prime(rand.Reader, 70)

		reply_time := time.Now().Format(time.RFC3339)

		_ = db.QueryRow("INSERT INTO replies(userId, fileId, commentId, replyId, replyTime, reply) values(?, ?, ?, ?, ?, ?)", user_id, file_id, comment_id, reply_id.String(), reply_time, reply)

		json.NewEncoder(rw).Encode(UpdateComplete{true})

		notification_rows, _ := db.Query("SELECT userId FROM comments WHERE commentId=?", comment_id)
		for notification_rows.Next() {
			var reply_to string
			notification_rows.Scan(&reply_to)
			Notify(r, "Replied to your comment\n'"+reply+"'", []string{reply_to})
		}
	} else {
		json.NewEncoder(rw).Encode(Error{"You need to login to proceed"})
	}
}

func PostCommentReplyLike(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)

	if user_id != "" {
		file_id := CommentForm(r).FileId
		comment_reply_id := CommentForm(r).CommentId

		_ = db.QueryRow("INSERT INTO commentReplyLikes(userId, fileId, commentReplyId) values(?, ?, ?)", user_id, file_id, comment_reply_id)

		json.NewEncoder(rw).Encode(UpdateComplete{true})
	} else {
		json.NewEncoder(rw).Encode(Error{"You need to login to proceed"})
	}
}

func DeleteComment(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)

	comment_id := CommentForm(r).CommentId

	delete_comment, _ := db.Query("DELETE FROM comments WHERE userId=? AND commentId=?", user_id, comment_id)
	defer delete_comment.Close()

	delete_replies, _ := db.Query("DELETE FROM replies WHERE commentId=?", comment_id)
	defer delete_replies.Close()

	delete_likes, _ := db.Query("DELETE FROM commentReplyLikes WHERE userId=? AND commentReplyId=?", user_id, comment_id)
	defer delete_likes.Close()

	json.NewEncoder(rw).Encode(UpdateComplete{true})
}

func DeleteReply(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)

	reply_id := CommentForm(r).CommentId

	delete_reply, _ := db.Query("DELETE FROM replies WHERE userId=? AND replyId=?", user_id, reply_id)
	defer delete_reply.Close()

	delete_likes, _ := db.Query("DELETE FROM commentReplyLikes WHERE userId=? AND commentReplyId=?", user_id, reply_id)
	defer delete_likes.Close()

	json.NewEncoder(rw).Encode(UpdateComplete{true})
}

func DeleteCommentReplyLike(rw http.ResponseWriter, r *http.Request) {
	user_id := authentication.FetchId(r)
	comment_reply_id := CommentForm(r).CommentId

	_ = db.QueryRow("DELETE FROM commentReplyLikes WHERE userId=? AND commentReplyId=?", user_id, comment_reply_id)

	json.NewEncoder(rw).Encode(UpdateComplete{true})
}
