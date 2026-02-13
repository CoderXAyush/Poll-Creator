export default function ResultsChart({ results, votedOptionId, colors }) {
  const maxVotes = Math.max(...results.options.map((o) => o.votes), 1);
  const winnerVotes = Math.max(...results.options.map((o) => o.votes));

  return (
    <div className="results-container">
      <h3>Results</h3>

      {results.options.map((opt, i) => {
        const pct = results.totalVotes > 0
          ? ((opt.votes / results.totalVotes) * 100).toFixed(1)
          : "0.0";
        const isWinner = opt.votes === winnerVotes && opt.votes > 0;
        const isVotedFor = votedOptionId === opt.id;
        const color = colors[i % colors.length];

        return (
          <div key={opt.id} className="result-bar-wrapper">
            <div className="result-bar-header">
              <span className={`result-bar-label ${isWinner ? "winner" : ""}`}>
                {opt.text}
                {isVotedFor && " (Your vote)"}
              </span>
              <span className="result-bar-stats">
                {pct}% · {opt.votes} vote{opt.votes !== 1 ? "s" : ""}
              </span>
            </div>
            <div className="result-bar-track">
              <div
                className="result-bar-fill"
                style={{
                  width: `${Math.max(parseFloat(pct), 1)}%`,
                  background: color,
                }}
              />
            </div>
          </div>
        );
      })}

      <div className="result-total">
        <strong>{results.totalVotes}</strong> total vote{results.totalVotes !== 1 ? "s" : ""}
      </div>
    </div>
  );
}
