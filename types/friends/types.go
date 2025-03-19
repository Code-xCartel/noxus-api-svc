package friends

type FriendResponse struct {
	NoxID    string `json:"noxId"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

type SearchUserResponse struct {
	NoxID    string `json:"noxId"`
	Username string `json:"username"`
}
