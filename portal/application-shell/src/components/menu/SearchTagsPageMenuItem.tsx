import React from "react";

import MenuItem from "./MenuItem";
import { LocalOffer } from "@mui/icons-material";
import { MenuItemProps } from "./Menu";
import { TagScope } from "../../budget/search-tags/domain/SearchTagRepository";

type SearchTagsPageMenuItemProps = MenuItemProps & {
    // Deep-links to the matching tab on the tag-management page. Omitted -> the
    // page defaults to the Expense tab.
    scope?: TagScope;
}

const SearchTagsPageMenuItem: React.FC<SearchTagsPageMenuItemProps> = ({ text, scope }) => {
    const link = scope ? `/budget/search-tags?scope=${scope}` : "/budget/search-tags"
    return <MenuItem
        icon={<LocalOffer />}
        link={link}
        text={text} />
}

export default SearchTagsPageMenuItem
