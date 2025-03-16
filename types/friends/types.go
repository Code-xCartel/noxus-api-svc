package friends

type FriendResponse struct {
	ID       int    `json:"id"`
	NoxID    string `json:"noxId"`
	Username string `json:"username"`
	Status   string `json:"status"`
}
