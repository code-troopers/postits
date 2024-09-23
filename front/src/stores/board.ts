import { Board } from './../models/Board';
import { defineStore } from 'pinia'
import { sendMessage } from '../services/websocketService'
import { Actions } from '../Actions'
import axios from 'axios';
import type { StickyNote } from '@/models/StickyNote';
import type { Message } from '@/models/Message';

const API_URL = import.meta.env.VITE_API_URL;

export const useBoardStore = defineStore('board', {
  state: () => ({ boards: [] as Board[] }),
  actions: {
  onMessage(message: Message) {
    let board: Board | undefined;
    switch (message.action) {
      case Actions.NEW_BOARD:
        board = new Board();
        board.id = message.id;
        board.name = message.text;
        this.boards.push(board)
        break
      case Actions.RENAME_BOARD:
        board = this.boards.find((b) => b.id === message.id)
        if (board) {
          board.name = message.text
        }
        break
      case Actions.DELETE_BOARD:
        this.boards = this.boards.filter((b) => b.id !== message.boardId)
        break
      case Actions.NEW_POSTIT:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          if (board.postits == null) {
            board.postits = []
          }
          board.postits.push({
            id: message.id,
            posX: message.posX,
            posY: message.posY,
            show: false,
            votes: 0,
            author: message.author,
            weight: message.weight ?? 0
          })
        }
        break
      case Actions.UPDATE_CONTENT:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.text = message.text
          }
        }
        break

      case Actions.MOVE_POSTIT:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.posX = message.posX
            postit.posY = message.posY
            postit.weight = message.weight ?? 0
          }
        }
        break

      case Actions.DELETE_POSTIT:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          board.postits = board.postits.filter((p) => p.id !== message.id)
        }
        break

      case Actions.ADD_VOTE:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.votes++
          }
        }
        break

      case Actions.REMOVE_VOTE:
        board = this.boards.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.votes--
          }
        }
        break

      case Actions.SHOW_POSTITS:
      case Actions.HIDE_POSTITS:
        if (message.boardId) {
          this.initPostits(message.boardId)
        }
        break

    }
  },
  newBoard() {
    try {
      sendMessage({
        action: Actions.NEW_BOARD,
        text: `Board ${this.boards.length + 1}`,
      })
    } catch (error) {
      console.error(error)
    }
  },
  renameBoard(boardId: string, text: string) {
    try {
      sendMessage({
        action: Actions.RENAME_BOARD,
        boardId: boardId,
        text: text
      })
    } catch (error) {
      console.error(error)
    }

  },
  deleteBoard(boardId: string | undefined) {
    if (boardId === undefined) {
      throw new Error('Board ID is required')
    }
    try {
      sendMessage({
        action: Actions.DELETE_BOARD,
        boardId: boardId,
      })
    } catch (error) {
      console.error(error)
    }
  },
  newPostit(boardId: string, x: number, y: number) {
    try {
      sendMessage({
        action: Actions.NEW_POSTIT,
        boardId: boardId,
        posX: x,
        posY: y,
      })
    } catch (error) {
      console.error(error)
    }
  },
  updateContent(boardId: string, id: string, text: string) {
    try {
      sendMessage({
        action: Actions.UPDATE_CONTENT,
        boardId: boardId,
        id: id,
        text: text
      })
    } catch (error) {
      console.error(error)
    }
  },
  deletePostit(boardId: string, id: string | undefined) {
    if (id === undefined) {
      throw new Error('postit id is required')
    }
    try {
      sendMessage({
        action: Actions.DELETE_POSTIT,
        id: id,
        boardId: boardId
      })
    } catch (error) {
      console.error(error)
    }
  },

  getBoards() {
    return axios.get(`${API_URL}/api/boards`).then((response: any) => {
      this.boards = response.data
    })
  },

  getPostits(boardId: string) {
    const board = this.boards.find((b) => b.id === boardId)
    if (board?.postits) {
      return board.postits.sort((a, b) => a.weight - b.weight)
    }
    return []
  },

  async initPostits(boardId: string) {
    if (this.boards.length === 0) {
      await this.getBoards()
    }
    const board = this.boards.find((b) => b.id === boardId)
    if (board) {
      await axios.get(`${API_URL}/api/boards/${boardId}/postits`).then((response: any) => {
        board.postits = response.data
      })
    }
  },

  addVote(id: string, boardId: string) {
    try {
      sendMessage({
        action: Actions.ADD_VOTE,
        id: id,
        boardId: boardId
      })
    } catch (error) {
      console.error(error)
    }
  },

  removeVote(id: string, boardId: string) {
    try {
      sendMessage({
        action: Actions.REMOVE_VOTE,
        id: id,
        boardId: boardId
      })
    } catch (error) {
      console.error(error)
    }
  },

  showPostits(boardId: string, authorId: string) {
    try {
      sendMessage({
        action: Actions.SHOW_POSTITS,
        boardId: boardId,
        authorId: authorId
      })
    } catch (error) {
      console.error(error)
    }
  },

  hidePostits(boardId: string, authorId: string) {
    try {
      sendMessage({
        action: Actions.HIDE_POSTITS,
        boardId: boardId,
        authorId: authorId
      })
    } catch (error) {
      console.error(error)
    }
  },

  movePostit(boardId: string, postit: StickyNote) {
    try {
      sendMessage({
        action: Actions.MOVE_POSTIT,
        boardId: boardId,
        id: postit.id,
        posX: postit.posX,
        posY: postit.posY
      })
    } catch (error) {
      console.error(error)
    }
 }
 }
})
