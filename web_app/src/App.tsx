import React from 'react';
import FileSearch from './components/FileSearch';

const App: React.FC = () => (
    <div className="containter">
        <div className="user-input">
            <h1 className="text-center font-bold clea">File Search</h1>
            <FileSearch/>
        </div>
        <div className="results-section">

        </div>
    </div>
);

export default App;
