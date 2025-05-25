import { WidgetFactory, Summary } from "./WidgetFactory";
import { AnalyzeLogsWidget, ViewGalleryWidget, ExploreSourceCodeWidget} from "./Widgets";

WidgetFactory.register({
    id: "analyzeLogs",
    name: "Analyze Logs",
    isActive: (_: string, summary: Summary) => {
        const totalFiles = Object.values(summary.fileTypes).reduce((a, b) => a + b, 0);
        if (totalFiles === 0) return false;
        const logCount = summary.fileTypes[".log"] || 0;
        return logCount / totalFiles > 0.5;
    },
    Component: AnalyzeLogsWidget,
});

WidgetFactory.register({
    id: "viewGallery",
    name: "View as Gallery",
    isActive: (query: string, summary: Summary) => {
        const imageExists = [".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp"];
        const totalFiles = Object.values(summary.fileTypes).reduce((a, b) => a + b, 0);

        if (totalFiles === 0)
            return false;

        const imageCount = imageExists.reduce((count, ext) => count + (summary.fileTypes[ext] || 0), 0);
        const queryHasImage = query.toLowerCase().includes("image");
        return imageCount / totalFiles > 0.3 || queryHasImage;
    },
    Component: ViewGalleryWidget,
});


WidgetFactory.register({
    id: "exploreCode",
    name: "Explore Source Code",
    isActive: (_: string, summary: Summary) => {
        const codeExtensions = [".js", ".ts", ".py", ".go", ".java", ".cpp"];
        const totalFiles = Object.values(summary.fileTypes).reduce((a, b) => a + b, 0);
        if (totalFiles === 0) return false;

        const codeCount = codeExtensions.reduce((count, ext) => count + (summary.fileTypes[ext] || 0), 0);
        return codeCount / totalFiles > 0.4;
    },
    Component: ExploreSourceCodeWidget,
});
