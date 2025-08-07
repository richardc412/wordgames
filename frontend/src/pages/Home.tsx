import React, { useState } from "react";
import { Button } from "../components/ui/button";
import { api } from "../lib/api";
import type { CreateMatchResponse } from "../types/api";

const Home: React.FC = () => {
  const [matchData, setMatchData] = useState<CreateMatchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreateRoom = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await api.createMatch();
      setMatchData(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create room");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-6 py-8 pt-20">
      <div className="space-y-6">
        <Button
          onClick={handleCreateRoom}
          disabled={loading}
          className="cursor-pointer"
        >
          {loading ? "Creating Room..." : "Create Room"}
        </Button>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}

        {matchData && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
            <h3 className="font-bold mb-2">Room Created Successfully!</h3>
            <div className="space-y-2 text-sm">
              <p>
                <strong>Match ID:</strong> {matchData.matchId}
              </p>
              <p>
                <strong>Player ID:</strong> {matchData.playerId}
              </p>
              <p>
                <strong>Token:</strong>{" "}
                <code className="bg-green-200 px-1 rounded">
                  {matchData.token.substring(0, 20)}...
                </code>
              </p>
              <p>
                <strong>Invite URL:</strong>{" "}
                <code className="bg-green-200 px-1 rounded">
                  {matchData.invite}
                </code>
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
