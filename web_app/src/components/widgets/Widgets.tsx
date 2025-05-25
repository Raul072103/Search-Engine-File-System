// src/components/widgets/Widgets.tsx
import React from "react";

const buttonStyle: React.CSSProperties = {
    backgroundColor: "#3498db",
    color: "white",
    padding: "12px 20px",
    margin: "10px",
    border: "none",
    borderRadius: "6px",
    cursor: "pointer",
    fontSize: "16px",
    fontWeight: "bold",
    transition: "background-color 0.3s ease",
};

const hoverStyle: React.CSSProperties = {
    backgroundColor: "#2980b9",
};

export const ViewGalleryWidget: React.FC = () => {
    const [hovered, setHovered] = React.useState(false);

    return (
        <button
            style={{ ...buttonStyle, ...(hovered ? hoverStyle : {}) }}
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
        >
            View as Gallery
        </button>
    );
};

export const AnalyzeLogsWidget: React.FC = () => {
    const [hovered, setHovered] = React.useState(false);

    return (
        <button
            style={{ ...buttonStyle, ...(hovered ? hoverStyle : {}) }}
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
        >
            Analyze Logs
        </button>
    );
};

export const ExploreSourceCodeWidget: React.FC = () => {
    const [hovered, setHovered] = React.useState(false);

    return (
        <button
            style={{ ...buttonStyle, ...(hovered ? hoverStyle : {}) }}
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
        >
            Explore Source Code
        </button>
    );
};
