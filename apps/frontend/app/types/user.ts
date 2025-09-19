// /apps/frontend/types/user.ts
import type { Profile } from "./profile";
import type { Permission, Role } from ".//user-shared"; // optional; see notes below
import type { ID } from "./board";

/**
 * Compatibility layer: previously some code used `UserProfile`.
 * Re-export Profile as UserProfile so existing imports work.
 */
export type UserProfile = Profile;

/**
 * RoleName and permission shapes.
 * (If you already have a role/permission file, adapt imports instead.)
 */
export type RoleName =
  | "user"
  | "moderator"
  | "admin"
  | "owner"
  | "co-owner"
  | "guest";

/**
 * Full User object for frontend usage.
 * Note: sensitive server-only properties (password hashes, reset tokens) MUST NOT be present here.
 */
export interface User extends Profile {
  // contact (email kept optional on the client; server never exposes password)
  email?: string | null;

  // role & permission information (can be names or full objects)
  roles?: RoleName[] | Role[];
  permissions?: Permission[];

  // last-seen and minimal metadata
  lastSeenAt?: string | null;
  isStaff?: boolean;
  // precise stats may be embedded in Profile.stats, but keep top-level ink for convenience
  ink?: number;
}

/**
 * Tokens returned from auth endpoints (frontend shape).
 * Access token short-lived; refresh token longer; server must support safe rotation + revocation.
 */
export interface AuthTokens {
  accessToken: string;
  refreshToken?: string;
  tokenType?: "Bearer";
  expiresAt?: string; // ISO timestamp when accessToken expires
}

/**
 * Signup DTO - used only when sending to backend to create an account.
 * Never persist raw password on client beyond transmission to server (use TLS).
 */
export interface SignupDTO {
  displayName: string;
  username: string;
  email: string;
  password: string;
}
