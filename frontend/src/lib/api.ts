import { SessionState } from './types'

export async function callApi<TBody extends object | undefined>(
  session: SessionState,
  path: string,
  method = 'GET',
  body?: TBody,
): Promise<unknown> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (session.accessToken) headers.Authorization = `Bearer ${session.accessToken}`

  const response = await fetch(`${session.apiBase}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  const payload = response.status !== 204 ? await response.json().catch(() => null) : { message: 'No content' }
  if (!response.ok) {
    throw new Error(JSON.stringify({ status: response.status, payload }, null, 2))
  }
  return payload
}
