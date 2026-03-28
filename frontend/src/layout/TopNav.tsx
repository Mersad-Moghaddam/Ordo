import { NavigateFunction } from 'react-router-dom'

interface TopNavProps {
  apiBase: string
  onApiBaseChange: (value: string) => void
  onNavigate: NavigateFunction
}

export function TopNav({ apiBase, onApiBaseChange, onNavigate }: TopNavProps) {
  return (
    <header className="top-nav">
      <div>
        <h1>Ordo Command Center</h1>
        <p>Jira-style workflow pages: Login → Create Account → Main operations.</p>
      </div>
      <label className="field api-field">
        <span>API Base</span>
        <input value={apiBase} onChange={(event) => onApiBaseChange(event.target.value)} />
      </label>
      <nav>
        <button className="tab" onClick={() => onNavigate('/login')}>Login</button>
        <button className="tab" onClick={() => onNavigate('/signup')}>Create Account</button>
        <button className="tab" onClick={() => onNavigate('/main')}>Main Page</button>
      </nav>
    </header>
  )
}
