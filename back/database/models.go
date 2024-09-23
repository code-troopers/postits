package database

type Board struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Postits []Postit `json:"postits"`
}

type Postit struct {
	ID      string `json:"id"`
	BoardID string `json:"boardId"`
	Text    string `json:"text"`
	PosX    int    `json:"posX"`
	PosY    int    `json:"posY"`
	Votes   int    `json:"votes"`
	Show    bool   `json:"show"`
	Author  User   `json:"author"`
	Weight  int    `json:"weight"`
}

type User struct {
	ID         string `json:"id"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}
