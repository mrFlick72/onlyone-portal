import React, { FC, useEffect, useState } from 'react'
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import SearchTagsTable from "./SearchTagsTable";
import SearchTagsForm from "./SearchTagsForm";
import { getSearchTagRegistry, saveSearchTag } from "./domain/SearchTagRepository";
import { Container, Paper, ThemeProvider } from "@mui/material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import Separator from "../../components/form/Separator";
import { MessageBundle } from '../../messages/MessageRepository';

interface SearchTagsPageProps {
    messageRegistry: MessageBundle;
}
const SearchTagsPage: FC<SearchTagsPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()

    const [searchTagsRegistry, setSearchTagsRegistry] = useState<SearchTag[]>([])

    const [searchTag, setSearchTag] = useState<SearchTag>({ key: "", value: "" })


    useEffect(() => {
        getSearchTagRegistry().then(registry => {
            setSearchTagsRegistry(registry)
        })
    }, [])

    const formHandler = {
        valueHandler: (value: React.ChangeEvent<HTMLInputElement>) => {
            setSearchTag((prevState) => ({
                ...prevState,
                value: value.target.value
            }))

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
        setSearchTag({ key: searchTagKey, value: searchTagValue })
    }

    let theme = themeProvider

    return <ThemeProvider theme={theme}>
        <Paper variant="outlined">
            <Menu messages={configMap.searchTags(messageRegistry).menuMessages} navBarItems={[]} />

            <Container>
                <SearchTagsForm searchTag={{ key: searchTag.key, value: searchTag.value }} handler={formHandler} />

                <Separator />

                <SearchTagsTable searchTagsRegistry={searchTagsRegistry} handler={tableHandler} />
            </Container>
        </Paper>
    </ThemeProvider>
}
export default SearchTagsPage