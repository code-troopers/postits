import type { User } from "./User";

export interface Message {
  action: string;
  boardId?: string;
  authorId?: string;
  id?: string;
  text?: string;
  posX?: number;
  posY?: number;
  token?: string;
  author?: User;
  weight?: number;
}
