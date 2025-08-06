import React from "react";
import WordleManager from "../components/wordle/WordleManager";

const Match: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-6 py-8 pt-20">
      <div className="text-center">
        <WordleManager />
      </div>
    </div>
  );
};

export default Match;
