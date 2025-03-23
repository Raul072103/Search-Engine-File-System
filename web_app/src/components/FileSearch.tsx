import React, { useState, useEffect } from "react";

const FileSearch: React.FC = () => {
    const [searchParams, setSearchParams] = useState({
        word_list: "",
        file_name: "",
        extensions: "",
    });
    const [results, setResults] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchResults = async () => {
            setLoading(true);
            setError(null);

            try {
                const query = new URLSearchParams();
                if (searchParams.word_list) {
                    searchParams.word_list.split(/\s+/).forEach(word => query.append("word_list", word));
                }
                if (searchParams.file_name) query.append("file_name", searchParams.file_name);
                if (searchParams.extensions) query.append("extensions", searchParams.extensions);

                const response = await fetch(`http://localhost:8080/v1/search?${query.toString()}`, {
                    method: "GET",
                    headers: { "Content-Type": "application/json" },
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }

                const data = await response.json();
                console.log("API Response:", data); // Check the structure of the response
                setResults(data.data || []);
            } catch (error) {
                setError("Error fetching search results.");
                console.error("Error fetching search results:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchResults();
    }, [searchParams]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setSearchParams((prev) => ({ ...prev, [name]: value }));
    };

    return (
        <div className="min-h-screen p-4 flex flex-col items-center bg-gray-100">
            <div className="w-full max-w-lg p-4 bg-white rounded-2xl shadow-lg space-y-4 mb-4">
                <input
                    type="text"
                    name="word_list"
                    placeholder="Word List"
                    value={searchParams.word_list}
                    onChange={handleChange}
                    className="w-full p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <input
                    type="text"
                    name="file_name"
                    placeholder="File Name"
                    value={searchParams.file_name}
                    onChange={handleChange}
                    className="w-full p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <input
                    type="text"
                    name="extensions"
                    placeholder="Extensions (comma-separated)"
                    value={searchParams.extensions}
                    onChange={handleChange}
                    className="w-full p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
            </div>

            {/* Results Section */}
            <div className="w-full max-w-lg mt-4 overflow-y-auto max-h-96">
                {loading ? (
                    <p className="text-gray-500 text-center">Loading...</p>
                ) : error ? (
                    <p className="text-red-500 text-center">{error}</p>
                ) : results.length > 0 ? (
                    <ul className="bg-white rounded-2xl shadow-lg p-4 divide-y divide-gray-200">
                        {results.map((result, index) => (
                            <li key={index} className="p-2">{result}</li>
                        ))}
                    </ul>
                ) : (
                    <p className="text-gray-500 text-center">No results found</p>
                )}
            </div>
        </div>
    );
};

export default FileSearch;
