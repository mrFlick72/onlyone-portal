import React from "react";
import MenuItem from "./MenuItem";
import { FormatListBulleted } from "@mui/icons-material";
import { MenuItemProps } from "./Menu";

const PlanPageMenuItem: React.FC<MenuItemProps> = ({ text }) => (
    <MenuItem
        icon={<FormatListBulleted />}
        link="/plan/index"
        text={text} />
);

export default PlanPageMenuItem
