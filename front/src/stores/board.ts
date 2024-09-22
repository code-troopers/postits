import { Board } from './../models/Board';
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { sendMessage } from '../services/websocketService'
import type { Message } from '@/models/Message'
import { Actions } from '../Actions'
import axios from 'axios';
import type { StickyNote } from '@/models/StickyNote';

const API_URL = import.meta.env.VITE_API_URL;

export const useBoardStore = defineStore('board', () => {
  const boards = ref([] as Board[])
  function onMessage(message: Message) {
    let board: Board | undefined;
    switch (message.action) {
      case Actions.NEW_BOARD:
        board = new Board();
        board.id = message.id;
        board.name = message.text;
        boards.value.push(board)
        break
      case Actions.RENAME_BOARD:
        board = boards.value.find((b) => b.id === message.id)
        if (board) {
          board.name = message.text
        }
        break
      case Actions.DELETE_BOARD:
        boards.value = boards.value.filter((b) => b.id !== message.boardId)
        break
      case Actions.NEW_POSTIT:
        board = boards.value.find((b) => b.id === message.boardId)
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
            author: {
              id: message.authorId
            }
          })
        }
        break
      case Actions.UPDATE_CONTENT:
        board = boards.value.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.text = message.text
          }
        }
        break

      case Actions.MOVE_POSTIT:
        board = boards.value.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.posX = message.posX
            postit.posY = message.posY
          }
        }
        break

      case Actions.DELETE_POSTIT:
        board = boards.value.find((b) => b.id === message.boardId)
        if (board) {
          board.postits = board.postits.filter((p) => p.id !== message.id)
        }
        break

      case Actions.ADD_VOTE:
        board = boards.value.find((b) => b.id === message.boardId)
        if (board) {
          const postit = board.postits.find((p) => p.id === message.id)
          if (postit) {
            postit.votes++
          }
        }
        break

      case Actions.REMOVE_VOTE:
        board = boards.value.find((b) => b.id === message.boardId)
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
          initPostits(message.boardId)
        }
        break

    }
  }
  function newBoard() {
    try {
      sendMessage({
        action: Actions.NEW_BOARD,
        text: `Board ${boards.value.length + 1}`,
      })
    } catch (error) {
      console.error(error)
    }
  }
  function renameBoard(boardId: string, text: string) {
    try {
      sendMessage({
        action: Actions.RENAME_BOARD,
        boardId: boardId,
        text: text
      })
    } catch (error) {
      console.error(error)
    }

  }
  function deleteBoard(boardId: string | undefined) {
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
  }
  function newPostit(boardId: string, x: number, y: number) {
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
  }
  function updateContent(boardId: string, id: string, text: string) {
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
  }
  function deletePostit(boardId: string, id: string | undefined) {
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
  }

  function getBoards() {
    return axios.get(`${API_URL}/api/boards`).then((response: any) => {
      boards.value = response.data
    })
  }

  function getPostits(boardId: string) {
    const board = boards.value.find((b) => b.id === boardId)
    if (board) {
      return board.postits
    }
    return []
  }

  async function initPostits(boardId: string) {
    if (boards.value.length === 0) {
      await getBoards()
    }
    const board = boards.value.find((b) => b.id === boardId)
    if (board) {
      axios.get(`${API_URL}/api/boards/${boardId}/postits`).then((response: any) => {
        board.postits = response.data
      })
    }
  }

  function addVote(id: string, boardId: string) {
    try {
      sendMessage({
        action: Actions.ADD_VOTE,
        id: id,
        boardId: boardId
      })
    } catch (error) {
      console.error(error)
    }
  }

  function removeVote(id: string, boardId: string) {
    try {
      sendMessage({
        action: Actions.REMOVE_VOTE,
        id: id,
        boardId: boardId
      })
    } catch (error) {
      console.error(error)
    }
  }

  function showPostits(boardId: string, authorId: string) {
    try {
      sendMessage({
        action: Actions.SHOW_POSTITS,
        boardId: boardId,
        authorId: authorId
      })
    } catch (error) {
      console.error(error)
    }
  }

  function hidePostits(boardId: string, authorId: string) {
    try {
      sendMessage({
        action: Actions.HIDE_POSTITS,
        boardId: boardId,
        authorId: authorId
      })
    } catch (error) {
      console.error(error)
    }
  }

  function movePostit(boardId: string, postit: StickyNote) {
    try {
      sendMessage({
        action: Actions.MOVE_POSTIT,
        boardId: postit.board?.id,
        id: postit.id,
        posX: postit.posX,
        posY: postit.posY
      })
    } catch (error) {
      console.error(error)
    }
  }

  return { boards, onMessage, newBoard, renameBoard,
     deleteBoard, newPostit, updateContent, deletePostit, getPostits, getBoards, initPostits,
      addVote, removeVote, showPostits, hidePostits, movePostit }
})
