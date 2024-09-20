package handlers

import (
	"context"

	"github.com/code-troopers/postitsonline/database"
	"github.com/gofrs/uuid/v5"
)

func GetAllBoards() ([]database.Board, error) {
	rows, err := database.DB.Query(context.Background(), "SELECT id, name FROM boards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []database.Board
	for rows.Next() {
		var board database.Board
		err := rows.Scan(&board.ID, &board.Name)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func createBoard(name string) (*database.Board, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	boardID := u2.String()
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO boards (id, name) VALUES ($1, $2)",
		boardID, name)
	if err != nil {
		return nil, err
	}
	return &database.Board{ID: boardID, Name: name, Postits: []database.Postit{}}, nil
}

func renameBoard(boardID, name string) error {
	_, err := database.DB.Exec(context.Background(),
		"UPDATE boards SET name = $1 WHERE id = $2",
		name, boardID)
	if err != nil {
		return err
	}
	return nil
}

func deleteBoard(boardID string) error {
	_, err := database.DB.Exec(context.Background(),
		"DELETE FROM postits WHERE board_id = $1",
		boardID)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec(context.Background(),
		"DELETE FROM boards WHERE id = $1",
		boardID)
	if err != nil {
		return err
	}
	return nil
}
