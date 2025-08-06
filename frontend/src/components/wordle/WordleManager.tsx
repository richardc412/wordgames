import { useState } from "react";
import WordleGrid from "./WordleGrid";
import type { WordleSquareState } from "./types";

const generateEmptyGrid = (): WordleSquareState[][] => {
  return Array(6)
    .fill(null)
    .map(() =>
      Array(5)
        .fill(null)
        .map(() => ({
          character: "",
          state: "white",
        }))
    );
};

const WordleManager = () => {
  const [grid, setGrid] = useState<WordleSquareState[][]>(generateEmptyGrid());

  return (
    <>
      <WordleGrid grid={grid} />
    </>
  );
};

export default WordleManager;
