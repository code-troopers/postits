package handlers

import (
	"context"

	"github.com/code-troopers/postitsonline/database"
	"github.com/gofrs/uuid/v5"
)

func GetAllPostitsByBoardId(boardID string) ([]database.Postit, error) {
	rows, err := database.DB.Query(context.Background(), "SELECT id, board_id, text, pos_x, pos_y, author_id, votes, show FROM postits WHERE board_id = $1", boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postits []database.Postit
	for rows.Next() {
		var postit database.Postit
		err := rows.Scan(&postit.ID, &postit.BoardID, &postit.Text, &postit.PosX, &postit.PosY, &postit.Author.ID, &postit.Votes, &postit.Show)
		if err != nil {
			return nil, err
		}
		postits = append(postits, postit)
	}
	return postits, nil
}

func CreatePostit(boardID, text string, posX, posY int, authorID string) (*database.Postit, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	postitID := u2.String()

	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO postits (id, board_id, text, pos_x, pos_y, author_id) VALUES ($1, $2, $3, $4, $5, $6)",
		postitID, boardID, text, posX, posY, authorID)
	if err != nil {
		return nil, err
	}
	return &database.Postit{ID: postitID, BoardID: boardID, Text: text, PosX: posX, PosY: posY, Author: database.User{ID: authorID}}, nil
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

func movePostit(postitID string, posX, posY int) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE postits SET pos_x = $1, pos_y = $2 WHERE id = $3",
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
