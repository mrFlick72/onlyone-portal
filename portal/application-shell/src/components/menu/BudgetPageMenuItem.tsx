import * as PropTypes from "prop-types";
import React from "react";
import MenuItem from "./MenuItem";
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import { MenuItemProps } from "./Menu";

const BudgetPageMenuItem: React.FC<MenuItemProps> = ({ text }) => {
    return <MenuItem
        icon={<ShoppingCartIcon />}
        link="/budget/index"
        text={text} />
}

export default BudgetPageMenuItem