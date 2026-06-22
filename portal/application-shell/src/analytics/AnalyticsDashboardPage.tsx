import React, { useEffect, useState } from "react"
import { Alert, Box, Button, CircularProgress, Container, Paper, Snackbar, ThemeProvider, Typography } from "@mui/material"
import themeProvider from "../theme/ThemeProvider"
import Menu from "../components/menu/Menu"
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap"
import { MessageBundle } from "../messages/MessageRepository"
import { getMonthRegistry } from "../time/MonthRepository"
import { getSearchTagRegistry } from "../budget/search-tags/domain/SearchTagRepository"
import FormSelect, { SelectOption } from "../components/form/FormSelect"
import FormInputTextField from "../components/form/FormInputTextField"
import type { Month } from "../time/months"
import { findTotalByTag, findTotalByYear, reindexBudgetExpense, TagTotal, YearTotal } from "./domain/AnalyticRepository"
import AnalyticBarChart from "./charts/AnalyticBarChart"

const MAX_YEAR_RANGE_SIZE = 20

type AnalyticsDashboardPageProps = {
    messageRegistry: MessageBundle
}

function isValidYear(year: string): boolean {
    return /^\d{4}$/.test(year)
}

const AnalyticsDashboardPage: React.FC<AnalyticsDashboardPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const messages = configMap.analytics(messageRegistry)
    const currentYear = new Date().getFullYear()

    const [monthRegistry, setMonthRegistry] = useState<Month[]>([])
    const [searchTagRegistry, setSearchTagRegistry] = useState<SearchTag[]>([])

    const [year, setYear] = useState(String(currentYear))
    const [month, setMonth] = useState("")
    const [selectedTags, setSelectedTags] = useState<SelectOption[]>([])
    const [tagTotals, setTagTotals] = useState<TagTotal[]>([])
    const [tagTotalsLoading, setTagTotalsLoading] = useState(true)
    const [tagTotalsError, setTagTotalsError] = useState(false)
    const [tagTotalsRetryToken, setTagTotalsRetryToken] = useState(0)

    const [fromYear, setFromYear] = useState(String(currentYear - 4))
    const [toYear, setToYear] = useState(String(currentYear))
    const [selectedTag, setSelectedTag] = useState<SelectOption | null>(null)
    const [yearTotals, setYearTotals] = useState<YearTotal[]>([])
    const [yearTotalsLoading, setYearTotalsLoading] = useState(true)
    const [yearTotalsError, setYearTotalsError] = useState(false)
    const [yearTotalsRetryToken, setYearTotalsRetryToken] = useState(0)

    const [reindexFromYear, setReindexFromYear] = useState(String(currentYear))
    const [reindexToYear, setReindexToYear] = useState(String(currentYear))
    const [reindexing, setReindexing] = useState(false)
    const [reindexFeedback, setReindexFeedback] = useState<{ severity: 'success' | 'error'; message: string } | null>(null)

    useEffect(() => {
        getMonthRegistry().then((data) => setMonthRegistry(data))
        getSearchTagRegistry("expense").then((data) => setSearchTagRegistry(data))
    }, [])

    useEffect(() => {
        if (!isValidYear(year)) {
            return
        }
        let cancelled = false
        setTagTotalsLoading(true)
        setTagTotalsError(false)
        findTotalByTag({
            year: Number(year),
            month: month ? Number(month) : undefined,
            tags: selectedTags.map((tag) => tag.value),
        }).then((data) => {
            if (!cancelled) {
                setTagTotals(data)
                setTagTotalsLoading(false)
            }
        }).catch(() => {
            if (!cancelled) {
                setTagTotalsError(true)
                setTagTotalsLoading(false)
            }
        })
        return () => {
            cancelled = true
        }
    }, [year, month, selectedTags, tagTotalsRetryToken])

    const invalidYearRange = isValidYear(fromYear) && isValidYear(toYear) &&
        (Number(fromYear) > Number(toYear) || Number(toYear) - Number(fromYear) + 1 > MAX_YEAR_RANGE_SIZE)

    useEffect(() => {
        if (!isValidYear(fromYear) || !isValidYear(toYear) || invalidYearRange) {
            return
        }
        let cancelled = false
        setYearTotalsLoading(true)
        setYearTotalsError(false)
        findTotalByYear({
            fromYear: Number(fromYear),
            toYear: Number(toYear),
            tag: selectedTag?.value || undefined,
        }).then((data) => {
            if (!cancelled) {
                setYearTotals(data)
                setYearTotalsLoading(false)
            }
        }).catch(() => {
            if (!cancelled) {
                setYearTotalsError(true)
                setYearTotalsLoading(false)
            }
        })
        return () => {
            cancelled = true
        }
    }, [fromYear, toYear, invalidYearRange, selectedTag, yearTotalsRetryToken])

    const monthOptions = [
        { value: "", label: messages.filters.allMonths },
        ...monthRegistry.map((item) => ({ value: String(item.monthValue), label: item.monthLabel })),
    ]
    const selectedMonthOption = monthOptions.find((option) => option.value === month) ?? monthOptions[0]

    const tagOptions = searchTagRegistry.map((tag) => ({ value: tag.key, label: tag.value }))
    const singleTagOptions = [{ value: "", label: messages.filters.allTags }, ...tagOptions]

    const chartStateMessages = {
        empty: messages.states.empty,
        error: messages.states.error,
        retry: messages.states.retry,
    }

    const reindexInvalidYearRange = isValidYear(reindexFromYear) && isValidYear(reindexToYear) &&
        (Number(reindexFromYear) > Number(reindexToYear) || Number(reindexToYear) - Number(reindexFromYear) + 1 > MAX_YEAR_RANGE_SIZE)
    const reindexDisabled = reindexing || !isValidYear(reindexFromYear) || !isValidYear(reindexToYear) || reindexInvalidYearRange

    const handleReindex = () => {
        setReindexing(true)
        reindexBudgetExpense({
            fromYear: Number(reindexFromYear),
            toYear: Number(reindexToYear),
        }).then((result) => {
            setReindexFeedback({ severity: 'success', message: `${messages.reindex.success} (${result.imported})` })
            // refresh both charts against the freshly rebuilt projection
            setTagTotalsRetryToken((token) => token + 1)
            setYearTotalsRetryToken((token) => token + 1)
        }).catch(() => {
            setReindexFeedback({ severity: 'error', message: messages.reindex.error })
        }).finally(() => {
            setReindexing(false)
        })
    }

    return (
        <ThemeProvider theme={themeProvider}>
            <Paper variant="outlined" sx={{ minHeight: '100vh' }}>
                <Menu messages={messages.menuMessages} navBarItems={[]} />

                <Container sx={{ pt: 6, pb: 4 }}>
                    <Paper sx={{ p: 3, mb: 4 }}>
                        <Typography variant="h6" gutterBottom>{messages.charts.totalByTagTitle}</Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, alignItems: 'flex-start', mb: 2 }}>
                            <Box sx={{ width: { xs: '100%', md: 150 } }}>
                                <FormInputTextField
                                    id="analyticsYear"
                                    label={messages.filters.year}
                                    type="number"
                                    value={year}
                                    handler={(event: React.ChangeEvent<HTMLInputElement>) => setYear(event.target.value)} />
                            </Box>
                            <Box sx={{ width: { xs: '100%', md: 220 } }}>
                                <FormSelect
                                    id="analyticsMonth"
                                    label={messages.filters.month}
                                    multi={false}
                                    options={monthOptions}
                                    value={selectedMonthOption}
                                    onChangeHandler={(option: SelectOption | null) => setMonth(option?.value ?? "")} />
                            </Box>
                            <Box sx={{ flex: 1, minWidth: { xs: '100%', md: 300 } }}>
                                <FormSelect
                                    id="analyticsTags"
                                    label={messages.filters.tags}
                                    multi={true}
                                    options={tagOptions}
                                    value={selectedTags}
                                    onChangeHandler={(options: readonly SelectOption[] | null) =>
                                        setSelectedTags(options ? [...options] : [])} />
                            </Box>
                        </Box>
                        <AnalyticBarChart
                            labels={tagTotals.map((item) => item.tag)}
                            values={tagTotals.map((item) => Number(item.total))}
                            seriesLabel={messages.charts.totalByTagTitle}
                            loading={tagTotalsLoading}
                            error={tagTotalsError}
                            messages={chartStateMessages}
                            retryHandler={() => setTagTotalsRetryToken((token) => token + 1)} />
                    </Paper>

                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6" gutterBottom>{messages.charts.totalByYearTitle}</Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, alignItems: 'flex-start', mb: 2 }}>
                            <Box sx={{ width: { xs: '100%', md: 150 } }}>
                                <FormInputTextField
                                    id="analyticsFromYear"
                                    label={messages.filters.fromYear}
                                    type="number"
                                    value={fromYear}
                                    handler={(event: React.ChangeEvent<HTMLInputElement>) => setFromYear(event.target.value)} />
                            </Box>
                            <Box sx={{ width: { xs: '100%', md: 150 } }}>
                                <FormInputTextField
                                    id="analyticsToYear"
                                    label={messages.filters.toYear}
                                    type="number"
                                    value={toYear}
                                    handler={(event: React.ChangeEvent<HTMLInputElement>) => setToYear(event.target.value)} />
                            </Box>
                            <Box sx={{ flex: 1, minWidth: { xs: '100%', md: 300 } }}>
                                <FormSelect
                                    id="analyticsTag"
                                    label={messages.filters.tag}
                                    multi={false}
                                    options={singleTagOptions}
                                    value={selectedTag ?? singleTagOptions[0]}
                                    onChangeHandler={(option: SelectOption | null) =>
                                        setSelectedTag(option && option.value ? option : null)} />
                            </Box>
                        </Box>
                        {invalidYearRange && <Alert severity="warning" sx={{ mb: 2 }}>{messages.filters.invalidYearRange}</Alert>}
                        <AnalyticBarChart
                            labels={yearTotals.map((item) => String(item.year))}
                            values={yearTotals.map((item) => Number(item.total))}
                            seriesLabel={messages.charts.totalByYearTitle}
                            loading={yearTotalsLoading}
                            error={yearTotalsError}
                            messages={chartStateMessages}
                            retryHandler={() => setYearTotalsRetryToken((token) => token + 1)} />
                    </Paper>

                    <Paper sx={{ p: 3, mt: 4 }}>
                        <Typography variant="h6" gutterBottom>{messages.reindex.title}</Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>{messages.reindex.description}</Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, alignItems: 'center' }}>
                            <Box sx={{ width: { xs: '100%', md: 150 } }}>
                                <FormInputTextField
                                    id="reindexFromYear"
                                    label={messages.filters.fromYear}
                                    type="number"
                                    value={reindexFromYear}
                                    handler={(event: React.ChangeEvent<HTMLInputElement>) => setReindexFromYear(event.target.value)} />
                            </Box>
                            <Box sx={{ width: { xs: '100%', md: 150 } }}>
                                <FormInputTextField
                                    id="reindexToYear"
                                    label={messages.filters.toYear}
                                    type="number"
                                    value={reindexToYear}
                                    handler={(event: React.ChangeEvent<HTMLInputElement>) => setReindexToYear(event.target.value)} />
                            </Box>
                            <Button
                                variant="contained"
                                onClick={handleReindex}
                                disabled={reindexDisabled}
                                startIcon={reindexing ? <CircularProgress size={16} color="inherit" /> : undefined}>
                                {reindexing ? messages.reindex.running : messages.reindex.button}
                            </Button>
                        </Box>
                        {reindexInvalidYearRange && <Alert severity="warning" sx={{ mt: 2 }}>{messages.filters.invalidYearRange}</Alert>}
                    </Paper>
                </Container>

                <Snackbar
                    open={reindexFeedback !== null}
                    autoHideDuration={6000}
                    onClose={() => setReindexFeedback(null)}
                    anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
                    {reindexFeedback
                        ? <Alert severity={reindexFeedback.severity} onClose={() => setReindexFeedback(null)} sx={{ width: '100%' }}>{reindexFeedback.message}</Alert>
                        : undefined}
                </Snackbar>
            </Paper>
        </ThemeProvider>
    )
}

export default AnalyticsDashboardPage
