import React from "react";

export interface Summary {
    fileTypes: Record<string, number>;
    modifiedYears: Record<string, number>;
}

export interface Widget {
    id: string;
    name: string;
    isActive: (query: string, summary: Summary) => boolean;
    Component: React.FC;
}

class WidgetFactoryClass {
    private widgets: Widget[] = [];

    register(widget: Widget) {
        this.widgets.push(widget);
    }

    getActiveWidgets(query: string, summary: Summary): Widget[] {
        return this.widgets.filter(widget => widget.isActive(query, summary));
    }
}

export const WidgetFactory = new WidgetFactoryClass();
