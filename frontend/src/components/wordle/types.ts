export type SquareState = "white" | "disabled" | "green" | "yellow" | "grey";

export type WordleSquareState = {
  character: string;
  state: SquareState;
};
