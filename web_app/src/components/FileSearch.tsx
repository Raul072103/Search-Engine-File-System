import React, { useState, useEffect } from "react";

const RESULTS_PER_PAGE = 3; // Number of results per page

const FileSearch: React.FC = () => {
    const [searchParams, setSearchParams] = useState({
        word_list: "",
        file_name: "",
        extensions: "",
    });

    const [searchQuery, setSearchQuery] = useState("");
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const [corrections, setCorrections] = useState<string[]>([]);
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
                    searchParams.word_list.split(/\s+/)
                        .filter(word => word != "")
                        .forEach(word => query.append("word_list", word));
                }
                if (searchParams.file_name) query.append("file_name", searchParams.file_name);

                if (searchParams.extensions) {
                    searchParams.extensions.split(/\s+/)
                        .filter(word => word != "")
                        .forEach(word => query.append("extensions", word));
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

    useEffect(() => {
        const parseSearchQuery = () => {
            const params =  {
                word_list: "",
                file_name: "",
                extensions: "" };

            const parts = searchQuery.split(/\s+/);

            parts.forEach(part => {
                const [key, ...values] = part.split(':');
                const  currParams= values.join("").trim()

                if (currParams !== "") {
                    switch (key?.trim().toLowerCase()) {
                        case "path":
                            // TODO() add file path search
                            break;
                        case "content":
                            const words = currParams.split(',').join(" ")
                            params.word_list = `${params.word_list} ${words}`;
                            break;
                        case "extensions":
                            const extensions = currParams.split(',').join(" ")
                            params.extensions = `${params.extensions} ${extensions}`;
                            break;
                        case "name":
                            params.file_name = currParams;
                            break;
                        default:
                            break;
                    }
                }
            });
            setSearchParams(params);
        };

        if (searchQuery.trim() != "") {
            parseSearchQuery();
        } else {
            setSearchParams({ word_list: "", file_name: "", extensions: "" });
        }

    }, [searchQuery]);

    useEffect(() => {
        const controller = new AbortController();
        const fetchSuggestions = async () => {
            try {
                if (!searchQuery.trim()) {
                    setSuggestions([]);
                    return;
                }

                const res = await fetch(`http://localhost:8080/v1/query-suggestions?query=${encodeURIComponent(searchQuery)}`, {
                    signal: controller.signal
                });

                if (!res.ok) {
                    throw new Error("Failed to fetch suggestions");
                }

                const data = await res.json();
                setSuggestions(data || []);
            } catch (err) {
                console.error("Error fetching search results:", err);
            }
        };

        fetchSuggestions();
        return () => controller.abort();
    }, [searchQuery]);

    useEffect(() => {
        const controller = new AbortController();

        const fetchCorrections = async () => {
            try {
                const query = new URLSearchParams();
                if (searchParams.word_list) {
                    searchParams.word_list.split(/\s+/)
                        .filter(word => word != "")
                        .forEach(word => query.append("word_list", word));
                }
                if (searchParams.file_name) query.append("file_name", searchParams.file_name);

                const res = await fetch(`http://localhost:8080/v1/query-spell-corrector?${query.toString()}`, {
                    method: "GET",
                    headers: { "Content-Type": "application/json" },
                });

                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }

                const data = await res.json();
                const { file_name_suggestion, word_list_suggestions } = data;

                // Reconstruct the suggested query
                const suggestedWords = word_list_suggestions.join(" ").trim();
                const suggestedQueryParts = [];

                if (suggestedWords) {
                    suggestedQueryParts.push(`content:${word_list_suggestions.join(",")}`);
                }

                if (file_name_suggestion)
                    suggestedQueryParts.push(`name:${file_name_suggestion}`);

                const suggestedQuery = suggestedQueryParts.join(" ");

                // Show only if it differs from the original query
                if (suggestedQuery.trim() && suggestedQuery.trim() !== searchQuery.trim()) {
                    setCorrections([suggestedQuery]);
                } else {
                    setCorrections([]);
                }

            } catch (err) {
                console.error("Error fetching corrections:", err);
            }
        };

        fetchCorrections();
        return () => controller.abort();
    }, [searchParams]);


    // Pagination Logic
    const totalPages = Math.ceil(results.length / RESULTS_PER_PAGE);
    const startIndex = (currentPage - 1) * RESULTS_PER_PAGE;
    const displayedResults = results.slice(startIndex, startIndex + RESULTS_PER_PAGE);

    return (
        <div>
            {/* Parser Input */}
            <div>
                <input
                    type="text"
                    name="search_bar"
                    placeholder="Enter search query"
                    className="search-bar"
                    value={searchQuery}
                    onChange={(e) => {
                        setSearchQuery(e.target.value)
                        console.log("searchQuery updated:", e.target.value); // Add this line
                    }}/>
            </div>
             {/* Input Fields */}
            <div>
                <input
                    type="text"
                    name="word_list"
                    placeholder="Word List"
                    value={searchParams.word_list}
                />
                <input
                    type="text"
                    name="file_name"
                    placeholder="File Name"
                    value={searchParams.file_name}
                />
                <input
                    type="text"
                    name="extensions"
                    placeholder="Extensions (comma-separated)"
                    value={searchParams.extensions}
                />
            </div>

            {/* Corrections */}
            {corrections.length > 0 && (
                <div style={{ marginTop: "8px", marginBottom: "16px" }}>
                    <strong>Did you mean:</strong>
                    <ul>
                        {corrections.map((correction, idx) => (
                            <li
                                key={idx}
                                style={{ cursor: "pointer", color: "darkorange" }}
                                onClick={() => setSearchQuery(correction)}
                            >
                                {correction}
                            </li>
                        ))}
                    </ul>
                </div>
            )}


            {/* Suggestions */}
            <div style={{ marginTop: "8px", marginBottom: "16px" }}>
                {suggestions.length > 0 && (
                    <div>
                        <strong>Suggestions:</strong>
                        <ul>
                            {suggestions.map((suggestion, idx) => (
                                <li
                                    key={idx}
                                    style={{ cursor: "pointer", color: "blue" }}
                                    onClick={() => setSearchQuery(suggestion)}
                                >
                                    {suggestion}
                                </li>
                            ))}
                        </ul>
                    </div>
                )}
            </div>

            {/* Display Results */}
            <div>
                {loading ? (
                    <p>Loading...</p>
                ) : error ? (
                    <p>{error}</p>
                ) : results.length > 0 ? (
                    <>
                        {/* Pagination Controls */}
                        <div>
                            <button
                                onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                                disabled={currentPage === 1}
                            >
                                Previous
                            </button>
                            <span>{currentPage} / {totalPages}</span>
                            <button
                                onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
                                disabled={currentPage === totalPages}
                            >
                                Next
                            </button>
                        </div>
                        <ul>
                            {displayedResults.map((result, index) => (
                                <li key={index}>
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
                    <p>No results found</p>
                )}
            </div>
        </div>
    );
};

export default FileSearch;
