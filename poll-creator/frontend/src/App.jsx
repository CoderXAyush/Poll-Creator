import { useState } from "react";
import { BrowserRouter, Routes, Route, useNavigate, useParams } from "react-router-dom";
import CreatePoll from "./components/CreatePoll";
import PollList from "./components/PollList";
import PollView from "./components/PollView";

function AppContent() {
  const [tab, setTab] = useState("create");
  const navigate = useNavigate();

  return (
    <div className="app">
      <header className="header">
        <h1>Poll Creator</h1>
        <p>Create polls, share them, and collect votes instantly</p>
      </header>

      <Routes>
        <Route
          path="/"
          element={
            <>
              <nav className="nav">
                <button
                  className={tab === "create" ? "active" : ""}
                  onClick={() => setTab("create")}
                >
                  Create
                </button>
                <button
                  className={tab === "browse" ? "active" : ""}
                  onClick={() => setTab("browse")}
                >
                  Browse
                </button>
              </nav>
              {tab === "create" ? (
                <CreatePoll
                  onCreated={(poll) => {
                    navigate(`/poll/${poll.id}`);
                  }}
                />
              ) : (
                <PollList
                  onSelect={(poll) => navigate(`/poll/${poll.id}`)}
                />
              )}
            </>
          }
        />
        <Route path="/poll/:id" element={<PollView />} />
      </Routes>
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
}
