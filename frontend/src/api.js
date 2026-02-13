const API_BASE = "/api";

// Generate or retrieve a session ID for duplicate vote prevention
function getSessionId() {
  let id = sessionStorage.getItem("poll-session-id");
  if (!id) {
    id = crypto.randomUUID?.() || Math.random().toString(36).slice(2);
    sessionStorage.setItem("poll-session-id", id);
  }
  return id;
}

function headers() {
  return {
    "Content-Type": "application/json",
    "x-session-id": getSessionId(),
  };
}

export async function createPoll(question, options) {
  const res = await fetch(`${API_BASE}/polls`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({ question, options }),
  });
  if (!res.ok) throw new Error((await res.json()).error);
  return res.json();
}

export async function getPolls() {
  const res = await fetch(`${API_BASE}/polls`, { headers: headers() });
  if (!res.ok) throw new Error("Failed to fetch polls");
  return res.json();
}

export async function getPoll(id) {
  const res = await fetch(`${API_BASE}/polls/${id}`, { headers: headers() });
  if (!res.ok) throw new Error("Poll not found");
  return res.json();
}

export async function submitVote(pollId, optionId) {
  const res = await fetch(`${API_BASE}/polls/${pollId}/vote`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({ optionId }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error);
  return data;
}

export async function getPollResults(id) {
  const res = await fetch(`${API_BASE}/polls/${id}/results`, {
    headers: headers(),
  });
  if (!res.ok) throw new Error("Failed to fetch results");
  return res.json();
}

export async function closePoll(id) {
  const res = await fetch(`${API_BASE}/polls/${id}/close`, {
    method: "PATCH",
    headers: headers(),
  });
  if (!res.ok) throw new Error("Failed to close poll");
  return res.json();
}
