
export interface Message {
  action: string;
  boardId?: string;
  authorId?: string;
  id?: string;
  text?: string;
  posX?: number;
  posY?: number;
  token?: string;
}
