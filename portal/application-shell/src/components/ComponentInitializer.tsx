import React, { ReactNode } from "react"
import { createRoot } from "react-dom/client";


export default (component: ReactNode) => {
    if (document.getElementById('app')) {
        const container = document.getElementById('app');
        const root = createRoot(container!!);
        root.render(component);
    }
}