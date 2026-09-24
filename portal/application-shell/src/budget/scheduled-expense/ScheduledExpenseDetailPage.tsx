import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { Alert, Box, Button, CircularProgress, Container, Divider, Paper, ThemeProvider, Typography, Snackbar } from "@mui/material";
import { ArrowBack, Save } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import { FormDateFormatPattern } from "../../components/form/FormDatePicker";
import { SelectOption } from "../../components/form/FormSelect";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { MessageBundle } from "../../messages/MessageRepository";
import { getSearchTagRegistry } from "../search-tags/domain/SearchTagRepository";
import selectUiAdapterFor from "../search-tags/SearchTagsUIAdapter";
import ScheduledExpenseForm from "./ScheduledExpenseForm";
import { createScheduledExpense, getScheduledExpense, updateScheduledExpense } from "./domain/ScheduledExpenseRepository";

type ScheduledExpenseDetailPageProps = {
    messageRegistry: MessageBundle;
}

// Create-mode when the URL carries no ?id=, edit-mode when it does (#52) —
// one form serves both, per the design session. Status is never loaded or
// submitted here: it stays a list-row (pause/resume) action, per ADR 0005.
const ScheduledExpenseDetailPage: React.FC<ScheduledExpenseDetailPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const id = new URLSearchParams(window.location.search).get("id") ?? ""

    const [description, setDescription] = useState("")
    const [amount, setAmount] = useState("0.00")
    const [note, setNote] = useState("")
    const [searchTags, setSearchTags] = useState<SelectOption[]>([])
    const [searchTagRegistry, setSearchTagRegistry] = useState<SearchTag[]>([])
    const [day, setDay] = useState("")
    const [month, setMonth] = useState("")
    const [hasEndDate, setHasEndDate] = useState(false)
    const [endDate, setEndDate] = useState(moment().format(FormDateFormatPattern))

    const [saving, setSaving] = useState(false)
    const [feedback, setFeedback] = useState<{ severity: 'success' | 'error'; message: string } | null>(null)

    const detailMessages = configMap.scheduledExpenseDetail(messageRegistry)

    useEffect(() => {
        getSearchTagRegistry("expense").then(setSearchTagRegistry)
    }, [])

    useEffect(() => {
        if (!id) {
            return
        }
        getScheduledExpense(id).then(scheduledExpense => {
            if (!scheduledExpense) {
                setFeedback({ severity: 'error', message: detailMessages.feedback.loadError })
                return
            }
            setDescription(scheduledExpense.description)
            setAmount(scheduledExpense.amount)
            setNote(scheduledExpense.notes)
            setSearchTags(scheduledExpense.tags.map(tag => ({ value: tag.tagKey, label: tag.tagValue })))
            setDay(String(scheduledExpense.day))
            setMonth(scheduledExpense.month !== undefined ? String(scheduledExpense.month) : "")
            setHasEndDate(scheduledExpense.endDate !== undefined)
            // endDate travels as DD/MM/YYYY — the same FormDateFormatPattern the
            // form state uses (budget-api's date.DateFor), so no conversion.
            if (scheduledExpense.endDate) {
                setEndDate(scheduledExpense.endDate)
            }
        }).catch(() => {
            setFeedback({ severity: 'error', message: detailMessages.feedback.loadError })
        })
        // id is stable for the page's lifetime (a full navigation is required
        // to change it) — fetch once, not on every messageRegistry update.
    }, [id])

    const save = useCallback(() => {
        setSaving(true)
        const payload = {
            description,
            amount,
            notes: note,
            tags: searchTags.map(tag => ({ tagKey: tag.value, tagValue: tag.label })),
            day: Number(day),
            month: month === "" ? undefined : Number(month),
            endDate: hasEndDate ? endDate : undefined,
        }
        const action = id ? updateScheduledExpense(id, payload) : createScheduledExpense(payload)
        action.then(response => {
            if (response.status === 201 || response.status === 204) {
                setFeedback({ severity: 'success', message: detailMessages.feedback.success })
            } else {
                setFeedback({ severity: 'error', message: detailMessages.feedback.error })
            }
        }).catch(() => {
            setFeedback({ severity: 'error', message: detailMessages.feedback.error })
        }).finally(() => {
            setSaving(false)
        })
    }, [id, description, amount, note, searchTags, day, month, hasEndDate, endDate, detailMessages.feedback])

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={detailMessages.menuMessages} navBarItems={[]}>
                <MenuItem
                    icon={<ArrowBack />}
                    text={detailMessages.menuMessages.backToList}
                    link="/budget/scheduled-expense/index" />
            </Menu>
            <Container>
                <Typography variant="h5" sx={{ my: 2 }}>
                    {id ? detailMessages.headingEdit : detailMessages.heading}
                </Typography>

                <ScheduledExpenseForm
                    data={{
                        description,
                        amount,
                        note,
                        searchTags,
                        day,
                        month,
                        hasEndDate,
                        endDate,
                    }}
                    handlers={{
                        description: (event) => setDescription(event.target.value),
                        amount: (event) => setAmount(event.target.value),
                        note: (event) => setNote(event.target.value),
                        searchTag: (selected) => setSearchTags(selected ?? []),
                        day: (event) => setDay(event.target.value),
                        month: (event) => setMonth(event.target.value),
                        toggleEndDate: (event) => setHasEndDate(event.target.checked),
                        endDate: (value) => setEndDate(value.format(FormDateFormatPattern)),
                    }}
                    searchTagRegistry={selectUiAdapterFor(searchTagRegistry)}
                    messages={detailMessages.form} />

                <Divider sx={{ my: 2 }} />

                <Box sx={{ display: "flex", gap: 2, mb: 2 }}>
                    <Button
                        variant="contained"
                        color="success"
                        onClick={save}
                        disabled={saving}
                        startIcon={saving ? <CircularProgress size={16} color="inherit" /> : <Save />}>
                        {detailMessages.saveButtonLabel}
                    </Button>
                    <Button
                        variant="outlined"
                        href="/budget/scheduled-expense/index"
                        startIcon={<ArrowBack />}>
                        {detailMessages.backButtonLabel}
                    </Button>
                </Box>
            </Container>

            <Snackbar
                open={feedback !== null}
                autoHideDuration={6000}
                onClose={() => setFeedback(null)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                {feedback
                    ? <Alert severity={feedback.severity} onClose={() => setFeedback(null)} sx={{ width: '100%' }}>{feedback.message}</Alert>
                    : undefined}
            </Snackbar>
        </Paper>
    </ThemeProvider>
}

export default ScheduledExpenseDetailPage;
