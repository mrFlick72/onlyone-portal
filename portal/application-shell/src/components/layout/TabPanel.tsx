import React, { ReactNode } from "react"

type TabPanelProps = {
    value: number
    index: number
    children: ReactNode
}

const TabPanel: React.FC<TabPanelProps> = ({ value, index, children }) => {
    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
        >
            {children}
        </div>
    );
}

export default TabPanel