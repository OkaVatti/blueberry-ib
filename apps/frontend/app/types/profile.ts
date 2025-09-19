// /apps/frontend/types/profile.ts
import type { ID } from "./board";

/**
 * Appearance & theming a user may configure.
 * Keep these small and serializable (no functions).
 */
export interface ProfileAppearance {
  theme?: string; // a named theme (e.g. 'fern', 'kitten', 'cafe', ...)
  accentColor?: string; // hex or CSS color token (e.g. '#1e90ff')
  compactMode?: boolean;
  customCss?: string | null; // optional user CSS (boxed server-side sanitization REQUIRED)
}

/**
 * Privacy options for profile visibility and contact.
 * Values chosen to be explicit and safe for policy enforcement.
 */
export type LastSeenVisibility = "everyone" | "followers" | "no-one";
export type ProfileVisibility = "public" | "unlisted" | "private";

export interface ProfilePrivacySettings {
  profileVisibility?: ProfileVisibility;
  showEmail?: boolean;
  showLastSeen?: LastSeenVisibility;
  allowDirectMessages?: boolean;
  allowMentions?: boolean;
  indexingAllowed?: boolean; // whether search engines / public indexers allowed
  dataExportAllowed?: boolean; // GDPR-style opt-in for exports (server enforces)
}

/**
 * Notification preferences -- keep boolean flags for server to honor
 */
export interface ProfileNotificationSettings {
  emailOnMention?: boolean;
  emailOnFollow?: boolean;
  pushOnMention?: boolean;
  pushOnDirectMessage?: boolean;
  digestFrequency?: "instant" | "hourly" | "daily" | "none";
}

/**
 * Trackable profile statistics (derived server-side).
 * These are read-only on the client.
 */
export interface ProfileStats {
  posts: number;
  comments: number;
  likes: number;
  dislikes: number;
  reposts: number;
  ink: number; // "karma"-like score
  followers?: number;
  following?: number;
  createdAt?: string;
}

/**
 * Core Profile shape used across the frontend.
 * This intentionally excludes sensitive server-only fields (passwordHash, email verification tokens, etc).
 */
export interface Profile {
  id: ID;
  displayName: string;
  username: string;
  bio?: string | null;
  avatarUrl?: string | null;
  bannerUrl?: string | null;
  pronouns?: string | null;
  website?: string | null;
  location?: string | null;

  // preferences & settings
  appearance?: ProfileAppearance;
  privacy?: ProfilePrivacySettings;
  notifications?: ProfileNotificationSettings;

  // derived / readonly
  stats?: ProfileStats;
  isVerified?: boolean;
  isActive?: boolean;

  // timestamps
  createdAt: string;
  updatedAt?: string | null;
}
