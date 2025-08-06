import React from "react";
import WordleManager from "../components/WordleManager";

const Match: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-6 py-8 pt-20">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Match Page</h1>
        <WordleManager />
      </div>
    </div>
  );
};

export default Match;
