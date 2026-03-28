export type RouteKey = 'auth' | 'workspace' | 'tasks' | 'collab' | 'admin'

export interface SessionState {
  accessToken: string
  refreshToken: string
  apiBase: string
}

export interface ApiResult {
  status: number
  payload: unknown
}
