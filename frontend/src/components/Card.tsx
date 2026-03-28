import { PropsWithChildren } from 'react'

interface CardProps extends PropsWithChildren {
  title: string
  subtitle?: string
}

export function Card({ title, subtitle, children }: CardProps) {
  return (
    <section className="card">
      <header className="card-head">
        <h3>{title}</h3>
        {subtitle ? <p>{subtitle}</p> : null}
      </header>
      {children}
    </section>
  )
}
