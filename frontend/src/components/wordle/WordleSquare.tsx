import type { WordleSquareState } from "./types";

interface WordleSquareProps {
  WordleSquareState: WordleSquareState;
}

const WordleSquare = ({ WordleSquareState }: WordleSquareProps) => {
  const getSquareClasses = () => {
    const baseClasses =
      "w-12 h-12 border-2 flex items-center justify-center text-xl font-bold uppercase transition-all duration-200";

    switch (WordleSquareState.state) {
      case "white":
        return `${baseClasses} border-gray-300 bg-white`;
      case "disabled":
        return `${baseClasses} border-gray-500 bg-gray-100`;
      case "grey":
        return `${baseClasses} border-gray-500 bg-gray-500 text-white`;
      case "yellow":
        return `${baseClasses} border-yellow-500 bg-yellow-500 text-white`;
      case "green":
        return `${baseClasses} border-green-500 bg-green-500 text-white`;
      default:
        return `${baseClasses} border-gray-300 bg-white`;
    }
  };

  return (
    <div className={getSquareClasses()}>{WordleSquareState.character}</div>
  );
};

export default WordleSquare;
