import React, { useState, useEffect } from "react";

const RESULTS_PER_PAGE = 3; // Number of results per page

const FileSearch: React.FC = () => {
    const [searchParams, setSearchParams] = useState({
        word_list: "",
        file_name: "",
        extensions: "",
    });

    const [results, setResults] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(1);

    useEffect(() => {
        // Don't fetch if all inputs are empty
        if (!searchParams.word_list && !searchParams.file_name && !searchParams.extensions) {
            setResults([]);  // Clear results if inputs are empty
            return;
        }

        const fetchResults = async () => {
            setLoading(true);
            setError(null);

            try {
                const query = new URLSearchParams();
                if (searchParams.word_list) {
                    searchParams.word_list.split(/\s+/).forEach(word => query.append("word_list", word));
                }
                if (searchParams.file_name) query.append("file_name", searchParams.file_name);
                if (searchParams.extensions) {
                    searchParams.extensions.split(/\s+/).forEach(word => query.append("extensions", word));
                }

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
                setCurrentPage(1); // Reset to the first page on new search
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

    // Pagination Logic
    const totalPages = Math.ceil(results.length / RESULTS_PER_PAGE);
    const startIndex = (currentPage - 1) * RESULTS_PER_PAGE;
    const displayedResults = results.slice(startIndex, startIndex + RESULTS_PER_PAGE);

    return (
        <div className="min-h-screen p-4 flex flex-col items-center bg-gray-100">
            {/* Input Fields */}
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

            {/* Display Results */}
            <div className="w-full max-w-lg bg-white rounded-2xl shadow-lg p-4">
                {loading ? (
                    <p className="text-center text-blue-500">Loading...</p>
                ) : error ? (
                    <p className="text-center text-red-500">{error}</p>
                ) : results.length > 0 ? (
                    <>
                        {/* Pagination Controls */}
                        <div className="flex justify-center mt-4 space-x-2">
                            <button
                                onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                                disabled={currentPage === 1}
                                className="px-3 py-1 border rounded disabled:opacity-50"
                            >
                                Previous
                            </button>
                            <span className="px-3 py-1">{currentPage} / {totalPages}</span>
                            <button
                                onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                                disabled={currentPage === totalPages}
                                className="px-3 py-1 border rounded disabled:opacity-50"
                            >
                                Next
                            </button>
                        </div>
                        <ul className="divide-y divide-gray-200 max-h-80 overflow-y-auto">
                            {displayedResults.map((result, index) => (
                                <li key={index} className="p-2 text-sm">
                                    <div><strong>Name:</strong> {result.Name}</div>
                                    <div><strong>Path:</strong> {result.Path}</div>
                                    <div><strong>Size:</strong> {result.Size} bytes</div>
                                    <div><strong>Extension:</strong> {result.Extension}</div>
                                    <div><strong>Updated At:</strong> {new Date(result.UpdatedAt).toLocaleString()}
                                    <div><strong>File Preview</strong> {result.Content.Text}</div>
                                    </div>
                                </li>
                            ))}
                        </ul>
                    </>
                ) : (
                    <p className="text-gray-500 text-center">No results found</p>
                )}
            </div>
        </div>
    );
};

export default FileSearch;
