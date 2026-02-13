import { useState } from "react";
import { createPoll } from "../api";

export default function CreatePoll({ onCreated }) {
  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const addOption = () => {
    if (options.length < 8) setOptions([...options, ""]);
  };

  const removeOption = (index) => {
    if (options.length > 2) {
      setOptions(options.filter((_, i) => i !== index));
    }
  };

  const updateOption = (index, value) => {
    const updated = [...options];
    updated[index] = value;
    setOptions(updated);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    const trimmedQ = question.trim();
    const trimmedOpts = options.map((o) => o.trim()).filter(Boolean);

    if (!trimmedQ) {
      setError("Please enter a question.");
      return;
    }
    if (trimmedOpts.length < 2) {
      setError("At least 2 options are required.");
      return;
    }

    setLoading(true);
    try {
      const poll = await createPoll(trimmedQ, trimmedOpts);
      setQuestion("");
      setOptions(["", ""]);
      onCreated(poll);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card">
      <h2>Create a New Poll</h2>

      {error && <div className="alert alert-error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Question</label>
          <input
            type="text"
            placeholder="e.g. What's your favorite programming language?"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            autoFocus
          />
        </div>

        <div className="form-group">
          <label>Options</label>
          <div className="options-list">
            {options.map((opt, i) => (
              <div key={i} className="option-row">
                <span className="option-number">{i + 1}</span>
                <input
                  type="text"
                  placeholder={`Option ${i + 1}`}
                  value={opt}
                  onChange={(e) => updateOption(i, e.target.value)}
                />
                {options.length > 2 && (
                  <button
                    type="button"
                    className="remove-btn"
                    onClick={() => removeOption(i)}
                    title="Remove option"
                  >
                    ×
                  </button>
                )}
              </div>
            ))}
          </div>

          {options.length < 8 && (
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              onClick={addOption}
            >
              + Add Option
            </button>
          )}
        </div>

        <button
          type="submit"
          className="btn btn-primary btn-block"
          disabled={loading}
        >
          {loading ? "Creating..." : "Create Poll"}
        </button>
      </form>
    </div>
  );
}
