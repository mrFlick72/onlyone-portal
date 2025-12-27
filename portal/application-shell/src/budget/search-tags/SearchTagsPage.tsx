import React, { FC, useEffect, useState } from 'react'
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import SearchTagsTable from "./SearchTagsTable";
import SearchTagsForm from "./SearchTagsForm";
import { getSearchTagRegistry, saveSearchTag } from "./domain/SearchTagRepository";
import { Container, Paper, ThemeProvider } from "@mui/material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import AccountPageMenuItem from "../../components/menu/AccountPageMenuItem";
import Separator from "../../components/form/Separator";
import { MessageBundle } from '../../messages/MessageRepository';

interface SearchTagsPageProps {
    messageRegistry: MessageBundle;
}
const SearchTagsPage: FC<SearchTagsPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()

    const [searchTagsRegistry, setSearchTagsRegistry] = useState<SearchTag[]>([])
    const [searchTagKey, setSearchTagKey] = useState("")
    const [searchTagValue, setSearchTagValue] = useState("")

    useEffect(() => {
        getSearchTagRegistry().then(registry => {
            setSearchTagsRegistry(registry)
        })
    }, [])

    const formHandler = {
        valueHandler: (value: React.ChangeEvent<HTMLInputElement>) => {
            setSearchTagValue(value.target.value);
        },
        submitHandler: (searchTagKey: string, searchTagValue: string) => {
            saveSearchTag({ key: searchTagKey, value: searchTagValue })
                .then(ignore => getSearchTagRegistry())
                .then(registry => {
                    setSearchTagsRegistry(registry)
                })
        }
    }



    const tableHandler = (searchTagKey: string, searchTagValue: string) => {
        setSearchTagKey(searchTagKey)
        setSearchTagValue(searchTagValue)
    }

    let theme = themeProvider

    return <ThemeProvider theme={theme}>
        <Paper variant="outlined">
            <Menu messages={configMap.searchTags(messageRegistry).menuMessages} navBarItems={[]}>
                <AccountPageMenuItem text={configMap.searchTags(messageRegistry).menuMessages.userProfileLabel} />
            </Menu>

            <Container>
                <SearchTagsForm searchTag={{ key: searchTagKey, value: searchTagValue }} handler={formHandler} />

                <Separator />

                <SearchTagsTable searchTagsRegistry={searchTagsRegistry} handler={tableHandler} />
            </Container>
        </Paper>
    </ThemeProvider>
}
export default SearchTagsPage