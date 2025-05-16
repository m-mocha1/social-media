package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func makeSQL() *sql.DB {
	db, err := sql.Open("sqlite3", "./SQL/data.db")
	if err != nil {
		fmt.Println("opening ERR", err)
		os.Exit(1)
	}
	return db
}

func getAllnoti(name string) (allNoti []map[string]interface{}) {
rows, err := db.Query("SELECT id, nottype,sender,status,created_at, actionId FROM not_requests WHERE receiver = ?", name)
	if err != nil {
		fmt.Println("error from getallnoti", err)
		return
	}

	for rows.Next() {
		var num int
		var s, st, createat, t, text string
		var actID sql.NullString 
		if err := rows.Scan(&num, &t, &s, &st, &createat,&actID); err != nil {
			fmt.Println("error scanning row: ", err)
			continue
		}
		fmt.Println("type", t)
		var groupID string
		var groupName string
		if t == "followreq" {
			text = "Requested to follow you "
		} else if strings.Contains(t,"join group")  {
			parts := strings.Split(t[11:], "groupId:")
			text = fmt.Sprintf("requested to join %s group", parts[0])
			groupID = parts[1]
			groupName = parts[0]
			t = t[:10]
		} else if  strings.Contains(t,"invite group") {
			parts := strings.Split(t[13:], "groupId:")
			text = fmt.Sprintf("invite to join %s group", parts[0])
			groupID = parts[1]
			groupName = parts[0]
			t = t[:12]
		}else if t == "startFollow" {
				text = "accepted your follow request "
			}else if t == "likePost" {
				text = "liked Your post "
			}else if t == "comment" {
				text = "commented on your post "
			}

		sPfp := getUserPfp(db, s)
		allNoti = append(allNoti, map[string]interface{}{
			"id":         num,
			"type":       t,
			"text":       text,
			"sender":     s,
			"sender_pfp": sPfp,
			"status":     st,
			"time":       createat,
			"groupId":    groupID,
			"groupName":  groupName,
			"actId":  actID.String,
		})
	}
	return allNoti
}




func createfollow(db *sql.DB, fer, fing string) {
	_, err := db.Exec("INSERT INTO follows (follower_username, following_username) VALUES ($1, $2) ON CONFLICT DO NOTHING", fer, fing)
	if err != nil {
		fmt.Println("error insert follow", err)
		return
	}
	fmt.Println("follow inserted")
}

// -------------------------user-------------------------------------------------------
func CreateUser(db *sql.DB, pfp []byte, username, email, password, gender, f_name, l_name, age, aboutme string,
	w http.ResponseWriter, r *http.Request,
) error {
	query := `INSERT INTO users (pfp,username, email, password, gender, first_name, last_name, age, aboutme, public) VALUES ( ?,?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, pfp, username, email, password, gender, f_name, l_name, age, aboutme, "public")
	if err != nil {
		fmt.Println("err creating", err)
		return err
	}
	return nil
}

func getUser(userName string) (u uInfo) {
	q := `SELECT  pfp,username, age, first_name, last_name, email, gender,aboutme, public FROM users
	WHERE username = ? `

	lines, err := db.Query(q, userName)
	if err != nil {
		fmt.Println("err geting username info", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		var temPfp []byte
		err := lines.Scan(&temPfp, &u.Username, &u.Age, &u.Fname, &u.Lname, &u.Email, &u.Gender, &u.AboutMe, &u.Public)
		if err != nil {
			fmt.Println("err scan", err)
			return
		}
		u.Pfp = encodeToBase64(temPfp)
	}
	return u
}

func getUserProfile(userName, name string) (u uInfo) {
	q := `SELECT  pfp,username, age, first_name, last_name, email, gender,aboutme, public FROM users
	WHERE username = ? `

	lines, err := db.Query(q, userName)
	if err != nil {
		fmt.Println("err geting username info", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		var temPfp []byte
		err := lines.Scan(&temPfp, &u.Username, &u.Age, &u.Fname, &u.Lname, &u.Email, &u.Gender, &u.AboutMe, &u.Public)
		if err != nil {
			fmt.Println("err scan", err)
			return
		}
		u.Pfp = encodeToBase64(temPfp)
		u.Me = userName == name
		u.Following = checkIfFollwed(userName, name)
		u.Followers = getfollows(userName)
		fmt.Println(u.Following)
	}
	return u
}

//--------------------------create post -----------------------------------------------------

func createPost(db *sql.DB, postContent string, username string, postType string, img []byte, groupId string) (int64, error) {
	var query string
	var res sql.Result
	var err error
	groupID, err2 := strconv.Atoi(groupId)
	if err2 != nil {
		query = `INSERT INTO posts (username, content, postType, img) VALUES (?,?,?,?)`
		res, err = db.Exec(query, username, postContent, postType, img)
	} else {
		query = `INSERT INTO posts (username, content, postType, group_id, img) VALUES (?,?,?,?,?)`
		res, err = db.Exec(query, username, postContent, postType, groupID, img)
	}
	if err != nil {
		fmt.Println("Error inserting post:", err)
		return 0, nil
	}
	postid, err := res.LastInsertId()
	if err != nil {
		fmt.Println("err geting id", err)
		return 0, err
	}
	fmt.Println("post saved")
	return postid, nil
}

// --likes--------------------------
func likingPost(db *sql.DB, username string, post_id string) int {
	id, _ := strconv.Atoi(post_id)
	query := `INSERT INTO likes (post_id, username) VALUES (?,?)`
	_, err := db.Exec(query, id, username)
	if err != nil {
		fmt.Println("Error inserting like:", err)
		return 0
	}
	return getlikes(db, id)
}

func disLikeCom(db *sql.DB, username string, comId int) int {
	_, err := db.Exec("DELETE FROM likes WHERE com_id = ? AND username = ?", comId, username)
	if err != nil {
		fmt.Println("Error removing like:", err)
		return 0
	}
	return getComlikes(db, comId)
}

func likingCom(db *sql.DB, username string, comId int) int {
	query := `INSERT INTO likes (com_id, username) VALUES (?,?)`
	_, err := db.Exec(query, comId, username)
	if err != nil {
		fmt.Println("Error inserting like:", err)
		return 0
	}
	likes := getComlikes(db, comId)
	fmt.Println("likeCounter", likes)

	return likes
}

func getComlikes(db *sql.DB, comId int) int {
	query := `SELECT COUNT(*) FROM likes WHERE com_id = ?`
	var likeCount int
	err := db.QueryRow(query, comId).Scan(&likeCount)
	if err != nil {
		fmt.Println("Error fetching like count:", err)
		return 0
	}
	return likeCount
}

func getlikes(db *sql.DB, post_id int) int {
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ?`
	var likeCount int
	err := db.QueryRow(query, post_id).Scan(&likeCount)
	if err != nil {
		fmt.Println("Error fetching like count:", err)
		return 0
	}

	return likeCount
}

func disLikePost(db *sql.DB, username string, post_id string) int {
	id, _ := strconv.Atoi(post_id)
	_, err := db.Exec("DELETE FROM likes WHERE post_id = ? AND username = ?", id, username)
	if err != nil {
		fmt.Println("Error removing like:", err)
		return 0
	}
	return getlikes(db, id)
}

// ----------------------------
func createComment(db *sql.DB, commentContent string, username string, post_id int, img []byte) int64 {
	query := `INSERT INTO comments (post_id, username, content,img) VALUES (?,?,?,?)`
	res, err := db.Exec(query, post_id, username, commentContent, img)
	if err != nil {
		fmt.Println("Error inserting comment:", err)
		return 0
	}
	postid, err := res.LastInsertId()
	fmt.Println("***commenttid", postid)
	return postid
}

func getUserAllPost(db *sql.DB, username string, r *http.Request) ([]Opost, error) {
	var Allposts []Opost
	name, _ := getSession(r)
	me := username == name
	lock := checkIfUserPubOrPri(username)
	following := checkIfFollwedUser(name, username)
	fmt.Println("Following:", following)

	if me || lock == "public" || (lock == "alPri" && following) || (lock == "pri" && me) {
		fmt.Println("inside")
		lines, err := db.Query("SELECT id, username, content, img, postType, created_at FROM posts WHERE group_id IS NULL AND username = ? ORDER BY posts.created_at DESC", username)
		if err != nil {
			fmt.Println("Error getting posts")
			return Allposts, err
		}
		defer lines.Close()

		for lines.Next() {
			var imgyte []byte
			var time time.Time
			var post Opost
			lines.Scan(&post.PostID, &post.Username, &post.Text, &imgyte, &post.PostType, &time)
			if !me {
				if post.PostType == "Pri" && !me {
					continue
				}
				if post.PostType == "Fol" && !following {
					continue
				}

			}
			post.Date = time.Format("15:04 01-02")
			post.Likes = getlikes(db, atoi(post.PostID))
			post.Spic = encodeToBase64(imgyte)
			post.UserPfp = getUserPfp(db, post.Username)
			post.PostLiked = checkIfLiked(db, atoi(post.PostID), name)
			Allposts = append(Allposts, post)
		}
	}
	return Allposts, nil

}

func getFollowingAllPost(db *sql.DB, name string) ([]Opost, error) {
	var Allposts []Opost
	lines, err := db.Query("SELECT id,username,content,img, postType,created_at FROM posts WHERE group_id IS NULL ORDER BY posts.created_at DESC")
	if err != nil {
		fmt.Println("err geting posts")
		return Allposts, err
	}
	defer lines.Close()
	for lines.Next() {
		var imgyte []byte
		var time time.Time
		var post Opost
		lines.Scan(&post.PostID, &post.Username, &post.Text, &imgyte, &post.PostType, &time)
		if name != post.Username {
			if !checkIfFollwedUser(name, post.Username) || post.PostType == "Pri" {
				continue

			}
		}
		post.Date = time.Format("15:04 01-02")
		post.Likes = getlikes(db, atoi(post.PostID))
		post.Spic = encodeToBase64(imgyte)
		post.UserPfp = getUserPfp(db, post.Username)
		post.PostLiked = checkIfLiked(db, atoi(post.PostID), name)
		Allposts = append(Allposts, post)
	}
	return Allposts, nil
}
func getAllPost(db *sql.DB, name string) ([]Opost, error) {
	var Allposts []Opost
	lines, err := db.Query("SELECT id,username,content,img, postType,created_at FROM posts WHERE group_id IS NULL ORDER BY posts.created_at DESC")
	if err != nil {
		fmt.Println("err geting posts")
		return Allposts, err
	}
	defer lines.Close()
	for lines.Next() {
		var imgyte []byte
		var time time.Time
		var post Opost
		lines.Scan(&post.PostID, &post.Username, &post.Text, &imgyte, &post.PostType, &time)
		if post.PostType == "Pub" {
			post.Date = time.Format("15:04 01-02")
			post.Likes = getlikes(db, atoi(post.PostID))
			post.Spic = encodeToBase64(imgyte)
			post.UserPfp = getUserPfp(db, post.Username)
			post.PostLiked = checkIfLiked(db, atoi(post.PostID), name)
			Allposts = append(Allposts, post)
		}
	}
	return Allposts, nil
}
func getGroupPost(db *sql.DB, groupId int, name string) ([]Opost, error) {
	var Allposts []Opost
	lines, err := db.Query("SELECT id,username,content,img, postType,created_at FROM posts WHERE group_id = ? ORDER BY posts.created_at DESC", groupId)
	if err != nil {
		fmt.Println("err geting posts")
		return Allposts, err
	}
	defer lines.Close()
	for lines.Next() {
		var imgyte []byte
		var time time.Time
		var post Opost
		lines.Scan(&post.PostID, &post.Username, &post.Text, &imgyte, &post.PostType, &time)
		post.Date = time.Format("15:04 01-02")
		post.Likes = getlikes(db, atoi(post.PostID))
		post.Spic = encodeToBase64(imgyte)
		post.UserPfp = getUserPfp(db, post.Username)
		post.PostLiked = checkIfLiked(db, atoi(post.PostID), name)
		Allposts = append(Allposts, post)
	}
	return Allposts, nil
}

func getUserPfp(db *sql.DB, username string) (pfp string) {
	q := (`SELECT pfp FROM users WHERE username = ?`)
	lines, err := db.Query(q, username)
	if err != nil {
		fmt.Println("err geting pfp pic", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		var tem []byte
		err := lines.Scan(&tem)
		if err != nil {
			fmt.Println("err scanning pfp", err)
			return
		}
		pfp = encodeToBase64(tem)
	}
	return pfp
}

func checkIfLiked(db *sql.DB, post_id int, username string) bool {
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ? AND username = ?`
	var i int
	err := db.QueryRow(query, post_id, username).Scan(&i)
	if err != nil {
		fmt.Println("err checking likes", err)
		return false
	}
	return i > 0
}

func atoi(a string) (b int) {
	b, _ = strconv.Atoi(a)
	return b
}

func getOnePost(db *sql.DB, post_id, name string) (post Opost, err error) {
	id := atoi(post_id)
	q := `SELECT id,username,content,img,created_at FROM posts WHERE id = ?`
	lines, err := db.Query(q, id)
	if err != nil {
		return post, err
	}
	defer lines.Close()
	for lines.Next() {
		var imgyte []byte
		var time time.Time
		lines.Scan(&post.PostID, &post.Username, &post.Text, &imgyte, &time)
		post.Spic = encodeToBase64(imgyte)
		post.Date = time.Format("15:04 01-02")
		post.Likes = getlikes(db, atoi(post.PostID))
		post.UserPfp = getUserPfp(db, post.Username)
		post.PostLiked = checkIfLiked(db, atoi(post.PostID), name)
		post.Comment = getComment(atoi(post.PostID), name)
	}
	return post, nil
}

func getComment(postID int, name string) []Comment {
	query := `SELECT id,content,created_at,username,img
	FROM comments
	WHERE comments.post_id = ?
	ORDER BY comments.created_at DESC`
	rows, err := db.Query(query, postID)
	if err != nil {
		fmt.Println("Error fetching comments:", err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var imgyte []byte
		var comment Comment
		err := rows.Scan(&comment.CommentId, &comment.Text, &comment.Date, &comment.User, &imgyte)
		if err != nil {
			fmt.Println("Error scanning comment:", err)
			return nil
		}
		comment.ComPic = encodeToBase64(imgyte)
		comment.Likes = getComlikes(db, atoi(comment.CommentId))
		comment.Pfpimg = getUserPfp(db, comment.User)
		comment.Liked = checkIfComLiked(db, atoi(comment.CommentId), name)
		comments = append(comments, comment)

	}

	return comments
}

func checkIfComLiked(db *sql.DB, com_id int, username string) bool {
	query := `SELECT COUNT(*) FROM likes WHERE com_id = ? AND username = ?`
	var i int
	err := db.QueryRow(query, com_id, username).Scan(&i)
	if err != nil {
		fmt.Println("Error fetching like count:", err)
		return false
	}
	return i > 0
}


func getAllUser(db *sql.DB, name string) (users []uInfo) {
	q := `SELECT u.pfp, u.username 
        FROM users u
        JOIN (
            SELECT f1.following_username AS username
            FROM follows f1
            JOIN follows f2 
            ON f1.follower_username = f2.following_username 
            AND f1.following_username = f2.follower_username
            WHERE f1.follower_username = ? AND f1.following_username != ?
        ) mutuals 
        ON u.username = mutuals.username
        LEFT JOIN messages m 
        ON (u.username = m.sender OR u.username = m.receiver) 
        AND (m.sender = ? OR m.receiver = ?)
        GROUP BY u.username
        ORDER BY MAX(m.created_at) DESC;`
    lines, err := db.Query(q, name, name, name, name)
	if err != nil {
		fmt.Println("err geting users ", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		var user uInfo
		var tem []byte
		err := lines.Scan(&tem, &user.Username)
		if err != nil {
			fmt.Println("err scanning pfp", err)
			return
		}
		if user.Username == name {
			continue
		}
		user.Pfp = encodeToBase64(tem)
		fmt.Println("user.username", user.Username)
		user.Online = checkIfUserOnline(user.Username)
		users = append(users, user)
		fmt.Println("user", user.Online)
	}
	return users
}
func getAllGroup(db *sql.DB, name string) (users []uInfo) {
	fmt.Println("all")
	q := `SELECT g.id, g.pfp, g.name 
	FROM groups g
	LEFT JOIN messages m 
	ON (g.name = m.sender OR g.name = m.receiver) 
	AND (m.sender = ? OR m.receiver = ?)
	GROUP BY g.name
	ORDER BY MAX(m.created_at) DESC`
	lines, err := db.Query(q, name, name)
	if err != nil {
		fmt.Println("err geting users ", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		var id int
		var user uInfo
		var tem []byte
		err := lines.Scan(&id, &tem, &user.Username)
		if err != nil {
			fmt.Println("err scanning pfp", err)
			return
		}
		strId := strconv.Itoa(id)
		if !checkIfMember(name, strId) {
			continue
		}
		if user.Username == name {
			continue
		}
		user.Pfp = encodeToBase64(tem)
		fmt.Println("user.username", user.Username)
		user.Online = true
		users = append(users, user)
		fmt.Println("user", user.Online)
	}
	return users
}

func checkIfUserOnline(username string) bool {
	for _, users := range sessionsMap {
		if users == username {
			return true
		}
	}
	return false
}

func saveMsg(db *sql.DB, sender, receiver, msg string) {
	query := `INSERT INTO messages (sender, receiver, message) VALUES(?,?,?)`
	_, err := db.Exec(query, sender, receiver, msg)
	if err != nil {
		fmt.Println("err inserting msg", err)
		return
	}
}

func getAllMsg(db *sql.DB, sender, receiver string) ([]Msg, error) {
	var allMsg []Msg
	var total int
	limit := "10"
	offset := "0"

	err := db.QueryRow(`SELECT COUNT(*) FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)`, sender, receiver, receiver, sender).Scan(&total)
	if err != nil {
		fmt.Println("Error getting total messages:", err)
		return allMsg, err
	}
	if total > 10 {
		offset = fmt.Sprintf("%d", total-10) // Skip older messages
	} else {
		offset = "0"
	}
	q := `SELECT sender,receiver,message
	 FROM messages
	WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)
	ORDER BY created_at ASC
	 LIMIT ? OFFSET ?`
	lines, err := db.Query(q, sender, receiver, receiver, sender, limit, offset)
	if err != nil {
		fmt.Println("err geting msg")
		return allMsg, err
	}
	defer lines.Close()
	for lines.Next() {

		var msg Msg
		msg.SenderPfp = getUserPfp(db, sender)
		msg.RecNamePfp = getUserPfp(db, receiver)
		err := lines.Scan(&msg.Sender, &msg.RecName, &msg.Text)
		if err != nil {
			fmt.Println("Err scan in msg", err)
		}
		allMsg = append(allMsg, msg)
	}

	return allMsg, nil
}

func getGroupMsg(db *sql.DB, groupName string) ([]Msg, error) {
	var allMsg []Msg

	q := `SELECT sender, message FROM messages WHERE receiver = ? ORDER BY created_at ASC`
	lines, err := db.Query(q, groupName)
	if err != nil {
		fmt.Println("Error fetching group messages:", err)
		return allMsg, err
	}
	defer lines.Close()

	for lines.Next() {
		var msg Msg
		err := lines.Scan(&msg.Sender, &msg.Text)
		if err != nil {
			fmt.Println("Error scanning group message:", err)
			continue
		}

		// Debugging: Check the sender before fetching the profile picture
		fmt.Println("Fetching profile picture for:", msg.Sender)

		// Fetch sender profile picture
		msg.RecNamePfp = getUserPfp(db, msg.Sender)

		fmt.Println("Sender:", msg.Sender)

		allMsg = append(allMsg, msg)
	}

	return allMsg, nil
}



func getUnFromEmail(email string) (username string) {
	q := (`SELECT username FROM users WHERE email = ?`)
	lines, err := db.Query(q, email)
	if err != nil {
		fmt.Println("err geting pfp pic", err)
		return
	}
	defer lines.Close()
	for lines.Next() {
		err := lines.Scan(&username)
		if err != nil {
			fmt.Println("err scanning email to username", err)
			return
		}
	}
	return username
}

func checkIfUserPubOrPri(username string) string {
	var temp string
	q := (`SELECT public FROM users WHERE username = ?`)
	lines, err := db.Query(q, username)
	if err != nil {
		fmt.Println("err", err)
	}
	defer lines.Close()
	for lines.Next() {
		err := lines.Scan(&temp)
		if err != nil {
			fmt.Println("err", err)
		}
	}
	return temp
}

// checks how many followers
func getfollows(username string) (follow int) {
	query := `
		SELECT COUNT(*) FROM follows WHERE following_username = ?;`

	err := db.QueryRow(query, username).Scan(&follow)
	if err != nil {
		fmt.Println("error getfoll", err)
		return
	}
	fmt.Println("follow", follow)
	return follow
}

// check if me follwing the other user
func checkIfFollwed(name, username string) bool {
	// name is me, username is the other user
	var count int
	fmt.Println("sss")
	query := `SELECT COUNT(*) FROM follows WHERE follower_username = ? AND following_username = ?;`

	err := db.QueryRow(query, username, name).Scan(&count)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	fmt.Println("Follow count:", count)
	return count > 0
}

func checkIfFollwedUser(name, username string) bool {
	//this to avoid conflict with above func cuz pointer
	// name is me, username is the other user
	var co int
	query := `SELECT COUNT(*) FROM follows WHERE following_username = ? AND follower_username = ?;`

	err := db.QueryRow(query, username, name).Scan(&co)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	fmt.Println("Follow count:", co)
	return co > 0
}

func unfollowUser(name, username string) error {
	query := `DELETE FROM follows WHERE follower_username = ? AND following_username = ?;`

	_, err := db.Exec(query, name, username)
	if err != nil {
		fmt.Println("Error while unfollowing:", err)
		return err
	}
	return nil
}

func createGroup(db *sql.DB, pfp []byte, groupName string, groupDesc string, name string) (int64, error) {
	query := `INSERT INTO groups (pfp, name, description, creator) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, pfp, groupName, groupDesc, name)
	if err != nil {
		fmt.Println("Error inserting group:", err)
		return 0, err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error getting group ID:", err)
		return 0, err
	}
	query1 := `INSERT INTO members (group_id, user) VALUES (?, ?)`
	res1, err := db.Exec(query1, groupID, name)
	if err != nil {
		fmt.Println("Error inserting group:", err)
		return 0, err
	}
	fmt.Println(res1)
	fmt.Println("Group saved:", groupID)
	return groupID, nil
}

func getUserPic(userName string) string {
	q := `SELECT  pfp FROM users
	WHERE username = ? `

	lines, err := db.Query(q, userName)
	if err != nil {
		fmt.Println("err geting username info", err)
		return ""
	}
	defer lines.Close()
	var tempic []byte
	for lines.Next() {
		err := lines.Scan(&tempic)
		if err != nil {
			fmt.Println("err scan", err)
			return ""
		}

	}
	pic := encodeToBase64(tempic)
	return pic
}
func deletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("err deleting post sql", err)
		return err
	}
	return nil
}
func deleteCom(id int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("err deleting  com sql", err)
		return err
	}
	return nil
}
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID, _ := strconv.Atoi(r.FormValue("eventID"))
	query := "DELETE FROM events WHERE id = ?"
	_, err := db.Exec(query, eventID)
	if err != nil {
		fmt.Println("err", err)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func countNoti(w http.ResponseWriter, r *http.Request) {
	username, _ := getSession(r) 

	fmt.Println("Checking notifications for:", username)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM not_requests WHERE receiver = $1", username).Scan(&count)
	if err != nil {
		fmt.Println("Error getting notification count:", err)
		return
	}

	fmt.Println("Notification count:", count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
}