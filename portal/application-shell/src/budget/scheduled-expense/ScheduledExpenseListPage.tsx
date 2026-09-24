import React, { useCallback, useEffect, useState } from "react"
import { Alert, Container, Paper, Snackbar, ThemeProvider } from "@mui/material";
import { EventRepeat } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { MessageBundle } from "../../messages/MessageRepository";
import ScheduledExpense from "./domain/ScheduledExpense";
import { changeScheduledExpenseStatus, deleteScheduledExpense, getAllScheduledExpenses } from "./domain/ScheduledExpenseRepository";
import ScheduledExpenseListContent from "./ScheduledExpenseListContent";
import DeleteScheduledExpenseConfirmationPopUp from "./DeleteScheduledExpenseConfirmationPopUp";

type ScheduledExpenseListPageProps = {
    messageRegistry: MessageBundle;
}

// Row actions: open (#51/#52), pause/resume toggle (#54) and delete (#53) —
// see the parent issue and ADR 0005.
const ScheduledExpenseListPage: React.FC<ScheduledExpenseListPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const listMessages = configMap.scheduledExpense(messageRegistry)
    const [scheduledExpenses, setScheduledExpenses] = useState<ScheduledExpense[]>([])

    const [deletable, setDeletable] = useState<ScheduledExpense | null>(null)
    const [openDeletePopUp, setOpenDeletePopUp] = useState(false)
    const [errorMessage, setErrorMessage] = useState<string | null>(null)

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
                setErrorMessage(listMessages.feedback.deleteError)
            }
        }).catch(() => {
            setErrorMessage(listMessages.feedback.deleteError)
        })
    }, [deletable, refresh, listMessages.feedback])

    // Sends the opposite of the row's current status. 404 means the row is
    // gone (deleted from another tab) — refresh drops it, no error needed.
    const toggleStatus = useCallback((scheduledExpense: ScheduledExpense) => {
        if (!scheduledExpense.id) {
            return
        }
        const target = scheduledExpense.status === "PAUSED" ? "ACTIVE" : "PAUSED"
        changeScheduledExpenseStatus(scheduledExpense.id, target).then(response => {
            if (response.status === 204 || response.status === 404) {
                refresh()
            } else {
                setErrorMessage(listMessages.feedback.statusError)
            }
        }).catch(() => {
            setErrorMessage(listMessages.feedback.statusError)
        })
    }, [refresh, listMessages.feedback])

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
                    toggleStatus={toggleStatus}
                    messages={listMessages.content} />
            </Container>
            <Snackbar
                open={errorMessage !== null}
                autoHideDuration={6000}
                onClose={() => setErrorMessage(null)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                <Alert severity="error" onClose={() => setErrorMessage(null)} sx={{ width: '100%' }}>
                    {errorMessage}
                </Alert>
            </Snackbar>
        </Paper>
    </ThemeProvider>
}

export default ScheduledExpenseListPage;
