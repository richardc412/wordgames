import type { SquareState } from "./constants";

interface WordleSquareProps {
  character: string;
  state: SquareState;
}

const WordleSquare = ({ character, state }: WordleSquareProps) => {
  const getSquareClasses = () => {
    const baseClasses =
      "w-12 h-12 border-2 flex items-center justify-center text-xl font-bold uppercase transition-all duration-200";

    switch (state) {
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

  return <div className={getSquareClasses()}>{character}</div>;
};

export default WordleSquare;
