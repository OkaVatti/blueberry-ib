// /apps/frontend/types/user.ts
import type { ID } from "./board";

export type RoleName =
  | "user"
  | "moderator"
  | "admin"
  | "owner"
  | "co-owner"
  | "guest";

export interface Permission {
  id: string;
  name: string;
  description?: string;
}

export interface Role {
  id: string;
  name: RoleName;
  permissions?: Permission[];
}

export interface UserProfile {
  id: ID;
  displayName: string;
  username: string;
  bio?: string;
  avatarUrl?: string;
  bannerUrl?: string;
  location?: string | null;
  website?: string | null;
  pronouns?: string | null;
  createdAt: string;
  updatedAt?: string;
  isVerified?: boolean;
  customTheme?: string | null;
}

export interface UserStats {
  ink: number;
  posts: number;
  comments: number;
  likes: number;
  dislikes: number;
  reposts: number;
  reputation?: number;
}

export interface User extends UserProfile {
  email?: string;
  roles?: RoleName[] | Role[];
  permissions?: Permission[];
  isActive?: boolean;
  lastSeenAt?: string | null;
  stats?: Partial<UserStats>;
  // Never store raw passwords in frontend types; only for DTOs during signup:
  // password?: string // only in signup DTOs and never persisted
}

export interface AuthTokens {
  accessToken: string;
  refreshToken?: string;
  expiresAt?: string;
}

export interface SignupDTO {
  displayName: string;
  username: string;
  email: string;
  password: string; // plaintext only when sending to backend to be hashed (argon2 server-side)
}
