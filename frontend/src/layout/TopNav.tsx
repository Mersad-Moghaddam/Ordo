import { RouteKey } from '../lib/types'

interface TopNavProps {
  active: RouteKey
  onSwitch: (route: RouteKey) => void
  apiBase: string
  onApiBaseChange: (value: string) => void
}

const tabs: Array<{ key: RouteKey; label: string }> = [
  { key: 'auth', label: 'Auth' },
  { key: 'workspace', label: 'Workspace' },
  { key: 'tasks', label: 'Tasks' },
  { key: 'collab', label: 'Collaboration' },
  { key: 'admin', label: 'Admin' },
]

export function TopNav({ active, onSwitch, apiBase, onApiBaseChange }: TopNavProps) {
  return (
    <header className="top-nav">
      <div>
        <h1>Ordo Command Center</h1>
        <p>Professional control panel mapped directly to backend APIs.</p>
      </div>
      <label className="field api-field">
        <span>API Base</span>
        <input value={apiBase} onChange={(event) => onApiBaseChange(event.target.value)} />
      </label>
      <nav>
        {tabs.map((tab) => (
          <button key={tab.key} className={tab.key === active ? 'tab active' : 'tab'} onClick={() => onSwitch(tab.key)}>
            {tab.label}
          </button>
        ))}
      </nav>
    </header>
  )
}
