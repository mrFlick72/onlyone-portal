import React, { useCallback, useEffect, useState } from "react"
import moment, { type Moment } from "moment"
import { Container, Paper, Tab, Tabs, ThemeProvider } from "@mui/material"
import { LocalGroceryStore, Search } from "@mui/icons-material"
import {
    getMonthSearchCriteria,
    getSearchTagsSearchCriteria,
    getYearSearchCriteria,
    setMonthSearchCriteria,
    setSearchTagsSearchCriteria,
    setYearSearchCriteria,
} from "../SearchCriteriaProvider"
import { deleteBudgetExpense, findBudgetExpense, saveBudgetExpense } from "./domain/BudgetExpenseRepository"
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap"
import { getMonthRegistry } from "../../time/MonthRepository"
import { getSearchTagRegistry } from "../search-tags/domain/SearchTagRepository"
import themeProvider from "../../theme/ThemeProvider"
import SaveBudgetExpensePopUp from "./popup/SaveBudgetExpensePopUp"
import SpentBudgetContent from "./budget/SpentBudgetContent"
import DeleteBudgetExpenseConfirmationPopUp from "./popup/DeleteBudgetExpenseConfirmationPopUp"
import TotalBySearchTags from "./budget/TotalBySearchTags"
import SpentBudgetTotalBanner from "./budget/SpentBudgetTotalBanner"
import SearchBudgetExpensePopUp from "./popup/SearchBudgetExpensePopUp"
import Menu from "../../components/menu/Menu"
import OpenPopUpMenuItem from "../../components/menu/OpenPopUpMenuItem"
import SearchTagsPageMenuItem from "../../components/menu/SearchTagsPageMenuItem"
import TabPanel from "../../components/layout/TabPanel"
import { FormDateFormatPattern } from "../../components/form/FormDatePicker"
import type { SelectOption } from "../../components/form/FormSelect"
import type { Month } from "../../time/months"
import type { BudgetExpense, SavedBudgetExpense, SpentBudget } from "./domain/BudgetExpense"

interface BudgetExpensePageProps {
    messageRegistry: any
}

const emptySpentBudget: SpentBudget = {
    dailyBudgetExpenseRepresentationList: [],
    totalDetailList: [],
    total: "0.00",
}

const emptySearchTag: SelectOption = {
    value: "",
    label: "",
}

const BudgetExpensePage: React.FC<BudgetExpensePageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()

    const [tabPanel, setTabPanel] = useState(0)
    const [id, setId] = useState("")
    const [date, setDate] = useState<Moment>(moment())
    const [amount, setAmount] = useState("0.00")
    const [note, setNote] = useState("")
    const [searchTag, setSearchTag] = useState<SelectOption>(emptySearchTag)
    const [spentBudget, setSpentBudget] = useState<SpentBudget>(emptySpentBudget)
    const [deletableItem, setDeletableItem] = useState<BudgetExpense | null>(null)
    const [searchTagRegistry, setSearchTagRegistry] = useState<SearchTag[]>([])
    const [monthRegistry, setMonthRegistry] = useState<Month[]>([])

    const [openSaveBudgetExpensePopUp, setOpenSaveBudgetExpensePopUp] = useState(false)
    const [openDeleteBudgetExpensePopUp, setOpenDeleteBudgetExpensePopUp] = useState(false)
    const [openSearchBudgetExpensePopUp, setOpenSearchBudgetExpensePopUp] = useState(false)

    const [selectedMonth, setSelectedMonth] = useState<string>(getMonthSearchCriteria())
    const [selectedYear, setSelectedYear] = useState<string>(getYearSearchCriteria())
    const [selectedSearchTags, setSelectedSearchTags] = useState<SelectOption[]>(
        getSearchTagsSearchCriteria()
            .filter(Boolean)
            .map((value) => ({ value, label: value }))
    )

    const getSpentBudget = useCallback(() => {
        findBudgetExpense({
            month: getMonthSearchCriteria(),
            year: getYearSearchCriteria(),
        }).then((data) => {
            setSpentBudget(data)
        })
    }, [])

    const loadCommonData = useCallback(() => {
        getMonthRegistry().then((data) => setMonthRegistry(data))
        getSearchTagRegistry().then((data) => {
            setSearchTagRegistry(data)
            setSelectedSearchTags((currentSelectedSearchTags) =>
                currentSelectedSearchTags.map((selectedTag) => {
                    const foundSearchTag = data.find((item) => item.key === selectedTag.value)
                    return foundSearchTag
                        ? { value: foundSearchTag.key, label: foundSearchTag.value }
                        : selectedTag
                })
            )
        })
    }, [])

    const makeNewBudgetExpensePopUpOpen = useCallback(() => {
        setId("")
        setDate(moment())
        setAmount("0.00")
        setNote("")
        setSearchTag(emptySearchTag)
        setOpenSaveBudgetExpensePopUp(true)
    }, [])

    const makeUpdateBudgetExpensePopUpOpen = useCallback((expense: SavedBudgetExpense) => {
        setId(expense.id ?? "")
        setDate(moment(expense.date, FormDateFormatPattern))
        setAmount(expense.amount.toString())
        setNote(expense.note)
        setSearchTag(expense.searchTag)
        setOpenSaveBudgetExpensePopUp(true)
    }, [])

    const saveBudgetExpensePopUpCloseHandler = useCallback(() => {
        setOpenSaveBudgetExpensePopUp(false)
    }, [])

    const makeDeleteBudgetExpensePopUpOpen = useCallback((dailyBudgetExpense: BudgetExpense) => {
        setOpenDeleteBudgetExpensePopUp(true)
        setDeletableItem(dailyBudgetExpense)
    }, [])

    const deleteBudgetExpensePopUpCloseHandler = useCallback(() => {
        setDeletableItem(null)
        setOpenDeleteBudgetExpensePopUp(false)
    }, [])

    const makeSearchBudgetExpensePopUpOpen = useCallback(() => {
        setOpenSearchBudgetExpensePopUp(true)
    }, [])

    const searchBudgetExpensePopUpCloseHandler = useCallback(() => {
        setOpenSearchBudgetExpensePopUp(false)
    }, [])

    const deleteItem = useCallback(() => {
        if (!deletableItem?.id) {
            return
        }

        deleteBudgetExpense(deletableItem.id).then((response) => {
            if (response.status === 204) {
                getSpentBudget()
                setOpenDeleteBudgetExpensePopUp(false)
            }
        })
    }, [deletableItem, getSpentBudget])

    const savePopupEventHandlers = {
        date: (value: Moment | null) => {
            setDate(value ?? moment())
        },
        amount: (event: React.ChangeEvent<HTMLInputElement>) => {
            setAmount(event.target.value)
        },
        searchTag: (selectedSearchTag: SelectOption | null) => {
            setSearchTag(selectedSearchTag ?? emptySearchTag)
        },
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => {
            setNote(event.target.value)
        },
    }

    const searchPopupEventHandlers = {
        searchTag: (selectedOptions: readonly SelectOption[] | null) => {
            setSelectedSearchTags(selectedOptions ? [...selectedOptions] : [])
        },
        month: (selectedOption: SelectOption | null) => {
            setSelectedMonth(selectedOption?.value ?? getMonthSearchCriteria())
        },
        year: (event: React.ChangeEvent<HTMLInputElement>) => {
            setSelectedYear(event.target.value)
        },
    }

    const saveExpense = useCallback(() => {
        saveBudgetExpense({
            id,
            date: date.format(FormDateFormatPattern),
            amount: Number(amount),
            note,
            tagKey: searchTag.value,
            tagValue: searchTag.label,
        }).then((response) => {
            if (response.status === 201 || response.status === 204) {
                getSpentBudget()
                setOpenSaveBudgetExpensePopUp(false)
            }
        })
    }, [amount, date, getSpentBudget, id, note, searchTag])

    useEffect(() => {
        loadCommonData()
        getSpentBudget()
    }, [getSpentBudget, loadCommonData])

    const handleTabPanelChange = (_event: React.SyntheticEvent, newValue: number) => {
        setTabPanel(newValue)
    }

    const navBarItems = [<SpentBudgetTotalBanner key="spent-budget-total" total={spentBudget.total} />]

    return (
        <ThemeProvider theme={themeProvider}>
            <Paper variant="outlined">
                <SearchBudgetExpensePopUp
                    month={selectedMonth}
                    year={selectedYear}
                    searchTags={selectedSearchTags}
                    monthRegistry={monthRegistry}
                    searchTagRegistry={searchTagRegistry}
                    handlers={searchPopupEventHandlers}
                    modal={configMap.budgetExpense(messageRegistry).searchFilterModal}
                    handleClose={searchBudgetExpensePopUpCloseHandler}
                    open={openSearchBudgetExpensePopUp}
                    saveCallback={() => {
                        setMonthSearchCriteria(selectedMonth)
                        setYearSearchCriteria(selectedYear)
                        setSearchTagsSearchCriteria(selectedSearchTags.map((tag) => tag.value).join(","))
                        setOpenSearchBudgetExpensePopUp(false)
                        getSpentBudget()
                    }}
                />

                <SaveBudgetExpensePopUp
                    open={openSaveBudgetExpensePopUp}
                    handleClose={saveBudgetExpensePopUpCloseHandler}
                    spentBudgetHandlers={savePopupEventHandlers}
                    budgetExpense={{
                        id,
                        date: date.format(FormDateFormatPattern),
                        amount: Number(amount),
                        note,
                        searchTag,
                    }}
                    saveCallback={saveExpense}
                    searchTagRegistry={searchTagRegistry}
                    modal={configMap.budgetExpense(messageRegistry).newBudgetExpenseModal}
                />

                <DeleteBudgetExpenseConfirmationPopUp
                    deleteBudgetExpenseAction={deleteItem}
                    open={openDeleteBudgetExpensePopUp}
                    handleClose={deleteBudgetExpensePopUpCloseHandler}
                    modal={configMap.budgetExpense(messageRegistry).deleteModal}
                />

                <Menu messages={configMap.searchTags(messageRegistry).menuMessages} navBarItems={navBarItems}>
                    <OpenPopUpMenuItem
                        icon={<LocalGroceryStore />}
                        openPopupHandler={makeNewBudgetExpensePopUpOpen}
                        text={configMap.budgetExpense(messageRegistry).menuMessages.insertBudgetModal}
                    />

                    <OpenPopUpMenuItem icon={<Search />} openPopupHandler={makeSearchBudgetExpensePopUpOpen} text="Search" />

                    <SearchTagsPageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.searchTags} />
                </Menu>

                <Container>
                    <Tabs
                        value={tabPanel}
                        onChange={handleTabPanelChange}
                        textColor="secondary"
                        indicatorColor="secondary"
                        aria-label="secondary tabs example"
                    >
                        <Tab value={0} label="Daily View" />
                        <Tab value={1} label="By Tags View" />
                    </Tabs>
                    <TabPanel value={tabPanel} index={0}>
                        <SpentBudgetContent
                            spentBudget={spentBudget}
                            openUpdateBudgetExpensePopUp={makeUpdateBudgetExpensePopUpOpen}
                            openDeleteBudgetExpensePopUp={makeDeleteBudgetExpensePopUpOpen}
                        />
                    </TabPanel>
                    <TabPanel value={tabPanel} index={1}>
                        <TotalBySearchTags totals={spentBudget.totalDetailList || []} />
                    </TabPanel>
                </Container>
            </Paper>
        </ThemeProvider>
    )
}

export default BudgetExpensePage
