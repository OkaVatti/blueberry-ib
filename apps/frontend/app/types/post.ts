// /apps/frontend/types/post.ts
import type { ID } from "./board";

// small media type used across posts/comments/threads
export type MediaType = "image" | "video" | "audio" | "file" | "embed";

export interface MediaItem {
  id: ID;
  url: string;
  type: MediaType;
  mime?: string;
  width?: number;
  height?: number;
  sizeBytes?: number;
  altText?: string;
  blurhash?: string;
}

// lightweight rich text block (expandable later)
export interface ContentBlock {
  type: "text" | "code" | "quote" | "html";
  data: string;
}

// reactions/engagement counts
export interface ReactionCounts {
  likes: number;
  dislikes: number;
  reposts: number;
  quotes: number;
  ink: number; // karma-like numeric score
}

// basic comment type (comments can have media and replies)
export interface Comment {
  id: ID;
  authorId: ID | null; // null for anonymous
  body: string;
  blocks?: ContentBlock[];
  media?: MediaItem[];
  createdAt: string;
  updatedAt?: string;
  parentCommentId?: ID; // threaded comments
  reactionCounts?: Partial<ReactionCounts>;
  reported?: boolean;
  deleted?: boolean;
}

// core Post type
export interface Post {
  id: ID;
  boardId?: ID; // which board (if any)
  threadId?: ID; // which thread (if any)
  authorId: ID | null; // null = anonymous post
  title?: string | null; // optional
  body: string;
  blocks?: ContentBlock[];
  media?: MediaItem[];
  tags?: string[];
  isSticky?: boolean;
  isLocked?: boolean;
  createdAt: string;
  updatedAt?: string;
  reactionCounts: ReactionCounts;
  commentCount: number;
  repostOfPostId?: ID | null; // for simple reposts
  quoteOfPostId?: ID | null; // for quote-repost
  deleted?: boolean;
  moderation?: {
    hidden?: boolean;
    hiddenBy?: ID;
    hiddenReason?: string;
    reports?: number;
  };
}
