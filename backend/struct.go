package main

import (
	"mime/multipart"
)

type uInfo struct {
	Pfp       string `json:"pfp"`
	Fname     string `json:"first_name"`
	Lname     string `json:"last_name"`
	Age       string `json:"Age"`
	Gender    string `json:"Gender"`
	Username  string `json:"Username"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Online    bool   `json:"online"`
	AboutMe   string `json:"AboutMe"`
	Public    string `json:"public"`
	Me        bool   `json:"Me"`
	Following bool   `json:"following"`
	Followers int    `json:"followers"`
}
type signin struct {
	UserNameEmail string `json:"username"`
	Password      string `json:"password"`
}
type RecivedPost struct {
	Username string                `json:"username"`
	Text     string                `json:"text"`
	Rpic     *multipart.FileHeader `jsom:"pic"`
}
type Opost struct {
	PostID    string `json:"postID"`
	UserPfp   string `json:"userPfp"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	Spic      string `json:"img"`
	Date      string `json:"time"`
	Likes     int    `json:"likes"`
	PostLiked bool   `json:"postLiked"`
	Comment   []Comment
	PostType  string `json:"postType"`
}
type Comment struct {
	ComPic   string `json:"comPic"`
	Pfpimg    string `json:"img"`
	User      string `json:"username"`
	CommentId string `json:"comId"`
	Text      string `json:"text"`
	Liked     bool   `json:"comLiked"`
	Date      string `json:"time"`
	Likes     int    `json:"likes"`
}
type Msg struct {
	SenderPfp  string `json:"pfp"`
	RecNamePfp string `json:"RecPfp"`
	Sender     string `json:"Sender"`
	RecName    string `json:"to"`
	Text       string `json:"text"`
}

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
