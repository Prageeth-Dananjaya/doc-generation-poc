import { useMemo, useState } from 'react'
import './App.css'

type BalanceResponse = {
  balances: Record<string, number>
}

const defaultPayload = {
  id: 'event-2',
  name: 'Weekend Trip',
  participants: ['alice', 'bob', 'carol'],
  expenses: [
    {
      description: 'Dinner',
      amount: 120,
      participants: ['alice', 'bob', 'carol'],
      payments: {
        alice: 120,
      },
    },
  ],
}

function App() {
  const [payload, setPayload] = useState(defaultPayload)
  const [result, setResult] = useState<BalanceResponse | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const summary = useMemo(() => {
    if (!result) return null
    return Object.entries(result.balances).sort(([a], [b]) => a.localeCompare(b))
  }, [result])

  const submitEvent = async () => {
    setIsLoading(true)
    setError('')

    try {
      const response = await fetch('/api/event', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      if (!response.ok) {
        throw new Error('Unable to calculate balances')
      }

      const data = (await response.json()) as BalanceResponse
      setResult(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <main className="app-shell">
      <section className="hero-card">
        <div className="hero-copy">
          <p className="eyebrow">Collaborative expense insights</p>
          <h1>Split group costs fairly, without the spreadsheet headache.</h1>
          <p className="lede">
            See who owes what in seconds when one or more people already covered the bill.
          </p>
        </div>

        <div className="panel-card">
          <label className="field">
            <span>Event name</span>
            <input
              value={payload.name}
              onChange={(e) => setPayload({ ...payload, name: e.target.value })}
            />
          </label>

          <label className="field">
            <span>Participants</span>
            <input
              value={payload.participants.join(', ')}
              onChange={(e) =>
                setPayload({
                  ...payload,
                  participants: e.target.value.split(',').map((item) => item.trim()).filter(Boolean),
                })
              }
            />
          </label>

          <label className="field">
            <span>Expense amount</span>
            <input
              type="number"
              value={payload.expenses[0].amount}
              onChange={(e) =>
                setPayload({
                  ...payload,
                  expenses: [
                    {
                      ...payload.expenses[0],
                      amount: Number(e.target.value),
                    },
                  ],
                })
              }
            />
          </label>

          <button type="button" onClick={submitEvent} disabled={isLoading}>
            {isLoading ? 'Calculating…' : 'Calculate balances'}
          </button>

          {error ? <p className="error">{error}</p> : null}
        </div>
      </section>

      <section className="results-card">
        <div className="results-header">
          <h2>Settlement view</h2>
          <p>Balances update instantly from the backend service.</p>
        </div>

        {summary ? (
          <div className="balance-grid">
            {summary.map(([name, amount]) => (
              <article key={name} className="balance-item">
                <h3>{name}</h3>
                <p className={amount >= 0 ? 'positive' : 'negative'}>
                  {amount >= 0 ? `Should receive ${amount.toFixed(2)}` : `Owes ${Math.abs(amount).toFixed(2)}`}
                </p>
              </article>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <p>Submit an event to see who should pay or receive money.</p>
          </div>
        )}
      </section>
    </main>
  )
}

export default App
