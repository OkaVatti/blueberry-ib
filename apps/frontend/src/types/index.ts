export interface User {
  id: string;
  username: string;
  role: "admin" | "moderator" | "user";
  ink: number;
  created_at: string;
  updated_at: string;
}

export interface Board {
  id: string;
  name: string;
  title: string;
  description: string;
  creator_id: string;
  creator: User;
  is_default: boolean;
  is_private: boolean;
  created_at: string;
  updated_at: string;
}

export interface Thread {
  id: string;
  title: string;
  board_id: string;
  board: Board;
  creator_id: string;
  creator: User;
  is_pinned: boolean;
  is_locked: boolean;
  posts?: Post[];
  likes: number;
  dislikes: number;
  views: number;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: string;
  content: string;
  thread_id: string;
  thread?: Thread;
  creator_id: string;
  creator: User;
  parent_id?: string;
  parent?: Post;
  comments?: Post[];
  images: string[];
  videos: string[];
  links: string[];
  likes: number;
  dislikes: number;
  views: number;
  is_anon: boolean;
  created_at: string;
  updated_at: string;
}

export interface Vote {
  id: string;
  user_id: string;
  user: User;
  post_id?: string;
  post?: Post;
  thread_id?: string;
  thread?: Thread;
  type: "like" | "dislike";
  created_at: string;
}

export interface CreateThreadRequest {
  title: string;
  content: string;
  is_anon: boolean;
}

export interface CreatePostRequest {
  content: string;
  parent_id?: string;
  is_anon: boolean;
  images?: string[];
  videos?: string[];
  links?: string[];
}

export interface CreateBoardRequest {
  name: string;
  title: string;
  description: string;
  is_private: boolean;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface ApiError {
  message: string;
  status: number;
}
