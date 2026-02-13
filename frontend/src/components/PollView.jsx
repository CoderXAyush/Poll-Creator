import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getPoll, submitVote, getPollResults, closePoll } from "../api";
import ResultsChart from "./ResultsChart";

const BAR_COLORS = [
  "var(--bar-1)",
  "var(--bar-2)",
  "var(--bar-3)",
  "var(--bar-4)",
  "var(--bar-5)",
  "var(--bar-6)",
  "var(--bar-7)",
  "var(--bar-8)",
];

export default function PollView() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [poll, setPoll] = useState(null);
  const [results, setResults] = useState(null);
  const [selectedOption, setSelectedOption] = useState(null);
  const [hasVoted, setHasVoted] = useState(false);
  const [votedOptionId, setVotedOptionId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [voting, setVoting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [showResults, setShowResults] = useState(false);
  const [copied, setCopied] = useState(false);

  // Fetch poll data
  useEffect(() => {
    fetchPoll();
  }, [id]);

  async function fetchPoll() {
    setLoading(true);
    try {
      const data = await getPoll(id);
      setPoll(data);
      setHasVoted(data.hasVoted);
      setVotedOptionId(data.votedOptionId);
      if (data.hasVoted || data.closed) {
        const resData = await getPollResults(id);
        setResults(resData);
        setShowResults(true);
      }
    } catch {
      setError("Poll not found");
    } finally {
      setLoading(false);
    }
  }

  async function handleVote() {
    if (selectedOption === null) return;
    setVoting(true);
    setError("");
    try {
      const data = await submitVote(id, selectedOption);
      setPoll(data);
      setHasVoted(true);
      setVotedOptionId(selectedOption);
      setSuccess("Vote submitted successfully!");
      // Fetch results
      const resData = await getPollResults(id);
      setResults(resData);
      setShowResults(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setVoting(false);
    }
  }

  async function handleClose() {
    try {
      await closePoll(id);
      fetchPoll();
    } catch {
      setError("Failed to close poll");
    }
  }

  async function handleViewResults() {
    const resData = await getPollResults(id);
    setResults(resData);
    setShowResults(true);
  }

  function handleCopyLink() {
    const url = `${window.location.origin}/poll/${id}`;
    navigator.clipboard.writeText(url).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner" />
        <p>Loading poll...</p>
      </div>
    );
  }

  if (error && !poll) {
    return (
      <div className="card">
        <div className="alert alert-error">⚠ {error}</div>
        <button className="btn btn-secondary" onClick={() => navigate("/")}>
          Back
        </button>
      </div>
    );
  }

  return (
    <div>
      <button className="back-link" onClick={() => navigate("/")}>
        ← Back
      </button>

      <div className="card">
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "start", marginBottom: 16 }}>
          <h2 style={{ margin: 0 }}>{poll.question}</h2>
          <span className={`badge ${poll.closed ? "badge-closed" : "badge-open"}`}>
            {poll.closed ? "Closed" : "Open"}
          </span>
        </div>

        {error && <div className="alert alert-error">{error}</div>}
        {success && <div className="alert alert-success">{success}</div>}

        {/* Voting Interface */}
        {!hasVoted && !poll.closed && (
          <>
            <div style={{ marginBottom: 16 }}>
              {poll.options.map((opt) => (
                <div
                  key={opt.id}
                  className={`vote-option ${selectedOption === opt.id ? "selected" : ""}`}
                  onClick={() => setSelectedOption(opt.id)}
                >
                  <div className="vote-radio">
                    <div className="vote-radio-inner" />
                  </div>
                  <span className="vote-option-text">{opt.text}</span>
                </div>
              ))}
            </div>

            <div className="btn-group">
              <button
                className="btn btn-primary"
                onClick={handleVote}
                disabled={selectedOption === null || voting}
              >
                {voting ? "Submitting..." : "Submit Vote"}
              </button>
              <button className="btn btn-secondary" onClick={handleViewResults}>
                View Results
              </button>
            </div>
          </>
        )}

        {/* Already voted or closed — show results */}
        {showResults && results && (
          <ResultsChart results={results} votedOptionId={votedOptionId} colors={BAR_COLORS} />
        )}

        {/* Actions */}
        <div className="btn-group" style={{ marginTop: 20 }}>
          {!poll.closed && (
            <button className="btn btn-danger btn-sm" onClick={handleClose}>
              Close Poll
            </button>
          )}
          {hasVoted && !showResults && (
            <button className="btn btn-secondary btn-sm" onClick={handleViewResults}>
              View Results
            </button>
          )}
        </div>

        {/* Share */}
        <div className="share-box">
          <label>Share this poll</label>
          <div className="share-url">
            <input
              readOnly
              value={`${window.location.origin}/poll/${id}`}
              onFocus={(e) => e.target.select()}
            />
            <button className="btn btn-success btn-sm" onClick={handleCopyLink}>
              {copied ? "Copied!" : "Copy"}
            </button>
          </div>
          {copied && <span className="copied-toast">Link copied to clipboard!</span>}
        </div>
      </div>
    </div>
  );
}
