import React from 'react';
import FileSearch from './components/FileSearch';

const App: React.FC = () => (
    <div className="min-h-screen flex flex-col justify-start items-center bg-gray-100">
        <h1 className="text-center font-bold ">File Search</h1>
        <FileSearch />
    </div>
);

export default App;
