import type { StickyNote } from "./StickyNote";

export class Board {
    id?: string;
    name?: string;
    postits: StickyNote[] = [];
}