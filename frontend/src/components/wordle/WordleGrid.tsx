import React from "react";
import WordleSquare from "./WordleSquare";
import type { WordleSquareState } from "./types";

interface WordleGridProps {
  grid: WordleSquareState[][];
}

const WordleGrid: React.FC<WordleGridProps> = ({ grid }) => {
  return (
    <div className="flex flex-col gap-2">
      {grid.map((row, rowIndex) => (
        <div key={rowIndex} className="flex gap-2 justify-center">
          {row.map((square, colIndex) => (
            <WordleSquare
              key={`${rowIndex}-${colIndex}`}
              WordleSquareState={square}
            />
          ))}
        </div>
      ))}
    </div>
  );
};

export default WordleGrid;
