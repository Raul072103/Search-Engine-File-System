import React from 'react';
import FileSearch from './components/FileSearch';

const App: React.FC = () => (
    <div>
        <div className="user-input">
            <h1 className="font-bold">File Search</h1>
            <FileSearch/>
        </div>
    </div>
);

export default App;
