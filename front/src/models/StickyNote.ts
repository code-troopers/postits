import type { Board } from "./Board";
import type { User } from "./User";

export class StickyNote {
    id?: string;
    author?: User;
    text?: string;
    board?: Board;
    posX?: number;
    posY?: number;
    show: boolean = false;
    votes: number = 0;
}