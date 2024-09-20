
export enum Actions {
    NEW_BOARD = 'NEW_BOARD', // params name
    RENAME_BOARD = 'RENAME_BOARD', // params id, text
    DELETE_BOARD = 'DELETE_BOARD', // params id

    NEW_POSTIT = 'NEW_POSTIT', // params {}
    UPDATE_CONTENT = 'UPDATE_CONTENT', // params {boardId, id, text}
    MOVE_POSTIT = 'MOVE_POSTIT', // params {boardId,id, from, to}
    DELETE_POSTIT = 'DELETE_POSTIT', // params boardId, id
    ADD_VOTE = 'ADD_VOTE', // params bardId, id
    REMOVE_VOTE = 'REMOVE_VOTE', // params boardId, id
    SHOW_POSTITS = 'SHOW_POSTITS', // params boardId, AuthorId
    HIDE_POSTITS = 'HIDE_POSTITS', // params boardId, AuthorId
}