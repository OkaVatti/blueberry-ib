// /apps/frontend/types/user-shared.ts
/**
 * Shared role & permission shapes and helpers for Blueberry frontend.
 * Fixed: ensure no circular references and no type-only symbols used as runtime values.
 */

/* ---------- Core Types ---------- */

/**
 * RoleName: union of known role names used across the application.
 */
export type RoleName =
  | "guest"
  | "anonymous"
  | "n00b13"
  | "user"
  | "moderator"
  | "admin"
  | "co-owner"
  | "owner";

/**
 * PermissionName: canonical list of permission keys used by the frontend
 */
export type PermissionName =
  | "posts.create"
  | "posts.edit"
  | "posts.delete"
  | "posts.moderate"
  | "comments.create"
  | "comments.delete"
  | "comments.moderate"
  | "boards.create"
  | "boards.edit"
  | "boards.delete"
  | "users.ban"
  | "users.promote"
  | "users.edit"
  | "moderation.view"
  | "moderation.take-action"
  | "site.admin"
  | "site.read"
  | "site.write";

/**
 * Permission object model.
 */
export interface Permission {
  id: string;
  name: PermissionName;
  description?: string;
}

/**
 * Role model: small denormalized version frontend can expect.
 */
export interface Role {
  id: string;
  name: RoleName;
  description?: string;
  permissions: Array<PermissionName | Permission>;
}

/* ---------- Default permission & role definitions ---------- */

/**
 * DEFAULT_PERMISSIONS: minimal permission catalog used by the frontend for policy checks / UI.
 * The server is the source of truth; these are fallbacks for client-side checks & dev.
 */
export const DEFAULT_PERMISSIONS: Readonly<Record<PermissionName, Permission>> =
  {
    "posts.create": {
      id: "posts.create",
      name: "posts.create",
      description: "Create new posts",
    },
    "posts.edit": {
      id: "posts.edit",
      name: "posts.edit",
      description: "Edit own posts",
    },
    "posts.delete": {
      id: "posts.delete",
      name: "posts.delete",
      description: "Delete posts",
    },
    "posts.moderate": {
      id: "posts.moderate",
      name: "posts.moderate",
      description: "Moderate posts",
    },

    "comments.create": {
      id: "comments.create",
      name: "comments.create",
      description: "Create comments",
    },
    "comments.delete": {
      id: "comments.delete",
      name: "comments.delete",
      description: "Delete comments",
    },
    "comments.moderate": {
      id: "comments.moderate",
      name: "comments.moderate",
      description: "Moderate comments",
    },

    "boards.create": {
      id: "boards.create",
      name: "boards.create",
      description: "Create new boards",
    },
    "boards.edit": {
      id: "boards.edit",
      name: "boards.edit",
      description: "Edit boards",
    },
    "boards.delete": {
      id: "boards.delete",
      name: "boards.delete",
      description: "Delete boards",
    },

    "users.ban": {
      id: "users.ban",
      name: "users.ban",
      description: "Ban users from site or boards",
    },
    "users.promote": {
      id: "users.promote",
      name: "users.promote",
      description: "Promote users to roles",
    },
    "users.edit": {
      id: "users.edit",
      name: "users.edit",
      description: "Edit user profiles",
    },

    "moderation.view": {
      id: "moderation.view",
      name: "moderation.view",
      description: "View moderation dashboards",
    },
    "moderation.take-action": {
      id: "moderation.take-action",
      name: "moderation.take-action",
      description: "Perform moderation actions (delete/ban)",
    },

    "site.admin": {
      id: "site.admin",
      name: "site.admin",
      description: "Full site administration",
    },
    "site.read": {
      id: "site.read",
      name: "site.read",
      description: "Read site content (general)",
    },
    "site.write": {
      id: "site.write",
      name: "site.write",
      description: "Write content site-wide",
    },
  } as const;

/**
 * DEFAULT_ROLES: frontend fallback definitions.
 * Construct each role individually without referencing DEFAULT_ROLES during build.
 */
export const DEFAULT_ROLES: Readonly<Record<RoleName, Role>> = {
  guest: {
    id: "role_guest",
    name: "guest",
    description: "Unauthenticated visitors",
    permissions: ["site.read"],
  },
  anonymous: {
    id: "role_anonymous",
    name: "anonymous",
    description: "anonymous user from the deep internet... probably a dickhead",
    permissions: ["site.read", "site.write", "posts.create", "comments.create"],
  },
  n00b13: {
      id: "role_noobie",
      name: "n00b13",
      description: "new user, fresh clothes, weird face tho... i hope they're nice! :3",
      permissions: ["site.read", "site.write", "comments.create"]
  },
  user: {
    id: "role_user",
    name: "user",
    description: "Regular authenticated user, also a dick-head...",
    permissions: ["site.read", "site.write", "posts.create", "comments.create"],
  },
  moderator: {
    id: "role_moderator",
    name: "moderator",
    description: "Moderation users with tools, somehow still also manages to be a dick-head",
    permissions: [
      "site.read",
      "site.write",
      "posts.moderate",
      "comments.moderate",
      "moderation.view",
      "moderation.take-action",
      "users.ban",
    ],
  },
  admin: {
    id: "role_admin",
    name: "admin",
    description: "Administrative role with broad control",
    permissions: [
      "site.read",
      "site.write",
      "site.admin",
      "boards.create",
      "boards.edit",
      "boards.delete",
      "users.promote",
      "users.edit",
      "moderation.take-action",
    ],
  },
  // co-owner: include admin permissions plus users.promote (avoids referencing DEFAULT_ROLES)
  "co-owner": {
    id: "role_coowner",
    name: "co-owner",
    description: "Co-owner with near-owner privileges, aka vice supreme loser",
    permissions: [
      // copy admin permissions
      "site.read",
      "site.write",
      "site.admin",
      "boards.create",
      "boards.edit",
      "boards.delete",
      "users.promote",
      "users.edit",
      "moderation.take-action",
      // ensure promote exists
      "users.promote",
    ],
  },
  // owner: grant everything from DEFAULT_PERMISSIONS
  owner: {
    id: "role_owner",
    name: "owner",
    description: "Project/site owner — full control, aka supreme loser",
    permissions: Object.keys(DEFAULT_PERMISSIONS) as PermissionName[],
  },
};

/* ---------- Utilities / helpers ---------- */

/**
 * Normalize permission value (Permission | PermissionName) -> PermissionName
 */
export function permissionNameFrom(
  p: Permission | PermissionName
): PermissionName {
  return typeof p === "string" ? p : p.name;
}

/**
 * Check whether a Role includes a permission.
 */
export function roleHasPermission(
  role: Role,
  perm: PermissionName | Permission
): boolean {
  const permName = permissionNameFrom(perm);
  return role.permissions.some((p) => permissionNameFrom(p) === permName);
}

/**
 * Check whether a list of Role or RoleName includes a specific RoleName
 */
export function rolesInclude(
  roles: Array<Role | RoleName> | undefined | null,
  roleName: RoleName
): boolean {
  if (!roles || roles.length === 0) return false;
  return roles.some((r) =>
    typeof r === "string" ? r === roleName : r.name === roleName
  );
}

/**
 * Check whether a user (roles can be RoleName[] or Role[]) has a permission.
 * Useful in components: `if (userHasPermission(user.roles, 'moderation.take-action')) {...}`
 */
export function userHasPermission(
  roles: Array<Role | RoleName> | undefined | null,
  perm: PermissionName | Permission
): boolean {
  if (!roles || roles.length === 0) return false;
  const permName = permissionNameFrom(perm);
  for (const r of roles) {
    if (typeof r === "string") {
      const fallback = DEFAULT_ROLES[r as RoleName];
      if (fallback && roleHasPermission(fallback, permName)) return true;
    } else {
      if (roleHasPermission(r, permName)) return true;
    }
  }
  return false;
}

/**
 * Convert mixed input into Role[].
 * Accepts Role[], RoleName[], single Role or RoleName.
 */
export function normalizeRoles(
  input: Role | RoleName | Array<Role | RoleName> | undefined | null
): Role[] {
  if (!input) return [];
  const arr = Array.isArray(input) ? input : [input];
  const out: Role[] = [];
  for (const item of arr) {
    if (typeof item === "string") {
      const r = DEFAULT_ROLES[item as RoleName];
      if (r) out.push(r);
      else
        out.push({
          id: `role_${item}`,
          name: item as RoleName,
          description: "Custom role",
          permissions: [],
        });
    } else {
      out.push(item);
    }
  }
  return out;
}

/* ---------- Runtime guards ---------- */

export function isRole(obj: unknown): obj is Role {
  if (!obj || typeof obj !== "object") return false;
  const asAny = obj as Partial<Role>;
  return (
    typeof asAny.id === "string" &&
    typeof asAny.name === "string" &&
    Array.isArray(asAny.permissions)
  );
}

/* ---------- Exports ---------- */
/* Exported values are the constants, types and functions declared above.
   Do NOT export type-only names as runtime values. */
export {
  DEFAULT_PERMISSIONS as DEFAULT_PERMISSIONS_CONST,
  DEFAULT_ROLES as DEFAULT_ROLES_CONST,
};
