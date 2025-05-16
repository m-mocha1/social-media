package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var (
	tbl *template.Template
	db  = makeSQL()
	fs  http.Handler
)

func main() {
	err := os.Mkdir("SQL", 0o755)
	if err != nil {
	}
	fmt.Println("http://localhost:8080/")

	// sql_tables(db)
	go brodcastUpdate()
	ApplyMigrations(db, "./migrate")

	// fs = http.FileServer(http.Dir("../frontEnd/dist"))
	// http.Handle("/", fs)

	//wep-socket
	http.HandleFunc("/ws", ws)

	//login
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)

	//userData
	http.HandleFunc("/data", userDataHandler)
	http.HandleFunc("/getFollowingPost", getFollowingPosts)
	http.HandleFunc("/allUsers", allUserHandler)

	//make a post
	http.HandleFunc("/posting", posthandler)

	//post actions
	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/unlike", dislikeHandler)
	http.HandleFunc("/deletPost", postDelet)

	// comment actions
	http.HandleFunc("/com", commentHandler)
	http.HandleFunc("/likeCom", likeCommentHandler)
	http.HandleFunc("/dislikeCom", dislikeCommentHandler)
	http.HandleFunc("/deletCom", comDelet)

	
	//profile actions
	http.HandleFunc("/get-followers", getFollowers)
	http.HandleFunc("/get-following", getFollowing)
	http.HandleFunc("/follower", follower)
	http.HandleFunc("/unfollow", unfollow)
	http.HandleFunc("/followReq", followReq)
	http.HandleFunc("/pri", pri)
	http.HandleFunc("/pub", pub)
	http.HandleFunc("/alPri", almostPri)

	//posts data
	http.HandleFunc("/getAllPost", getAllPosts)
	http.HandleFunc("/getUserPost", getUserPost)
	http.HandleFunc("/onepost", clickedPostHandler)

	//Auth
	http.HandleFunc("/Auth", checkAuthHandler)

	// msg actions
	http.HandleFunc("/msg", messegesHandler)
	http.HandleFunc("/getMsg", getMsgHandler)

	//groups
	http.HandleFunc("/createGroup", createGroupHandler)
	http.HandleFunc("/getGroups", getGroupsHandler)
	http.HandleFunc("/getGroup", getGroupDetailsHandler)
	http.HandleFunc("/joinGroup", joinGroup)
	http.HandleFunc("/notJoinGroup", notJoinGroup)
	http.HandleFunc("/getGroupRequests", getGroupRequestsHandler)
	http.HandleFunc("/getGroupMembers", getGroupMembersHandler)
	http.HandleFunc("/getMutualFollowers", getMutualFollowers)
	http.HandleFunc("/acceptInvite", acceptInviteHandler)
	http.HandleFunc("/rejectInvite", rejectInviteHandler)
	http.HandleFunc("/outOfGroup", outOfGroupHandler)
	http.HandleFunc("/getGroupPosts", getGroupPosts)
	http.HandleFunc("/InviteFriend", inviteFriend)
	http.HandleFunc("/UnInviteFriend", uninviteFriend)
	http.HandleFunc("/inviteFrAcc", acceptInviteFriendHandler)
	http.HandleFunc("/inviteFrRej", rejectInviteFriendHandler)

	//events
	http.HandleFunc("/createEvent", createEventHandler)
	http.HandleFunc("/getEvents", getEventsHandler)
	http.HandleFunc("/respondToEvent", respondToEventHandler)
	http.HandleFunc("/deletEvent", deleteEvent)

	//notifications
	http.HandleFunc("/sendNotification", sendNotification)
	http.HandleFunc("/deleteNot", deleteNot)
	http.HandleFunc("/getNoti",getNoti)
	http.HandleFunc("/countNoti",countNoti)

	
	//follow requests
	http.HandleFunc("/followerRej", followerRej)
	http.HandleFunc("/followerAcc", followerAcc)

	//logout
	http.HandleFunc("/out", logout_handler)

	// tbl = template.Must(template.ParseFiles("../frontEnd/dist/index.html"))
	http.ListenAndServe(":8080", nil)
}
