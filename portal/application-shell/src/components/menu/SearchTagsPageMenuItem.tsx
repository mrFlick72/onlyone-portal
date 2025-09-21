import React from "react";

import MenuItem from "./MenuItem";
import { LocalOffer } from "@mui/icons-material";
import { MenuItemProps } from "./Menu";

const SearchTagsPageMenuItem: React.FC<MenuItemProps> = ({ text }) => {
    return <MenuItem
        icon={<LocalOffer />}
        link="/budget/index#/search-tags"
        text={text} />
}

export default SearchTagsPageMenuItem