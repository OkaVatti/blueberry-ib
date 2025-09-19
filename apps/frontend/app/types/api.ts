// strong API types used by the frontend stores
export interface ApiError {
  message: string;
  status?: number;
  code?: string;
  details?: string;
}

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: ApiError };

export interface Paginated<T> {
  items: T[];
  nextCursor?: string | null;
  prevCursor?: string | null;
  total?: number;
}
