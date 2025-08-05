import React, { useState } from 'react';

const ServerButton: React.FC = () => {
  const [response, setResponse] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const handleClick = async () => {
    setLoading(true);
    setError('');
    setResponse('');
    
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.text();
      setResponse(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-4">
      <button
        onClick={handleClick}
        disabled={loading}
        className="px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {loading ? 'Loading...' : 'Fetch from Server'}
      </button>
      
      {response && (
        <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
          <h3 className="font-medium text-green-800 mb-2">Server Response:</h3>
          <p className="text-green-700">{response}</p>
        </div>
      )}
      
      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
          <h3 className="font-medium text-red-800 mb-2">Error:</h3>
          <p className="text-red-700">{error}</p>
        </div>
      )}
    </div>
  );
};

export default ServerButton; 