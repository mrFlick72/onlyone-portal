import {createRoot} from 'react-dom/client';
import HomePage from "./HomePage";

if (document.getElementById('app')) {
    const container = document.getElementById('app');
    if (container) {
        const root = createRoot(container);
        root.render(<HomePage />);
    }
}