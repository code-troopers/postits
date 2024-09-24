package handlers

import (
	"context"

	"github.com/code-troopers/postitsonline/database"
	"github.com/gofrs/uuid/v5"
)

func GetAllPostitsByBoardId(boardID string, user *database.User) ([]database.Postit, error) {
	rows, err := database.DB.Query(context.Background(),
		`SELECT p.id, p.board_id, p.text, p.pos_x, p.pos_y, p.author_id, p.votes, p.show, p.weight, u.given_name, u.family_name, u.email, u.picture
		 FROM postits p
		 JOIN users u ON p.author_id = u.id
		 WHERE board_id = $1`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postits []database.Postit
	for rows.Next() {
		var postit database.Postit
		err := rows.Scan(&postit.ID, &postit.BoardID, &postit.Text, &postit.PosX, &postit.PosY, &postit.Author.ID, &postit.Votes, &postit.Show, &postit.Weight,
			&postit.Author.GivenName, &postit.Author.FamilyName, &postit.Author.Email, &postit.Author.Picture)
		if err != nil {
			return nil, err
		}
		if postit.Author.ID != user.ID && !postit.Show {
			postit.Text = "**********"
		}
		postits = append(postits, postit)
	}
	return postits, nil
}

func getBiggestWeight(posX int, posY int, postitID string) int {
	rows, err := database.DB.Query(context.Background(),
		`SELECT p.weight
		 FROM postits p
		 WHERE pos_x > $1 - 200 and pos_x < $1 + 200 and pos_y > $2 - 200 and pos_y < $2 + 200 and id != $3`, posX, posY, postitID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	weight := 0
	for rows.Next() {
		var w int
		err := rows.Scan(&w)
		if err != nil {
			return 0
		}
		if w > weight {
			weight = w
		}
	}
	return weight
}

func CreatePostit(boardID, text string, posX, posY int, authorID string, show bool) (*database.Postit, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	postitID := u2.String()

	w := getBiggestWeight(posX, posY, "") + 1

	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO postits (id, board_id, text, pos_x, pos_y, author_id, weight, show) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		postitID, boardID, text, posX, posY, authorID, w, show)
	if err != nil {
		return nil, err
	}
	return &database.Postit{ID: postitID, BoardID: boardID, Text: text, PosX: posX, PosY: posY, Author: database.User{ID: authorID}, Weight: w, Show: show}, nil
}

func updatePostitContent(postitID, text, userId string) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET text = $1 WHERE id = $2 AND author_id = $3",
		text, postitID, userId)
	if err != nil {
		return err
	}
	return nil
}

func deletePostit(postitID string) error {
	_, err := database.DB.Exec(context.Background(),
		"DELETE FROM postits WHERE id = $1",
		postitID)
	if err != nil {
		return err
	}
	return nil
}

func endMovePostit(postitID string, posX, posY int) (int, error) {
	w := getBiggestWeight(posX, posY, postitID) + 1
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET pos_x = $1, pos_y = $2, weight = $3 WHERE id = $4",
		posX, posY, w, postitID)
	if err != nil {
		return 0, err
	}
	return w, nil
}
func movePostit(postitID string, posX, posY int) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET pos_x = $1, pos_y = $2 WHERE id = $4",
		posX, posY, postitID)
	if err != nil {
		return err
	}
	return nil
}

func addVote(postitId string) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET votes = votes + 1 WHERE id = $1",
		postitId)
	if err != nil {
		return err
	}
	return nil
}

func removeVote(postitId string) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET votes = votes - 1 WHERE id = $1",
		postitId)
	if err != nil {
		return err
	}
	return nil
}

func showPostits(boardId string, authorID string, show bool) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET show = $1 WHERE board_id = $2 AND author_id = $3",
		show, boardId, authorID)
	if err != nil {
		return err
	}
	return nil
}
