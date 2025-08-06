import { useState, useEffect, useCallback } from "react";
import WordleGrid from "./WordleGrid";
import type { WordleSquareState } from "./types";
import {
  generateEmptyGrid,
  isWordValid,
  isWordInDictionary,
  withinBounds,
  generateHints,
} from "./helpers";

const ANSWER = "hello";

const WordleManager = () => {
  const [grid, setGrid] = useState<WordleSquareState[][]>(generateEmptyGrid());
  const [currentRow, setCurrentRow] = useState<number>(0);
  const [currentColumn, setCurrentColumn] = useState<number>(0);

  const getCurrentWord = useCallback((): string => {
    return grid[currentRow]
      .map((square) => square.character)
      .join("")
      .toLowerCase();
  }, [grid, currentRow]);

  const handleKeyPress = useCallback(
    (key: string): void => {
      if (!withinBounds(currentRow, currentColumn)) return;
      setGrid(
        grid.map((row, rowIndex) =>
          rowIndex === currentRow
            ? row.map((square, colIndex) =>
                colIndex === currentColumn
                  ? { ...square, character: key }
                  : square
              )
            : row
        )
      );
      setCurrentColumn(currentColumn + 1);
    },
    [currentRow, currentColumn, grid]
  );

  const handleEnterPress = useCallback((): void => {
    if (!withinBounds(currentRow, 0)) return;
    let currentWord = getCurrentWord();
    if (isWordValid(currentWord) && isWordInDictionary(currentWord)) {
      setGrid(
        grid.map((row, rowIndex) =>
          rowIndex === currentRow ? generateHints(currentWord, ANSWER) : row
        )
      );
      setCurrentRow(currentRow + 1);
      setCurrentColumn(0);
    }
  }, [getCurrentWord]);

  const handleBackspacePress = useCallback((): void => {
    if (!withinBounds(currentRow, currentColumn - 1)) return;
    setGrid(
      grid.map((row, rowIndex) =>
        rowIndex === currentRow
          ? row.map((square, colIndex) =>
              colIndex === currentColumn - 1
                ? { ...square, character: "" }
                : square
            )
          : row
      )
    );
    setCurrentColumn(currentColumn - 1);
  }, [currentRow, currentColumn, grid]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Handle Enter key
      if (event.key === "Enter") {
        handleEnterPress();
        return;
      }

      // Handle Backspace key
      if (event.key === "Backspace") {
        handleBackspacePress();
        return;
      }

      // Handle Latin alphabet characters (a-z, A-Z)
      if (/^[a-zA-Z]$/.test(event.key)) {
        handleKeyPress(event.key);
        return;
      }
    };

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [handleKeyPress, handleEnterPress, handleBackspacePress]);

  return (
    <>
      <WordleGrid grid={grid} />
    </>
  );
};

export default WordleManager;
