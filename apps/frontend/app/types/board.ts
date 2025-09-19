// /apps/frontend/types/board.ts
export type ID = string;

export type BoardVisibility = "public" | "restricted" | "private";

export interface BoardSettings {
  allowAnonymousPosts?: boolean;
  allowMedia?: boolean;
  maxMediaAttachments?: number;
  allowReposts?: boolean;
  defaultSort?: "new" | "hot" | "top" | "active";
  nsfw?: boolean;
  customTheme?: string | null;
}

export interface BoardStats {
  postsCount: number;
  threadsCount: number;
  activeUsers?: number;
  dailyAvgPosts?: number;
}

export interface BoardSummary {
  id: ID;
  name: string;
  slug: string;
  description?: string;
  visibility: BoardVisibility;
  createdAt: string;
  updatedAt?: string;
  stats?: Partial<BoardStats>;
  tags?: string[];
  pinnedThreadIds?: ID[];
  creatorId?: ID;
}

export interface Board extends BoardSummary {
  settings: BoardSettings;
  moderators?: ID[]; // user ids
  admins?: ID[]; // user ids
  rules?: string[];
  featuredMedia?: string[]; // urls
}
