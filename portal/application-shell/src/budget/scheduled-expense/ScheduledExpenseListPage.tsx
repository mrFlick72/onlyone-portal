import React, { useCallback, useEffect, useState } from "react"
import { Alert, Container, Paper, Snackbar, ThemeProvider } from "@mui/material";
import { EventRepeat } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { MessageBundle } from "../../messages/MessageRepository";
import ScheduledExpense from "./domain/ScheduledExpense";
import { deleteScheduledExpense, getAllScheduledExpenses } from "./domain/ScheduledExpenseRepository";
import ScheduledExpenseListContent from "./ScheduledExpenseListContent";
import DeleteScheduledExpenseConfirmationPopUp from "./DeleteScheduledExpenseConfirmationPopUp";

type ScheduledExpenseListPageProps = {
    messageRegistry: MessageBundle;
}

// Row actions: open (#51/#52) and delete (#53). Pause/resume lands in #54 —
// see the parent issue and ADR 0005.
const ScheduledExpenseListPage: React.FC<ScheduledExpenseListPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const [scheduledExpenses, setScheduledExpenses] = useState<ScheduledExpense[]>([])

    const [deletable, setDeletable] = useState<ScheduledExpense | null>(null)
    const [openDeletePopUp, setOpenDeletePopUp] = useState(false)
    const [deleteError, setDeleteError] = useState(false)

    const refresh = useCallback(() => {
        getAllScheduledExpenses().then(setScheduledExpenses)
    }, [])

    useEffect(() => { refresh() }, [refresh])

    const openDetail = useCallback((scheduledExpense: ScheduledExpense) => {
        window.location.href = `/budget/scheduled-expense/detail?id=${encodeURIComponent(scheduledExpense.id ?? "")}`
    }, [])

    const openDelete = useCallback((scheduledExpense: ScheduledExpense) => {
        setDeletable(scheduledExpense)
        setOpenDeletePopUp(true)
    }, [])

    const closeDelete = useCallback(() => setOpenDeletePopUp(false), [])

    const confirmDelete = useCallback(() => {
        if (!deletable?.id) {
            return
        }
        deleteScheduledExpense(deletable.id).then(response => {
            // 404 means it's already gone (e.g. deleted from another tab) —
            // same end state as a successful delete, so don't strand the popup.
            if (response.status === 204 || response.status === 404) {
                setOpenDeletePopUp(false)
                refresh()
            } else {
                setDeleteError(true)
            }
        }).catch(() => {
            setDeleteError(true)
        })
    }, [deletable, refresh])

    const listMessages = configMap.scheduledExpense(messageRegistry)

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={listMessages.menuMessages} navBarItems={[]}>
                <MenuItem
                    icon={<EventRepeat />}
                    text={listMessages.menuMessages.newScheduledExpense}
                    link="/budget/scheduled-expense/detail" />
            </Menu>
            <Container>
                <DeleteScheduledExpenseConfirmationPopUp
                    open={openDeletePopUp}
                    handleClose={closeDelete}
                    saveCallback={confirmDelete}
                    modal={listMessages.deleteModal} />

                <ScheduledExpenseListContent
                    scheduledExpenses={scheduledExpenses}
                    openDetail={openDetail}
                    openDelete={openDelete}
                    messages={listMessages.content} />
            </Container>
            <Snackbar
                open={deleteError}
                autoHideDuration={6000}
                onClose={() => setDeleteError(false)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                <Alert severity="error" onClose={() => setDeleteError(false)} sx={{ width: '100%' }}>
                    {listMessages.feedback.deleteError}
                </Alert>
            </Snackbar>
        </Paper>
    </ThemeProvider>
}

export default ScheduledExpenseListPage;
