import React from 'react';
import { Button } from "../components/ui/button";

const Home: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-6 py-8 pt-20">
      <Button className="cursor-pointer">Create Room</Button>
    </div>
  );
};

export default Home; 