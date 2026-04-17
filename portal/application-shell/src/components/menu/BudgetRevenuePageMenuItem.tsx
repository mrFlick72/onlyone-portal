import * as PropTypes from "prop-types";
import MenuItem from "./MenuItem";
import { Money } from "@mui/icons-material";
import React from "react";
import { MenuItemProps } from "./Menu";

const BudgetRevenuePageMenuItem: React.FC<MenuItemProps> = ({ text }) => {
    return <MenuItem
        icon={<Money />}
        link="/budget/revenue/index"
        text={text} />
}

export default BudgetRevenuePageMenuItem
