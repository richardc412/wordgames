import "./App.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Header from "./components/Header";
import Home from "./pages/Home";
import Match from "./pages/Match";

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Header title="Worduel" />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/match" element={<Match />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
