import React, { useEffect, useState } from "react"
import { Alert, Box, Container, Paper, ThemeProvider, Typography } from "@mui/material"
import themeProvider from "../theme/ThemeProvider"
import Menu from "../components/menu/Menu"
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap"
import { MessageBundle } from "../messages/MessageRepository"
import { getMonthRegistry } from "../time/MonthRepository"
import { getSearchTagRegistry } from "../budget/search-tags/domain/SearchTagRepository"
import FormSelect, { SelectOption } from "../components/form/FormSelect"
import FormInputTextField from "../components/form/FormInputTextField"
import type { Month } from "../time/months"
import { findTotalByTag, findTotalByYear, TagTotal, YearTotal } from "./domain/AnalyticRepository"
import AnalyticBarChart from "./charts/AnalyticBarChart"

const MAX_YEAR_RANGE_SIZE = 20

interface AnalyticsDashboardPageProps {
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

    useEffect(() => {
        getMonthRegistry().then((data) => setMonthRegistry(data))
        getSearchTagRegistry().then((data) => setSearchTagRegistry(data))
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
                </Container>
            </Paper>
        </ThemeProvider>
    )
}

export default AnalyticsDashboardPage
