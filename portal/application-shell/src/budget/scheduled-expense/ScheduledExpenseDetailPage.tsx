import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { Alert, Box, Button, CircularProgress, Container, Divider, Paper, Snackbar, ThemeProvider } from "@mui/material";
import { ArrowBack, Save } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../../components/form/FormDatePicker";
import { SelectOption } from "../../components/form/FormSelect";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { MessageBundle } from "../../messages/MessageRepository";
import { getSearchTagRegistry } from "../search-tags/domain/SearchTagRepository";
import selectUiAdapterFor from "../search-tags/SearchTagsUIAdapter";
import ScheduledExpenseForm from "./ScheduledExpenseForm";
import { createScheduledExpense } from "./domain/ScheduledExpenseRepository";

type ScheduledExpenseDetailPageProps = {
    messageRegistry: MessageBundle;
}

// #51 ships create-mode only. Edit-mode (reading ?id=, loading the existing
// values, submitting via PUT) lands in #52 once budget-api has an Update
// action and a single-item lookup — see the parent issue.
const ScheduledExpenseDetailPage: React.FC<ScheduledExpenseDetailPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()

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

    useEffect(() => {
        getSearchTagRegistry("expense").then(setSearchTagRegistry)
    }, [])

    const detailMessages = configMap.scheduledExpenseDetail(messageRegistry)

    // Stays on the page and confirms via toast, rather than navigating back to
    // the list — the list link is a separate, explicit action (menu and the
    // bottom "back" button), not an implicit side effect of saving.
    const save = useCallback(() => {
        setSaving(true)
        createScheduledExpense({
            description,
            amount,
            notes: note,
            tags: searchTags.map(tag => ({ tagKey: tag.value, tagValue: tag.label })),
            day: Number(day),
            month: month === "" ? undefined : Number(month),
            endDate: hasEndDate ? moment(endDate, FormDateFormatPattern).format(ApiDateFormatPattern) : undefined,
        }).then(response => {
            if (response.status === 201) {
                setFeedback({ severity: 'success', message: detailMessages.feedback.success })
            } else {
                setFeedback({ severity: 'error', message: detailMessages.feedback.error })
            }
        }).catch(() => {
            setFeedback({ severity: 'error', message: detailMessages.feedback.error })
        }).finally(() => {
            setSaving(false)
        })
    }, [description, amount, note, searchTags, day, month, hasEndDate, endDate, detailMessages.feedback])

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={detailMessages.menuMessages} navBarItems={[]}>
                <MenuItem
                    icon={<ArrowBack />}
                    text={detailMessages.menuMessages.backToList}
                    link="/budget/scheduled-expense/index" />
            </Menu>
            <Container>
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
