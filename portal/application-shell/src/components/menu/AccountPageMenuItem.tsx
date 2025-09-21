import React from "react";

import MenuItem from "./MenuItem";
import { Person } from "@mui/icons-material";
import { MenuItemProps } from "./Menu";

const AccountPageMenuItem: React.FC<MenuItemProps> = ({ text }) => {
    return <MenuItem
        icon={<Person />}
        link="/account/index"
        text={text} />
}

export default AccountPageMenuItem