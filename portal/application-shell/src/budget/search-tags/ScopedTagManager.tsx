import React, { FC, useCallback, useEffect, useState } from 'react'
import { Container } from "@mui/material";
import SearchTagsTable from "./SearchTagsTable";
import SearchTagsForm from "./SearchTagsForm";
import { getSearchTagRegistry, saveSearchTag, TagScope, UNKNOWN_SENTINEL_KEY } from "./domain/SearchTagRepository";
import Separator from "../../components/form/Separator";
import { SearchTagFormMessageBundle, SearchTagsTableMessageBundle } from '../../messages/MessageBundles';

type ScopedTagManagerProps = {
    scope: TagScope;
    formMessages: SearchTagFormMessageBundle;
    tableMessages: SearchTagsTableMessageBundle;
}

/**
 * Manages the tag catalog for a single {@link TagScope}: lists the scope's
 * tags and creates new ones. Rendered once per tab (Expense / Revenue), so the
 * two tabs are guaranteed to behave identically — same component, two scopes.
 */
const ScopedTagManager: FC<ScopedTagManagerProps> = ({ scope, formMessages, tableMessages }) => {
    const [searchTagsRegistry, setSearchTagsRegistry] = useState<SearchTag[]>([])
    const [searchTag, setSearchTag] = useState<SearchTag>({ key: "", value: "" })

    // Load the scope's catalog, hiding the UNKNOWN technical sentinel — this is
    // a catalog editor, not a tag selector, so the non-persisted catch-all has
    // no place here.
    const loadRegistry = useCallback(() =>
        getSearchTagRegistry(scope).then(registry =>
            setSearchTagsRegistry(registry.filter(tag => tag.key !== UNKNOWN_SENTINEL_KEY))
        ), [scope])

    useEffect(() => { loadRegistry() }, [loadRegistry])

    const formHandler = {
        valueHandler: (value: React.ChangeEvent<HTMLInputElement>) => {
            setSearchTag((prevState) => ({
                ...prevState,
                value: value.target.value
            }))
        },
        submitHandler: (searchTagKey: string, searchTagValue: string) => {
            saveSearchTag({ key: searchTagKey, value: searchTagValue }, scope)
                .then(() => loadRegistry())
        }
    }

    const tableHandler = (searchTagKey: string, searchTagValue: string) => {
        setSearchTag({ key: searchTagKey, value: searchTagValue })
    }

    return <Container>
        <SearchTagsForm
            searchTag={{ key: searchTag.key, value: searchTag.value }}
            handler={formHandler}
            messages={formMessages} />

        <Separator />

        <SearchTagsTable
            searchTagsRegistry={searchTagsRegistry}
            handler={tableHandler}
            messages={tableMessages} />
    </Container>
}

export default ScopedTagManager
