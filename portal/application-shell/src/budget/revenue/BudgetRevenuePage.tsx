import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import * as searchCriteria from "../SearchCriteriaProvider";
import { deleteBudgetRevenue, findBudgetRevenue, saveBudgetRevenue } from "./domain/BudgetRevenueRepository";
import themeProvider from "../../theme/ThemeProvider";
import { Container, Paper, ThemeProvider } from "@mui/material";
import { Money, Search } from "@mui/icons-material";
import BudgetRevenueContent from "./BudgetRevenueContent";
import RevenueSearchPopUp from "./RevenueSearchPopUp";
import DeleteBudgetRevenueConfirmationPopUp from "./DeleteBudgetRevenueConfirmationPopUp";
import SaveBudgetRevenuePopUp from "./SaveBudgetRevenuePopUp";
import Menu from "../../components/menu/Menu";
import OpenPopUpMenuItem from "../../components/menu/OpenPopUpMenuItem";
import { FormDateFormatPattern } from "../../components/form/FormDatePicker";
import BudgetRevenue from "./domain/BudgetRevenue";
interface BudgetRevenuePageProps {
    messageRegistry: any;
}

const BudgetRevenuePage: React.FC<BudgetRevenuePageProps> = ({ messageRegistry }) => {

    const [deletableItem, setDeletableItem] = useState<BudgetRevenue | null>(null)
    const [revenues, setRevenues] = useState<BudgetRevenue[]>([])
    const [total, setTotal] = useState("0.00")
    const [currentBudgetRevenueId, setCurrentBudgetRevenueId] = useState("")
    const [currentBudgetRevenueDate, setCurrentBudgetRevenueDate] = useState(moment())
    const [currentBudgetRevenueAmount, setCurrentBudgetRevenueAmount] = useState("0.00")
    const [currentBudgetRevenueNote, setCurrentBudgetRevenueNote] = useState("")

    const [openSaveBudgetRevenuePopUp, setOpenSaveBudgetRevenuePopUp] = useState(false)
    const makeSaveBudgetRevenuePopUpOpen = useCallback(() => {
        setCurrentBudgetRevenueId("")
        setCurrentBudgetRevenueDate(moment())
        setCurrentBudgetRevenueAmount("0.00")
        setCurrentBudgetRevenueNote("")
        setOpenSaveBudgetRevenuePopUp(true)
    }, [])
    const saveBudgetRevenuePopUpCloseHandler = useCallback(() => {
        setOpenSaveBudgetRevenuePopUp(false)
    }, [])

    const makeUpdateBudgetRevenuePopUpOpen = useCallback((revenue: BudgetRevenue) => {
        setCurrentBudgetRevenueId(revenue.id!)
        setCurrentBudgetRevenueDate(moment(revenue.date, "DD/MM/YYYY"))
        setCurrentBudgetRevenueAmount(revenue.amount!.toString())
        setCurrentBudgetRevenueNote(revenue.note)
        setOpenSaveBudgetRevenuePopUp(true)
    }, [])

    const [selectedYear, setSelectedYear] = useState(searchCriteria.getYearSearchCriteria())

    const [openRevenueSearchPopUp, setOpenRevenueSearchPopUp] = useState(false)
    const makeRevenueSearchPopUpOpen = useCallback(() => {
        setSelectedYear(searchCriteria.getYearSearchCriteria())
        setOpenRevenueSearchPopUp(true)
    }, [])
    const revenueSearchPopUpCloseHandler = useCallback(() => {
        setOpenRevenueSearchPopUp(false)
    }, [])

    const yearHandler = (event: React.ChangeEvent<HTMLInputElement>) =>
        setSelectedYear(event.target.value)

    const [openDeleteBudgetRevenuePopUp, setOpenDeleteBudgetRevenuePopUp] = useState(false)
    const makeDeleteBudgetRevenuePopUpOpen = useCallback((revenue: BudgetRevenue) => {
        setDeletableItem(revenue)
        setOpenDeleteBudgetRevenuePopUp(true)
    }, [])
    const deleteBudgetRevenuePopUpOpenCloseHandler = useCallback(() => {
        setOpenDeleteBudgetRevenuePopUp(false)
    }, [])

    let configMap = new OnlyonePortalPagesConfigMap()
    useEffect(() => budgetRevenue(), [])

    const budgetRevenue = () => {
        findBudgetRevenue(searchCriteria.getYearSearchCriteria())
            .then(({ revenues, total }) => {
                setRevenues(revenues)
                setTotal(total)
            })
    }

    const saveYear = useCallback(() => {
        searchCriteria.setYearSearchCriteria(selectedYear)
        setOpenRevenueSearchPopUp(false)
        budgetRevenue()
    }, [selectedYear])

    const budgetRevenueHandlers = {
        date: (value:moment.Moment) => setCurrentBudgetRevenueDate(value),
        amount: (event: React.ChangeEvent<HTMLInputElement>) => setCurrentBudgetRevenueAmount(event.target.value),
        note: (event: React.ChangeEvent<HTMLTextAreaElement>) => setCurrentBudgetRevenueNote(event.target.value)
    }

    const deleteItem = useCallback(() => {
        deleteBudgetRevenue(deletableItem?.id!)
            .then((response) => {
                if (response.status === 204) {
                    budgetRevenue()
                    setOpenDeleteBudgetRevenuePopUp(false)
                }
            })
    }, [deletableItem])

    const saveRevenue = useCallback(() => {
        saveBudgetRevenue({
            id: currentBudgetRevenueId,
            date: currentBudgetRevenueDate.format(FormDateFormatPattern),
            amount: currentBudgetRevenueAmount,
            note: currentBudgetRevenueNote
        }).then(response => {
            if (response.status === 201 || response.status === 204) {
                setOpenSaveBudgetRevenuePopUp(false)
                budgetRevenue()
            }
        })
    }, [currentBudgetRevenueId, currentBudgetRevenueDate, currentBudgetRevenueAmount, currentBudgetRevenueNote])

    let deleteConfirmationPopupMessages = {
        id: "deleteBudgetRevenueModal",
        title: "Delete Budget Revenue",
        message: "Are you sure of delete the Budget Revenue from the list?"
    };
    let theme = themeProvider
    return <ThemeProvider theme={theme}>
        <Paper variant="outlined">

            <Menu messages={configMap.budgetRevenue(messageRegistry).menuMessages} navBarItems={[]}>

                <OpenPopUpMenuItem icon={<Money />}
                    openPopupHandler={makeSaveBudgetRevenuePopUpOpen}
                    text={configMap.budgetExpense(messageRegistry).menuMessages.insertBudgetModal} />

                <OpenPopUpMenuItem icon={<Search />}
                    openPopupHandler={makeRevenueSearchPopUpOpen}
                    text={configMap.budgetRevenue(messageRegistry).menuMessages.changeYearModal} />

            </Menu>
            <Container>
                <SaveBudgetRevenuePopUp
                    open={openSaveBudgetRevenuePopUp}
                    handleClose={saveBudgetRevenuePopUpCloseHandler}
                    budgetRevenue={{
                        date: currentBudgetRevenueDate.format(FormDateFormatPattern),
                        amount: currentBudgetRevenueAmount,
                        note: currentBudgetRevenueNote
                    }}
                    handlers={budgetRevenueHandlers}
                    modal={configMap.budgetRevenue(messageRegistry).saveBudgetRevenueModal}
                    saveCallback={saveRevenue} />

                <DeleteBudgetRevenueConfirmationPopUp
                    open={openDeleteBudgetRevenuePopUp}
                    handleClose={deleteBudgetRevenuePopUpOpenCloseHandler}
                    modal={deleteConfirmationPopupMessages}
                    saveCallback={deleteItem} />

                <RevenueSearchPopUp
                    open={openRevenueSearchPopUp}
                    handleClose={revenueSearchPopUpCloseHandler}
                    year={selectedYear}
                    handler={yearHandler}
                    modal={configMap.budgetRevenue(messageRegistry).changeYearModal}
                    saveCallback={saveYear} />

                <BudgetRevenueContent openUpdatePopUp={makeUpdateBudgetRevenuePopUpOpen}
                    openDeletePopUp={makeDeleteBudgetRevenuePopUpOpen}
                    revenues={revenues}
                    total={total} />
            </Container>

        </Paper>
    </ThemeProvider>

}

export default BudgetRevenuePage
