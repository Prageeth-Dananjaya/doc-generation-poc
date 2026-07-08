import { useEffect, useMemo, useState } from "react";
import "./App.css";

type ExpensePayload = {
  id: string;
  description: string;
  amount: number;
  participants: string[];
  payments: Record<string, number>;
};

type EventPayload = {
  id: string;
  name: string;
  participants: string[];
  expenses: ExpensePayload[];
};

type BalanceResponse = {
  balances: Record<string, number>;
};

const defaultPayload: EventPayload = {
  id: "event-2",
  name: "Weekend Trip",
  participants: ["alice", "bob", "carol"],
  expenses: [
    {
      id: "expense-1",
      description: "Dinner",
      amount: 120,
      participants: ["alice", "bob", "carol"],
      payments: {
        alice: 120,
      },
    },
  ],
};

const createId = () => {
  return typeof crypto !== "undefined" && "randomUUID" in crypto
    ? crypto.randomUUID()
    : `event-${Date.now()}`;
};

function App() {
  const [payload, setPayload] = useState<EventPayload>(defaultPayload);
  const [result, setResult] = useState<BalanceResponse | null>(null);
  const [events, setEvents] = useState<EventPayload[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const summary = useMemo(() => {
    if (!result) return null;
    return Object.entries(result.balances).sort(([a], [b]) =>
      a.localeCompare(b),
    );
  }, [result]);

  const sanitizePayload = (payload: EventPayload): EventPayload => {
    const participants = payload.participants
      .map((name) => name.trim())
      .filter(Boolean);
    const participantSet = new Set(participants);

    const expenses = payload.expenses.map((expense) => {
      const validParticipants = expense.participants.filter((name) =>
        participantSet.has(name),
      );
      const payments: Record<string, number> = {};

      Object.entries(expense.payments).forEach(([payer, amount]) => {
        if (participantSet.has(payer) && amount > 0) {
          payments[payer] = amount;
        }
      });

      return {
        ...expense,
        id: expense.id || createId(),
        participants: validParticipants,
        payments,
      };
    });

    return {
      ...payload,
      id: payload.id || createId(),
      participants,
      expenses,
    };
  };

  const updateParticipant = (index: number, value: string) => {
    setPayload((current) => {
      const nextParticipants = [...current.participants];
      nextParticipants[index] = value;
      const participantSet = new Set(
        nextParticipants.filter((name) => name.trim() !== ""),
      );

      const nextExpenses = current.expenses.map((expense) => ({
        ...expense,
        participants: expense.participants.filter((name) =>
          participantSet.has(name),
        ),
        payments: Object.fromEntries(
          Object.entries(expense.payments).filter(([payer]) =>
            participantSet.has(payer),
          ),
        ),
      }));

      return {
        ...current,
        participants: nextParticipants,
        expenses: nextExpenses,
      };
    });
  };

  const addParticipant = () => {
    setPayload((current) => ({
      ...current,
      participants: [...current.participants, ""],
    }));
  };

  const removeParticipant = (index: number) => {
    setPayload((current) => {
      const nextParticipants = current.participants.filter(
        (_, participantIndex) => participantIndex !== index,
      );
      const participantSet = new Set(nextParticipants.filter(Boolean));

      const nextExpenses = current.expenses.map((expense) => ({
        ...expense,
        participants: expense.participants.filter((name) =>
          participantSet.has(name),
        ),
        payments: Object.fromEntries(
          Object.entries(expense.payments).filter(([payer]) =>
            participantSet.has(payer),
          ),
        ),
      }));

      return {
        ...current,
        participants: nextParticipants,
        expenses: nextExpenses,
      };
    });
  };

  const updateExpenseField = (
    index: number,
    field: keyof Omit<ExpensePayload, "id" | "payments" | "participants">,
    value: string | number,
  ) => {
    setPayload((current) => {
      const nextExpenses = [...current.expenses];
      nextExpenses[index] = {
        ...nextExpenses[index],
        [field]: field === "amount" ? Number(value) : value,
      } as ExpensePayload;
      return {
        ...current,
        expenses: nextExpenses,
      };
    });
  };

  const toggleExpenseParticipant = (expenseIndex: number, name: string) => {
    setPayload((current) => {
      const nextExpenses = [...current.expenses];
      const expense = nextExpenses[expenseIndex];
      const hasParticipant = expense.participants.includes(name);
      nextExpenses[expenseIndex] = {
        ...expense,
        participants: hasParticipant
          ? expense.participants.filter((participant) => participant !== name)
          : [...expense.participants, name],
      };
      return {
        ...current,
        expenses: nextExpenses,
      };
    });
  };

  const updateExpensePayment = (
    expenseIndex: number,
    payer: string,
    value: number,
  ) => {
    setPayload((current) => {
      const nextExpenses = [...current.expenses];
      const expense = nextExpenses[expenseIndex];
      nextExpenses[expenseIndex] = {
        ...expense,
        payments: {
          ...expense.payments,
          [payer]: Number(value),
        },
      };
      return {
        ...current,
        expenses: nextExpenses,
      };
    });
  };

  const addExpense = () => {
    setPayload((current) => ({
      ...current,
      expenses: [
        ...current.expenses,
        {
          id: createId(),
          description: "",
          amount: 0,
          participants: [...current.participants],
          payments: {},
        },
      ],
    }));
  };

  const removeExpense = (index: number) => {
    setPayload((current) => ({
      ...current,
      expenses: current.expenses.filter(
        (_, expenseIndex) => expenseIndex !== index,
      ),
    }));
  };

  const submitEvent = async () => {
    setIsLoading(true);
    setError("");

    try {
      const event = sanitizePayload(payload);
      const response = await fetch("/api/events", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(event),
      });

      if (!response.ok) {
        const body = await response.json();
        throw new Error(body.error || "Unable to create event");
      }

      const createdEvent = (await response.json()) as EventPayload;
      const balancesResponse = await fetch(
        `/api/events/${createdEvent.id}/balances`,
      );
      if (!balancesResponse.ok) {
        const body = await balancesResponse.json();
        throw new Error(body.error || "Unable to calculate balances");
      }

      const balanceData = (await balancesResponse.json()) as BalanceResponse;
      setResult(balanceData);
      setPayload(createdEvent);
      await loadEvents();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setIsLoading(false);
    }
  };

  const loadEvents = async () => {
    try {
      const res = await fetch("/api/events");
      if (!res.ok) return;
      const data = await res.json();
      setEvents(data.events || []);
    } catch (_) {
      // ignore
    }
  };

  useEffect(() => {
    loadEvents();
  }, []);

  const loadEvent = async (id: string) => {
    try {
      const res = await fetch(`/api/events/${id}`);
      if (!res.ok) throw new Error("failed to load event");
      const evt = (await res.json()) as EventPayload;
      setPayload(evt);

      const balancesRes = await fetch(`/api/events/${id}/balances`);
      if (balancesRes.ok) {
        const bd = (await balancesRes.json()) as BalanceResponse;
        setResult(bd);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "load failed");
    }
  };

  const deleteEvent = async (id: string) => {
    try {
      const res = await fetch(`/api/events/${id}`, { method: "DELETE" });
      if (!res.ok) throw new Error("delete failed");
      await loadEvents();
      // if currently viewing deleted event, clear
      if (payload.id === id) {
        setPayload({ ...defaultPayload, id: createId() });
        setResult(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "delete failed");
    }
  };

  const newEvent = () => {
    setPayload({ id: createId(), name: "", participants: [""], expenses: [] });
    setResult(null);
  };

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="subsection-header">
          <h3>Events</h3>
          <button type="button" className="secondary-button" onClick={newEvent}>
            New
          </button>
        </div>
        <div className="event-list">
          {events.map((evt) => (
            <div key={evt.id} className="event-item">
              <button className="text-button" onClick={() => loadEvent(evt.id)}>
                {evt.name || evt.id}
              </button>
              <button
                className="text-button"
                onClick={() => deleteEvent(evt.id)}
              >
                Delete
              </button>
            </div>
          ))}
        </div>
      </aside>
      <section className="hero-card">
        <div className="hero-copy">
          <p className="eyebrow">Collaborative expense insights</p>
          <h1>Split group costs fairly, without the spreadsheet headache.</h1>
          <p className="lede">
            Build an event with participants and expenses, then calculate fair
            balances instantly.
          </p>
        </div>

        <div className="panel-card">
          <label className="field">
            <span>Event name</span>
            <input
              value={payload.name}
              onChange={(e) => setPayload({ ...payload, name: e.target.value })}
              placeholder="Weekend Trip"
            />
          </label>

          <div className="subsection">
            <div className="subsection-header">
              <h3>Participants</h3>
              <button
                type="button"
                className="secondary-button"
                onClick={addParticipant}
              >
                Add participant
              </button>
            </div>
            {payload.participants.map((participant, index) => (
              <div key={`${participant}-${index}`} className="list-row">
                <input
                  value={participant}
                  onChange={(e) => updateParticipant(index, e.target.value)}
                  placeholder={index === 0 ? "alice" : "participant name"}
                />
                <button
                  type="button"
                  className="text-button"
                  onClick={() => removeParticipant(index)}
                >
                  Remove
                </button>
              </div>
            ))}
          </div>

          <div className="subsection">
            <div className="subsection-header">
              <h3>Expenses</h3>
              <button
                type="button"
                className="secondary-button"
                onClick={addExpense}
              >
                Add expense
              </button>
            </div>
            {payload.expenses.map((expense, expenseIndex) => (
              <div key={expense.id} className="expense-card">
                <div className="expense-header">
                  <h4>Expense {expenseIndex + 1}</h4>
                  <button
                    type="button"
                    className="text-button"
                    onClick={() => removeExpense(expenseIndex)}
                  >
                    Remove
                  </button>
                </div>
                <label className="field">
                  <span>Description</span>
                  <input
                    value={expense.description}
                    onChange={(e) =>
                      updateExpenseField(
                        expenseIndex,
                        "description",
                        e.target.value,
                      )
                    }
                    placeholder="Dinner, taxi, supplies…"
                  />
                </label>
                <label className="field">
                  <span>Amount</span>
                  <input
                    type="number"
                    min="0"
                    value={expense.amount}
                    onChange={(e) =>
                      updateExpenseField(
                        expenseIndex,
                        "amount",
                        Number(e.target.value),
                      )
                    }
                  />
                </label>
                <div className="expense-choices">
                  <span>Shared by</span>
                  <div className="participant-grid">
                    {payload.participants.map((participant) => (
                      <label
                        key={`${expense.id}-${participant}`}
                        className="checkbox-label"
                      >
                        <input
                          type="checkbox"
                          checked={expense.participants.includes(participant)}
                          onChange={() =>
                            toggleExpenseParticipant(expenseIndex, participant)
                          }
                        />
                        {participant || "Unnamed"}
                      </label>
                    ))}
                  </div>
                </div>
                <div className="expense-choices">
                  <span>Payments</span>
                  <div className="participant-grid">
                    {payload.participants.map((participant) => (
                      <label
                        key={`pay-${expense.id}-${participant}`}
                        className="field small-field"
                      >
                        <span>{participant || "Unnamed"}</span>
                        <input
                          type="number"
                          min="0"
                          value={expense.payments[participant] ?? 0}
                          onChange={(e) =>
                            updateExpensePayment(
                              expenseIndex,
                              participant,
                              Number(e.target.value),
                            )
                          }
                        />
                      </label>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>

          <button type="button" onClick={submitEvent} disabled={isLoading}>
            {isLoading ? "Calculating…" : "Calculate balances"}
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
                <p className={amount >= 0 ? "positive" : "negative"}>
                  {amount >= 0
                    ? `Should receive ${amount.toFixed(2)}`
                    : `Owes ${Math.abs(amount).toFixed(2)}`}
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
  );
}

export default App;
