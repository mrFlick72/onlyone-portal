import React, { FC, useCallback, useEffect, useState } from 'react'
import { Container } from "@mui/material";
import SearchTagsTable from "./SearchTagsTable";
import SearchTagsForm from "./SearchTagsForm";
import UpdateSearchTagPopUp from "./UpdateSearchTagPopUp";
import DeleteSearchTagConfirmationPopUp from "./DeleteSearchTagConfirmationPopUp";
import { deleteSearchTag, getSearchTagRegistry, saveSearchTag, TagScope, UNKNOWN_SENTINEL_KEY, updateSearchTagValue } from "./domain/SearchTagRepository";
import Separator from "../../components/form/Separator";
import { DeleteModalMessageBundle, SaveModalMessageBundle, SearchTagFormMessageBundle, SearchTagsTableMessageBundle } from '../../messages/MessageBundles';

type ScopedTagManagerProps = {
    scope: TagScope;
    formMessages: SearchTagFormMessageBundle;
    tableMessages: SearchTagsTableMessageBundle;
    editModal: SaveModalMessageBundle;
    deleteModal: DeleteModalMessageBundle;
}

const emptySearchTag: SearchTag = { key: "", value: "" }

/**
 * Manages the tag catalog for a single {@link TagScope}: lists the scope's
 * tags, creates new ones, and edits/deletes existing ones. Rendered once per
 * tab (Expense / Revenue), so the two tabs are guaranteed to behave
 * identically — same component, two scopes.
 */
const ScopedTagManager: FC<ScopedTagManagerProps> = ({ scope, formMessages, tableMessages, editModal, deleteModal }) => {
    const [searchTagsRegistry, setSearchTagsRegistry] = useState<SearchTag[]>([])

    // The create form always starts empty and resets after every save — it
    // never reflects a tag selected for editing (that goes through the
    // dedicated edit popup below).
    const [newSearchTag, setNewSearchTag] = useState<SearchTag>(emptySearchTag)

    const [editingSearchTag, setEditingSearchTag] = useState<SearchTag>(emptySearchTag)
    const [openEditPopUp, setOpenEditPopUp] = useState(false)

    const [deletingSearchTag, setDeletingSearchTag] = useState<SearchTag>(emptySearchTag)
    const [openDeletePopUp, setOpenDeletePopUp] = useState(false)

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
            setNewSearchTag((prevState) => ({
                ...prevState,
                value: value.target.value
            }))
        },
        submitHandler: (_searchTagKey: string, searchTagValue: string) => {
            saveSearchTag({ key: "", value: searchTagValue }, scope)
                .then(() => {
                    setNewSearchTag(emptySearchTag)
                    loadRegistry()
                })
        }
    }

    const editValueHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEditingSearchTag((prevState) => ({
            ...prevState,
            value: event.target.value
        }))
    }

    const openEdit = (searchTag: SearchTag) => {
        setEditingSearchTag(searchTag)
        setOpenEditPopUp(true)
    }
    const closeEdit = () => setOpenEditPopUp(false)
    const saveEdit = () => {
        updateSearchTagValue(editingSearchTag.key, editingSearchTag.value, scope)
            .then(() => {
                setOpenEditPopUp(false)
                loadRegistry()
            })
    }

    const openDelete = (searchTag: SearchTag) => {
        setDeletingSearchTag(searchTag)
        setOpenDeletePopUp(true)
    }
    const closeDelete = () => setOpenDeletePopUp(false)
    const confirmDelete = () => {
        deleteSearchTag(deletingSearchTag.key, scope)
            .then(() => {
                setOpenDeletePopUp(false)
                loadRegistry()
            })
    }

    return <Container>
        <SearchTagsForm
            searchTag={newSearchTag}
            handler={formHandler}
            messages={formMessages} />

        <Separator />

        <SearchTagsTable
            searchTagsRegistry={searchTagsRegistry}
            handler={{ onEdit: openEdit, onDelete: openDelete }}
            messages={tableMessages} />

        <UpdateSearchTagPopUp
            open={openEditPopUp}
            handleClose={closeEdit}
            value={editingSearchTag.value}
            valueHandler={editValueHandler}
            saveCallback={saveEdit}
            modal={editModal}
            formMessages={formMessages} />

        <DeleteSearchTagConfirmationPopUp
            open={openDeletePopUp}
            handleClose={closeDelete}
            saveCallback={confirmDelete}
            modal={deleteModal} />
    </Container>
}

export default ScopedTagManager
