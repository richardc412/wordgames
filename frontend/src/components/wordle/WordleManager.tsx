import { useState, useEffect, useCallback } from "react";
import WordleGrid from "./WordleGrid";
import type { WordleSquareState } from "./types";

const GRID_WIDTH = 5;
const GRID_HEIGHT = 6;

const generateEmptyGrid = (): WordleSquareState[][] => {
  return Array(GRID_HEIGHT)
    .fill(null)
    .map(() =>
      Array(GRID_WIDTH)
        .fill(null)
        .map(() => ({
          character: "",
          state: "white",
        }))
    );
};

const isWordValid = (word: string) => {
  return word.length === GRID_WIDTH;
};

const isWordInDictionary = (word: string) => {
  // TODO: Implement this
  return word.length === GRID_WIDTH;
};

const withinBounds = (row: number, column: number) => {
  return row >= 0 && row < GRID_HEIGHT && column >= 0 && column < GRID_WIDTH;
};

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

  const replaceCharAtIndex = (
    word: string,
    index: number,
    replacement: string
  ): string => {
    return word.slice(0, index) + replacement + word.slice(index + 1);
  };

  const generateHints = useCallback(
    (currentWord: string, answer: string): WordleSquareState[] => {
      const hints: WordleSquareState[] = Array(GRID_WIDTH)
        .fill(null)
        .map(() => ({
          character: "",
          state: "white",
        }));

      let remainingAnswer = answer;
      let remainingWord = currentWord;

      // First pass: find exact matches (green)
      for (let i = 0; i < GRID_WIDTH; i++) {
        if (remainingWord[i] === remainingAnswer[i]) {
          hints[i] = { character: remainingWord[i], state: "green" };
          remainingAnswer = replaceCharAtIndex(remainingAnswer, i, " ");
          remainingWord = replaceCharAtIndex(remainingWord, i, " ");
        }
      }

      // Second pass: find misplaced letters (yellow)
      for (let i = 0; i < GRID_WIDTH; i++) {
        if (
          remainingWord[i] !== " " &&
          remainingAnswer.includes(remainingWord[i])
        ) {
          hints[i] = { character: remainingWord[i], state: "yellow" };
          const answerIndex = remainingAnswer.indexOf(remainingWord[i]);
          remainingAnswer = replaceCharAtIndex(
            remainingAnswer,
            answerIndex,
            " "
          );
          remainingWord = replaceCharAtIndex(remainingWord, i, " ");
        }
      }

      // Third pass: mark remaining as incorrect (grey)
      for (let i = 0; i < GRID_WIDTH; i++) {
        if (hints[i].state === "white") {
          hints[i] = { character: currentWord[i], state: "grey" };
        }
      }

      return hints;
    },
    []
  );

  const handleKeyPress = useCallback(
    (key: string): void => {
      if (withinBounds(currentRow, currentColumn)) {
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
      }
    },
    [currentRow, currentColumn, grid]
  );

  const handleEnterPress = useCallback((): void => {
    if (currentRow >= GRID_HEIGHT) {
      return;
    }
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

  const handleBackspacePress = useCallback(() => {
    if (withinBounds(currentRow, currentColumn - 1)) {
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
    }
  }, [currentRow, currentColumn, grid]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const key = event.key.toLowerCase();

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
