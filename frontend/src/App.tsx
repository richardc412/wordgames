import "./App.css";
import Header from "./components/Header";
import ServerButton from "./components/ServerButton";

function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header title="Worduel" />
      <div className="max-w-7xl mx-auto px-6 py-8">
        <ServerButton />
      </div>
    </div>
  );
}

export default App;
