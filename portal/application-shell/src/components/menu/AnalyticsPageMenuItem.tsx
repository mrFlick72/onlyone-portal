import React from "react";
import MenuItem from "./MenuItem";
import { BarChart } from "@mui/icons-material";
import { MenuItemProps } from "./Menu";

const AnalyticsPageMenuItem: React.FC<MenuItemProps> = ({ text }) => (
    <MenuItem
        icon={<BarChart />}
        link="/analytics/index"
        text={text} />
);

export default AnalyticsPageMenuItem
