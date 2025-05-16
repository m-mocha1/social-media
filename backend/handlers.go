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
)

var (
	userInfo uInfo
)

//	func home(w http.ResponseWriter, r *http.Request) {
//		path := "../frontend/dist" + r.URL.Path
//		if _, err := os.Stat(path); os.IsNotExist(err) {
//			http.ServeFile(w, r, "../frontend/dist/index.html")
//		} else {
//			fs.ServeHTTP(w, r)
//		}
//	}
func pri(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	q := `UPDATE users SET public = 'pri' WHERE username = ?`
	_, err := db.Exec(q, name)
	if err != nil {
		fmt.Println("err updating privecy", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func pub(w http.ResponseWriter, r *http.Request) {

	name, _ := getSession(r)
	q := `UPDATE users SET public = 'public' WHERE username = ?`
	_, err := db.Exec(q, name)
	if err != nil {
		fmt.Println("err updating privecy", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func almostPri(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	q := `UPDATE users SET public = 'alPri' WHERE username = ?`
	_, err := db.Exec(q, name)
	if err != nil {
		fmt.Println("err updating privecy", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getUserPost(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	username := r.FormValue("username")
	allPosts, err := getUserAllPost(db, username, r)
	userData := getUserProfile(username, name)
	if err != nil {
		fmt.Println("err geting posts", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": allPosts,
		"user":  userData,
	})

}

func userDataHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	user2 := getUser(name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user2)
}

func getFollowingPosts(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	allPosts, err := getFollowingAllPost(db, name)
	if err != nil {
		fmt.Println("err geting posts", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allPosts)
}

// for the explore page
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	allPosts, err := getAllPost(db, name)
	if err != nil {
		fmt.Println("err geting posts", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allPosts)
}

func getGroupPosts(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	groupId, _ := strconv.Atoi(r.FormValue("groupId"))
	allPosts, err := getGroupPost(db, groupId, name)
	if err != nil {
		fmt.Println("err geting posts", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allPosts)
}

func posthandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	var imgByte []byte
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	text := r.FormValue("postText")
	postType := r.FormValue("postType")
	img, _, err := r.FormFile("postImg")
	if err == nil {
		imgByte, _ = readImg("", img)
	}
	groupId := r.FormValue("groupId")
	fmt.Println("postType", postType)
	t := time.Now()
	time := t.Format("15:04 01-02")
	postid, _ := createPost(db, text, name, postType, imgByte, groupId)

	w.WriteHeader(http.StatusOK)

	imgBase64 := encodeToBase64(imgByte)
	userPfpPic := getUserPfp(db, name)
	brodcast <- map[string]interface{}{
		"type":     "post",
		"postID":   postid,
		"postType": postType,
		"userPfp":  userPfpPic,
		"username": name,
		"text":     text,
		"img":      imgBase64,
		"time":     time,
		"likes":    0,
	}

}
func commentHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	pId := r.FormValue("postId")
	postId:= atoi(pId)
	post,err := getOnePost(db,pId,name)
	if err != nil {
		fmt.Println("post deleted", err)
	}
	var imgByte []byte
	fmt.Println("-----atoi-----post id ", postId)
	text := r.FormValue("comment")
	img, _, err := r.FormFile("comImg")
	if err == nil {
		imgByte, _ = readImg("", img)
	}
	fmt.Println("come", text)
	t := time.Now()
	time := t.Format("15:04 01-02")
	comId := createComment(db, text, name, postId, imgByte)
	stringComId := strconv.Itoa(int(comId))
	likes := 0
	liked := false
	user := getUser(name)
	com := map[string]interface{}{
		"post":     int64(postId),
		"username": name,
		"img":      user.Pfp,
		"text":     text,
		"time":     time,
		"comId":    stringComId,
		"likes":    likes,
		"liked":    liked,
		"comPic":   encodeToBase64(imgByte),
	}
	brodcast <- map[string]interface{}{
		"type": "comment",
		"data": com,
	}

			if post.Username != name{
				_, err = db.Exec("INSERT INTO not_requests (sender, receiver,nottype,actionId) VALUES ($1, $2,$3,$4)", name, post.Username, "comment",pId)
		}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(com)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userInput := strings.ToLower(r.FormValue("username"))
		passInput := r.FormValue("password")
		var dbUsername string
		var dbPassword string

		query := `SELECT  username, password FROM users WHERE username = ? OR email = ?`
		err := db.QueryRow(query, userInput, userInput).Scan(&dbUsername, &dbPassword)
		if err != nil {
			fmt.Println("user msh mojod")
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		passInput = strings.TrimSpace(passInput)
		fmt.Println("input password hashed", passInput)

		if !checkPass(passInput, dbPassword) {
			fmt.Println("wrong password")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if strings.Contains(userInput, "@") {
			username := getUnFromEmail(userInput)
			setSession(w, username)
		} else {
			setSession(w, userInput)
		}
	}
	w.WriteHeader(http.StatusOK)
}
func regHandler(w http.ResponseWriter, r *http.Request) {
	var imgByte []byte
	if r.Method == http.MethodPost {
		cwd, _ := os.Getwd()
		img, _, err := r.FormFile("img")
		if err != nil {
			fmt.Println("defult img")
			imgByte, err = readImg(cwd+"/profile.jpg", nil)
			if err != nil {
				fmt.Println("err in img ", err)
			}
		} else {
			imgByte, _ = readImg("", img)
		}
		password := strings.TrimSpace(r.FormValue("password"))
		username := strings.ToLower(r.FormValue("username"))

		password, _ = passwordHash(r.FormValue("password"))
		fmt.Println("passhashed", password)
		err = CreateUser(db, imgByte, username, r.FormValue("email"),
			password, r.FormValue("gender"),
			r.FormValue("firstname"), r.FormValue("lastname"),
			r.FormValue("age"), r.FormValue("aboutMe"), w, r)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			setSession(w, username)
			w.WriteHeader(http.StatusOK)
		}
	}

}
func clickedPostHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)

	fmt.Println("in one post")
	if r.Method != http.MethodPost {
		return
	}
	postId := r.FormValue("postid")
	fmt.Println("post id ", postId)
	onePost, err := getOnePost(db, postId, name)
	if err != nil {
		fmt.Println("err geting one post", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(onePost)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)

	fmt.Println("in like")
	if r.Method != http.MethodPost {
		return
	}
	postId := r.FormValue("postID")
	fmt.Println("post lie id ", postId)
	postLikes := likingPost(db, name, postId)
	post,_ := getOnePost(db,postId,name)
		sPfp := getUserPfp(db, name)
	
	fmt.Println("post.user",post.Username)
	brodcast <- map[string]interface{}{
		"type": "like",
		"actId":   postId,
		"likeCounter": postLikes,
		"sender": name,
		"sender_pfp" : sPfp,
		"text" : "liked your post",
	}
		if post.Username == name{
			w.WriteHeader(http.StatusOK)
			return
		}
		_, err := db.Exec("INSERT INTO not_requests (sender, receiver,nottype,actionId) VALUES ($1, $2,$3,$4)", name, post.Username, "likePost",postId)
		if err != nil{
			fmt.Println("err in insert noti for liked post",err)
			return
		}
		w.WriteHeader(http.StatusOK)
}
func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)

	fmt.Println("in dislike")
	if r.Method != http.MethodPost {
		return
	}
	postId := r.FormValue("postID")
	post,err := getOnePost(db,postId,name)
	if err != nil {
		fmt.Println("post is deleted",err)
		return
	} 

	postLikes := disLikePost(db, name, postId)

	brodcast <- map[string]interface{}{
		"type":        "dislike",
		"postId":      postId,
		"likeCounter": postLikes,
	}
		_, err = db.Exec("DELETE FROM not_requests WHERE sender = ? AND receiver = ? AND actionId = ?", name,post.Username,postId)
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func dislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in c")
	name, _ := getSession(r)

	fmt.Println("in c dislike")
	if r.Method != http.MethodPost {
		return
	}
	comId := r.FormValue("comId")
	fmt.Println("comid", comId)

	comLikes := disLikeCom(db, name, atoi(comId))

	brodcast <- map[string]interface{}{
		"type":        "comLike",
		"comId":       comId,
		"likeCounter": comLikes,
	}
	w.WriteHeader(http.StatusOK)
}
func likeCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("like coms")
	name, _ := getSession(r)

	fmt.Println("in c like")
	if r.Method != http.MethodPost {
		return
	}
	comId := r.FormValue("comId")

	fmt.Println("comid", comId)
	comLikes := likingCom(db, name, atoi(comId))

	fmt.Println("comid", comId)

	brodcast <- map[string]interface{}{
		"type":        "comLike",
		"comId":       comId,
		"likeCounter": comLikes,
	}
	w.WriteHeader(http.StatusOK)
}
func allUserHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	users := getAllUser(db, name)
	groups := getAllGroup(db, name)
	users = append(users, groups...)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	_, logedin := getSession(r)
	if logedin {
		http.Error(w, "user is loged in", http.StatusAccepted)
	} else {
		http.Error(w, "user not loged in", http.StatusUnauthorized)
	}
}
func messegesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in msg")
	name, _ := getSession(r) // the sender
	if r.Method != http.MethodPost {
		return
	}
	rec := r.FormValue("rec") // Receiver (User or Group)
	msg := r.FormValue("text")

	fmt.Println("Message:", msg, "Sender:", name, "Receiver:", rec)

	// Check if the receiver is a group and get its group_id
	var groupID int
	err := db.QueryRow(`SELECT id FROM groups WHERE name = ?`, rec).Scan(&groupID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error checking if receiver is a group:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pfp := getUserPfp(db, name)

	saveMsg(db, name, rec, msg)
	if groupID > 0 {
		sendGroupMessage(rec, groupID, msg, pfp, name)
		full := map[string]interface{}{
			"type":   "msg",
			"Sender": name,
			"to":     rec,
			"text":   msg,
			"pfp":    pfp,
			"RecPfp": "",
		}

		brodcast <- full
	} else {
		RecNamePfp := getUserPfp(db, rec)

		full := map[string]interface{}{
			"type":   "msg",
			"Sender": name,
			"to":     rec,
			"text":   msg,
			"pfp":    pfp,
			"RecPfp": RecNamePfp,
		}

		brodcast <- full
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(full)
	}
}

func sendGroupMessage(sender string, groupID int, msg string, senderPfp string, name string) {
	// Fetch all members of the group (excluding sender)
	rows, err := db.Query(`SELECT user FROM members WHERE group_id = ? AND user != ?`, groupID, name)
	if err != nil {
		fmt.Println("Error fetching group members:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var member string
		err := rows.Scan(&member)
		if err != nil {
			fmt.Println("Error scanning group member:", err)
			continue
		}

		// Get recipient profile picture
		RecNamePfp := getUserPfp(db, member)

		// Broadcast the message to each member
		full := map[string]interface{}{
			"type":   "msg",
			"Sender": sender,
			"to":     member, // Send to individual members
			"text":   msg,
			"pfp":    senderPfp,
			"RecPfp": RecNamePfp,
		}

		brodcast <- full
	}
}
func getMsgHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	rec := r.FormValue("receiver")
	var isGroup bool
	err := db.QueryRow(`SELECT COUNT(*) FROM groups WHERE name = ?`, rec).Scan(&isGroup)
	if err != nil {
		fmt.Println("Error checking if receiver is a group:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var allMsgs []Msg

	if isGroup {
		allMsgs, err = getGroupMsg(db, rec)
	} else {
		allMsgs, err = getAllMsg(db, name, rec)
	}

	if err != nil {
		fmt.Println("Error getting messages:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msgs": allMsgs,
	})
}

func logout_handler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)

	cookie, err := r.Cookie("sID")
	if err != nil {
		http.Error(w, "user not loged in", http.StatusSeeOther)
		return
	}
	delete(userSession, cookie.Value)
	delete(sessionsMap, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "sID",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	brodcast <- map[string]interface{}{
		"type": "off",
		"user": name,
	}

	http.Error(w, "user not loged in", http.StatusSeeOther)
	fmt.Println("sessins after logout", sessionsMap)
}

func followerAcc(w http.ResponseWriter, r *http.Request) {
	follower_username, _ := getSession(r)     //making the requset is the one who follow
	following_username := r.FormValue("fing") //the guy u want to follow
	createfollow(db, following_username, follower_username)
	idacc := r.FormValue("id")
	fmt.Println("id",idacc)
	_, err := db.Exec("DELETE FROM not_requests WHERE id =?", idacc)
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}
	_, err = db.Exec("INSERT INTO not_requests (sender, receiver,nottype) VALUES ($1, $2,$3)", follower_username, following_username, "startFollow")




	w.WriteHeader(http.StatusOK)

}
func followerRej(w http.ResponseWriter, r *http.Request) {
	//  follower_username, _ := getSession(r)     //making the requset is the one who follow
	//  following_username := r.FormValue("fing") //the guy u want to follow
	//just delet the notification from the table
	idacc := r.FormValue("id")
	_, err := db.Exec("DELETE FROM not_requests WHERE id =?", idacc)
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func follower(w http.ResponseWriter, r *http.Request) {

	follower_username, _ := getSession(r)     //making the requset is the one who follow
	following_username := r.FormValue("fing") //the guy u want to follow

	fmt.Println("follower_username", follower_username)
	fmt.Println("following_username", following_username)

	createfollow(db, follower_username, following_username)
}
func unfollow(w http.ResponseWriter, r *http.Request) {
	follower_username, _ := getSession(r)     //making the requset is the one who follow
	following_username := r.FormValue("fing") //the guy u want to follow

	fmt.Println("follower_username", follower_username)
	fmt.Println("following_username", following_username)

	unfollowUser(follower_username, following_username)
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	var imgByte []byte
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	cwd, _ := os.Getwd()
	img, _, err := r.FormFile("groupImage")
	if err != nil {
		fmt.Println("defult img")
		imgByte, err = readImg(cwd+"/profile.jpg", nil)
		if err != nil {
			fmt.Println("err in img ", err)
		}
	} else {
		imgByte, _ = readImg("", img)
	}
	groupName := r.FormValue("groupName")
	groupDesc := r.FormValue("groupDesc")

	fmt.Println(groupName, groupDesc)
	groupID, err := createGroup(db, imgByte, groupName, groupDesc, name)
	if err != nil {
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Broadcast new group creation (like how you did for posts)
	brodcast <- map[string]interface{}{
		"type":       "group",
		"groupID":    groupID,
		"creator":    name,
		"groupName":  groupName,
		"groupDesc":  groupDesc,
		"created_at": time.Now().Format("15:04 01-02"),
	}
}

func getGroupsHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getSession(r)

	rows, err := db.Query("SELECT id, pfp, name, description, created_at FROM groups ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var myGroups []map[string]interface{}
	var allGroups []map[string]interface{}

	for rows.Next() {
		var id int
		var name, description, createdAt string
		var tempic []byte
		err := rows.Scan(&id, &tempic, &name, &description, &createdAt)
		if err != nil {
			http.Error(w, "Error reading groups", http.StatusInternalServerError)
			return
		}

		pic := encodeToBase64(tempic)
		var exists int
		err = db.QueryRow("SELECT COUNT(*) FROM members WHERE `user` = ? AND group_id = ?", user, id).Scan(&exists)
		if err != nil {
			http.Error(w, "Error checking membership", http.StatusInternalServerError)
			return
		}

		group := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"created_at":  createdAt,
			"pic":         pic,
		}

		if exists > 0 {
			myGroups = append(myGroups, group)
		} else {
			allGroups = append(allGroups, group)
		}
	}

	response := map[string]interface{}{
		"myGroups":  myGroups,
		"allGroups": allGroups,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getGroupDetailsHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getSession(r)
	groupID := r.FormValue("groupId")
	var name, description, createdAt, admin string
	var tempic []byte
	query := `
	SELECT g.pfp, g.name, g.description, g.created_at, g.creator
	FROM groups g
	WHERE g.id = ?
`
	err := db.QueryRow(query, groupID).Scan(&tempic, &name, &description, &createdAt, &admin)
	if err != nil {
		fmt.Println("Error: Group not found", err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}
	var count int
	query2 := `SELECT COUNT(*) FROM members WHERE group_id = ?;`

	err2 := db.QueryRow(query2, groupID).Scan(&count)
	if err2 != nil {
		fmt.Println("impossible to have err")
	}
	pic := encodeToBase64(tempic)
	IsMem := checkIfMember(user, groupID)
	IsInvitee := checkIfInvitee(user, groupID)
	response := map[string]interface{}{
		"id":          groupID,
		"name":        name,
		"description": description,
		"created_at":  createdAt,
		"admin":       admin,
		"IsAdmin":     admin == user,
		"NumOfMems":   count,
		"IsMem":       IsMem,
		"IsInv":       IsInvitee,
		"pic":         pic,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkIfMember(name, grId string) bool {
	var count int
	query := `SELECT COUNT(*) FROM members WHERE group_id = ? AND user = ?;`

	err := db.QueryRow(query, grId, name).Scan(&count)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	return count > 0
}

func checkIfInvitee(name, grId string) bool {
	var count int
	query := `SELECT COUNT(*) FROM invites WHERE group_id = ? AND invitee = ? AND status = ?;`

	err := db.QueryRow(query, grId, name, "pending").Scan(&count)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	return count > 0
}

func joinGroup(w http.ResponseWriter, r *http.Request) {

	user, _ := getSession(r)
	groupId := r.FormValue("groupId")

	_, err := db.Exec("INSERT INTO invites (group_id, invitee) VALUES (?, ?) ON CONFLICT DO NOTHING", groupId, user)
	if err != nil {
		fmt.Println("error insert to join", err)
		return
	}

}

func notJoinGroup(w http.ResponseWriter, r *http.Request) {

	user, _ := getSession(r)
	groupId := r.FormValue("groupId")
	groupName := r.FormValue("groupName")

	_, err := db.Exec("DELETE FROM invites WHERE group_id = ? AND invitee = ?", groupId, user)
	if err != nil {
		fmt.Println("error insert to join", err)
		return
	}

	var admin string
	err = db.QueryRow("SELECT creator FROM groups WHERE id = ?", groupId).Scan(&admin)
	if err != nil {
		fmt.Println("Error fetching group admin:", err)
		return
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE sender =? AND receiver =? AND nottype =?", user, admin, fmt.Sprintf("join group %s groupId:%s", groupName, groupId))
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}

}

func getGroupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("groupId")
	rows, err := db.Query("SELECT id, invitee, inviter FROM invites WHERE group_id = ? AND status = 'pending'", groupID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id int
		var inviter string
		var username string
		rows.Scan(&id, &username, &inviter)
		pfp := getUserPic(username)
		requests = append(requests, map[string]interface{}{
			"id":       id,
			"username": username,
			"pfp":      pfp,
			"inviter":  inviter,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func getGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("groupId")
	rows, err := db.Query("SELECT user FROM members WHERE group_id = ?", groupID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []map[string]interface{}
	for rows.Next() {
		var username string
		rows.Scan(&username)
		userPfp := getUserPfp(db,username)
		members = append(members, map[string]interface{}{
			"username": username,
			"pfp": userPfp,
		})

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func acceptInviteHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getSession(r)
	groupID := r.FormValue("groupId")
	userName := r.FormValue("userName")
	requestId := r.FormValue("requestId")
	groupName := r.FormValue("groupName")
	fmt.Println(groupID)
	fmt.Println(groupName)

	_, err := db.Exec("INSERT INTO members (group_id, user) VALUES (?, ?)", groupID, userName)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM invites WHERE id = ?", requestId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE sender =? AND receiver =? AND nottype =?", userName, user, fmt.Sprintf("join group %s groupId:%s", groupName, groupID))
	if err != nil {
		fmt.Println("abd is")
		fmt.Println("error from followeracc", err)
		return
	}
}

func rejectInviteHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getSession(r)
	groupID := r.FormValue("groupId")
	userName := r.FormValue("userName")
	requestId := r.FormValue("requestId")
	groupName := r.FormValue("groupName")

	_, err := db.Exec("DELETE FROM invites WHERE id = ?", requestId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE sender =? AND receiver =? AND nottype =?", userName, user, fmt.Sprintf("join group %s groupId:%s", groupName, groupID))
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}
}

func outOfGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("groupId")
	user, _ := getSession(r)
	fmt.Println("groupID:", groupID)
	fmt.Println("userName:", user)

	_, err := db.Exec("DELETE FROM members WHERE group_id = ? AND user = ?", groupID, user)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	groupID := r.FormValue("groupId")
	creator, _ := getSession(r)
	title := r.FormValue("title")
	description := r.FormValue("description")
	eventDatetime := r.FormValue("dateTime")
	fmt.Println(groupID)
	fmt.Println(title)
	fmt.Println(eventDatetime)
	fmt.Println(creator)
	_, err := db.Exec(`
        INSERT INTO events (group_id, creator, title, description, event_datetime) 
        VALUES (?, ?, ?, ?, ?)`, groupID, creator, title, description, eventDatetime)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Event created successfully")
}

// Fetch Events for a Group
func getEventsHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("groupId")
	user, _ := getSession(r)

	rows, err := db.Query("SELECT id, creator, title, description, event_datetime, created_at FROM events WHERE group_id = ? ORDER BY events.created_at DESC", groupID)
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var id int
		var title, creator, description, eventDatetime, createdAt string
		err = rows.Scan(&id, &creator, &title, &description, &eventDatetime, &createdAt)
		if err != nil {
			http.Error(w, "Failed to scan events", http.StatusInternalServerError)
			return
		}
		goingCount, notGoingCount := countEventRes((id))
		pic := getUserPic(creator)
		events = append(events, map[string]interface{}{
			"id":            id,
			"title":         title,
			"description":   description,
			"eventDatetime": eventDatetime,
			"createdAt":     createdAt,
			"goingCount":    goingCount,
			"notGoingCount": notGoingCount,
			"creator":       creator,
			"pic":           pic,
			"me":            user == creator,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Respond to an Event
func respondToEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.FormValue("eventId")
	user, _ := getSession(r)
	response := r.FormValue("response") // Should be 'Going' or 'Not going'
	fmt.Println(eventID)
	fmt.Println(user)
	fmt.Println(response)

	_, err := db.Exec(`
        INSERT INTO event_responses (event_id, user, response) 
        VALUES (?, ?, ?) 
        ON CONFLICT(event_id, user) DO UPDATE SET response = excluded.response, responded_at = CURRENT_TIMESTAMP`,
		eventID, user, response)

	if err != nil {
		http.Error(w, "Failed to respond to event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Response recorded successfully")
}

func countEventRes(eventID int) (int, int) {

	var goingCount, notGoingCount int

	err := db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN response = 'Going' THEN 1 END) AS going,
			COUNT(CASE WHEN response = 'Not going' THEN 1 END) AS notGoing
		FROM event_responses
		WHERE event_id = ?`, eventID).Scan(&goingCount, &notGoingCount)

	if err != nil {
		fmt.Println("Failed to fetch response counts")
		return 0, 0
	}
	return goingCount, notGoingCount
}

func followReq(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	username := r.FormValue("username")
	namePfp := getUserPfp(db, name)

	id, err := db.Exec("INSERT INTO not_requests (sender, receiver,nottype) VALUES ($1, $2,$3)", name, username, "followreq")
	if err != nil {
		fmt.Println("eroor in followreq", err)
		return
	}
	reqId, _:= id.LastInsertId()


	brodcast <- map[string]interface{}{
		"type":    "followreq",
		"id":reqId,
		"sender_pfp":  namePfp,
		"sender": name,
		"to":      username,
		"text":    "Requested to follow you ",
		"time":    time.Now().Format("15:04 01-02"),
	}
}

func postDelet(w http.ResponseWriter, r *http.Request) {
	post_id := r.FormValue("id")

	err := deletePost(atoi(post_id))
	if err != nil {
		i := 0
		fmt.Println("error deleting post trying again", err)
		i++
		if i == 2 {
			fmt.Println("couldn't delet post", err)
			return

		}
		_ = deletePost(atoi(post_id))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}



func getNoti(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	allNoti := getAllnoti(name)
	//return all notification in a resp
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allNoti)
}

func comDelet(w http.ResponseWriter, r *http.Request) {
	com_id := r.FormValue("id")
	fmt.Println("com_id", com_id)
	err := deleteCom(atoi(com_id))
	if err != nil {
		i := 0
		fmt.Println("error deleting com trying again", err)
		i++
		if i == 2 {
			fmt.Println("couldn't delet com", err)
			return

		}
		_ = deleteCom(atoi(com_id))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func sendNotification(w http.ResponseWriter, r *http.Request) {
	name, _ := getSession(r)
	receiver := r.FormValue("receiver")
	notiType := r.FormValue("type")

	_, err := db.Exec("INSERT INTO not_requests (sender, receiver, nottype) VALUES (?, ?, ?)", name, receiver, notiType)
	if err != nil {
		fmt.Println("Error inserting notification:", err)
		return
	}

	brodcast <- map[string]interface{}{
		"type":    notiType,
		"sender":  name,
		"to":      receiver,
		"message": name + " " + notiType + " you",
		"time":    time.Now().Format("15:04 01-02"),
	}

	w.WriteHeader(http.StatusOK)
}


func deleteNot(w http.ResponseWriter, r *http.Request){
	id := r.FormValue("id")
	_, err := db.Exec("DELETE FROM not_requests WHERE id =?", id)
		if err != nil {
			fmt.Println("err deleting not",err)
			return
		}
    w.WriteHeader(http.StatusOK)
			
}

func getFollowers(w http.ResponseWriter, r *http.Request) {
	username,_ := getSession(r) 
	
	fmt.Println("Fetching followers for:", username)

	rows, err := db.Query("SELECT follower_username FROM follows WHERE following_username = $1", username)
	if err != nil {
		fmt.Println("Error fetching followers:", err)
		return
	}
	defer rows.Close()

	var followers []map[string]interface{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			fmt.Println("Error scanning follower:", err)
			continue
		}
		userPfp := getUserPfp(db,username)
		followers = append(followers, map[string]interface{}{
    	"username": username,
    	"pfp":      userPfp,
})
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
	w.WriteHeader(http.StatusOK)

}

func getFollowing(w http.ResponseWriter, r *http.Request) {
		username,_ := getSession(r) 
	fmt.Println("Fetching following list for:", username)

	rows, err := db.Query("SELECT following_username FROM follows WHERE follower_username = $1", username)
	if err != nil {
		fmt.Println("Error fetching following:", err)
		return
	}
	defer rows.Close()

	var following []map[string]interface{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			fmt.Println("Error scanning following:", err)
			continue
		}
		fmt.Println("username after scan", username)
		userPfp := getUserPfp(db,username)
		following = append(following, map[string]interface{}{
    	"username": username,
    	"pfp":      userPfp,
})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
	w.WriteHeader(http.StatusOK)
}
func getMutualFollowers(w http.ResponseWriter, r *http.Request) {
	username, _ := getSession(r)
	groupId := r.FormValue("groupId")
	query := `
	SELECT f1.following_username
	FROM follows f1
	JOIN follows f2 
	ON f1.follower_username = f2.following_username 
	AND f1.following_username = f2.follower_username
	WHERE f1.follower_username = ? AND f1.following_username != ?;
	`

	rows, err := db.Query(query, username, username)
	if err != nil {
		fmt.Println("err to read Mutual Followers")
	}
	defer rows.Close()

	var followers []map[string]interface{}
	for rows.Next() {
		var follower string
		rows.Scan(&follower)
		var count int
		query2 := `SELECT COUNT(*) FROM members WHERE group_id = ? AND user = ?;`

		err2 := db.QueryRow(query2, groupId, follower).Scan(&count)
		if err2 != nil {
			fmt.Println("impossible to have err")
		}

		var count2 int
		query3 := `SELECT COUNT(*) FROM invites WHERE group_id = ? AND invitee = ?;`

		err3 := db.QueryRow(query3, groupId, follower).Scan(&count2)
		if err3 != nil {
			fmt.Println("impossible to have err")
		}
		if count == 0 && count2 == 0 {
			pfp := getUserPic(follower)
			followers = append(followers, map[string]interface{}{
				"username": follower,
				"pfp":      pfp,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func inviteFriend(w http.ResponseWriter, r *http.Request) {

	user, _ := getSession(r)
	groupId := r.FormValue("groupId")
	userName := r.FormValue("userName")
	namePfp := getUserPfp(db, user)
	groupName := r.FormValue("groupName")

	_, err := db.Exec("INSERT INTO invites (group_id, inviter, invitee, status) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING", groupId, user, userName, "invited")
	if err != nil {
		fmt.Println("error insert to join", err)
		return
	}

	result, err := db.Exec(`
    INSERT OR IGNORE INTO not_requests (sender, receiver, nottype) 
    VALUES ($1, $2, $3)`, user, userName, fmt.Sprintf("invite group %s groupId:%s", groupName, groupId))

	if err != nil {
		fmt.Println("eroor in followreq", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		brodcast <- map[string]interface{}{
			"type":    "invite group",
			"hisPfp":  namePfp,
			"hisName": user,
			"to":      userName,
			"text":    fmt.Sprintf("invite to join %s group", groupName),
			"time":    time.Now().Format("15:04 01-02"),
		}
	}

}

func uninviteFriend(w http.ResponseWriter, r *http.Request) {

	user, _ := getSession(r)
	groupId := r.FormValue("groupId")
	userName := r.FormValue("userName")
	groupName := r.FormValue("groupName")
	// var admin string
	// err := db.QueryRow("SELECT creator FROM groups WHERE id = ?", groupId).Scan(&admin)
	// if err != nil {
	// 	fmt.Println("Error fetching group admin:", err)
	// 	return
	// }

	_, err := db.Exec("DELETE FROM invites WHERE group_id = ? AND inviter = ? AND invitee = ?",
		groupId, user, userName)
	if err != nil {
		fmt.Println("error insert to join", err)
		return
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE sender =? AND receiver = ? AND nottype = ?", user, userName, fmt.Sprintf("invite group %s groupId:%s", groupName, groupId))
	if err != nil {
		fmt.Println("error 5")
		fmt.Println("error from followeracc", err)
		return
	}

}

func acceptInviteFriendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("suiiiiiiiiiiiiiii")
	user, _ := getSession(r)
	groupID := r.FormValue("groupId")
	userName := r.FormValue("inviter")

	groupName := r.FormValue("groupName")
	idacc := r.FormValue("id")
	fmt.Println(user)
	namePfp := getUserPfp(db, user)
	fmt.Println(groupID)
	fmt.Println(userName)
	var admin string
	err := db.QueryRow("SELECT creator FROM groups WHERE id = ?", groupID).Scan(&admin)
	if err != nil {
		fmt.Println("Error fetching group admin:", err)
		return
	}
	fmt.Println(admin)
	if userName == admin {
		_, err = db.Exec("INSERT INTO members (group_id, user) VALUES (?, ?)", groupID, user)
		if err != nil {
			fmt.Println("error 0")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("DELETE FROM invites WHERE group_id = ? AND inviter = ? AND invitee = ?",
			groupID, userName, user)
		if err != nil {
			fmt.Println("error 1")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {

		_, err = db.Exec("UPDATE invites SET status = 'pending' WHERE group_id = ? AND inviter = ? AND invitee = ?",
			groupID, userName, user)
		if err != nil {
			fmt.Println("error 3")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		result, err := db.Exec(`
    INSERT OR IGNORE INTO not_requests (sender, receiver, nottype) 
    VALUES ($1, $2, $3)`, user, admin, fmt.Sprintf("join group %s groupId:%s", groupName, groupID))

		if err != nil {
			fmt.Println("eroor in followreq", err)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			brodcast <- map[string]interface{}{
				"type":    "join group",
				"hisPfp":  namePfp,
				"hisName": user,
				"to":      admin,
				"text":    fmt.Sprintf("requested to join %s group", groupName),
				"time":    time.Now().Format("15:04 01-02"),
			}
		}
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE id =?", idacc)
	if err != nil {
		fmt.Println("error from followeracc", err)
		return
	}
}

func rejectInviteFriendHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := getSession(r)
	groupID := r.FormValue("groupId")
	userName := r.FormValue("userName")
	// requestId := r.FormValue("requestId")
	// groupName := r.FormValue("groupName")
	idacc := r.FormValue("id")

	_, err := db.Exec("DELETE FROM invites WHERE group_id = ? AND inviter = ? AND invitee = ?",
		groupID, userName, user)
	if err != nil {
		fmt.Println("error 4")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM not_requests WHERE id =?", idacc)
	if err != nil {
		fmt.Println("error 5")
		fmt.Println("error from followeracc", err)
		return
	}
}