import { useState, useEffect } from "react";
import { getPolls } from "../api";

export default function PollList({ onSelect }) {
  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getPolls()
      .then(setPolls)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner" />
        <p>Loading polls...</p>
      </div>
    );
  }

  if (polls.length === 0) {
    return (
      <div className="empty-state">
        <p>No polls yet. Create the first one.</p>
      </div>
    );
  }

  return (
    <div>
      {polls.map((poll) => (
        <div
          key={poll.id}
          className="poll-list-item"
          onClick={() => onSelect(poll)}
        >
          <h3>{poll.question}</h3>
          <div className="poll-meta">
            <span>
              <span className={`badge ${poll.closed ? "badge-closed" : "badge-open"}`}>
                {poll.closed ? "Closed" : "Open"}
              </span>
            </span>
            <span>{poll.totalVotes} votes</span>
            <span>{poll.optionCount} options</span>
            <span>
              {new Date(poll.createdAt).toLocaleDateString()}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}
