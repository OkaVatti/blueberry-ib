// /apps/frontend/types/thread.ts
import type { ID } from "./board";
import type { Post } from "./post";

export interface ThreadSummary {
  id: ID;
  boardId?: ID;
  title?: string | null;
  starterPostId: ID;
  createdAt: string;
  updatedAt?: string;
  postCount: number;
  bumpAt?: string; // last bump timestamp
  pinned?: boolean;
  closed?: boolean;
  sticky?: boolean;
  tags?: string[];
  previewMedia?: string[]; // urls
}

export interface Thread extends ThreadSummary {
  posts?: Post[]; // Optional: the server can deliver the posts inline
  starter?: Post; // optional denormalized starter post
  participants?: ID[]; // users who replied
  moderation?: {
    reported?: boolean;
    reportCount?: number;
  };
}
