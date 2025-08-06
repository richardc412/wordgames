import React from 'react';

interface HeaderProps {
  title?: string;
}

const Header: React.FC<HeaderProps> = ({ title = "Word Games" }) => {
  return (
    <header className="fixed top-0 left-0 right-0 bg-white shadow-sm border-b border-gray-200 px-6 py-4 z-50">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900">
          {title}
        </h1>
      </div>
    </header>
  );
};

export default Header; 