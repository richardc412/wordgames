import type { WordleSquareState } from "./types";

export const GRID_WIDTH = 5;
export const GRID_HEIGHT = 6;

export const generateEmptyGrid = (): WordleSquareState[][] => {
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

export const isWordValid = (word: string) => {
  return word.length === GRID_WIDTH;
};

export const isWordInDictionary = (word: string) => {
  // TODO: Implement this
  return word.length === GRID_WIDTH;
};

export const withinBounds = (row: number, column: number) => {
  return row >= 0 && row < GRID_HEIGHT && column >= 0 && column < GRID_WIDTH;
};

export const generateHints = (
  currentWord: string,
  answer: string
): WordleSquareState[] => {
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
      remainingAnswer = replaceCharAtIndex(remainingAnswer, answerIndex, " ");
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
};

export const replaceCharAtIndex = (
  word: string,
  index: number,
  replacement: string
): string => {
  return word.slice(0, index) + replacement + word.slice(index + 1);
};
