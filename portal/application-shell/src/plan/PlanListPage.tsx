import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { Container, Paper, ThemeProvider } from "@mui/material";
import { PlaylistAdd } from "@mui/icons-material";
import themeProvider from "../theme/ThemeProvider";
import Menu from "../components/menu/Menu";
import OpenPopUpMenuItem from "../components/menu/OpenPopUpMenuItem";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap";
import Plan from "./domain/Plan";
import { createPlan, deletePlan, getAllPlans } from "./domain/PlanRepository";
import PlanListContent from "./PlanListContent";
import SavePlanPopUp from "./SavePlanPopUp";
import DeletePlanConfirmationPopUp from "./DeletePlanConfirmationPopUp";

interface PlanListPageProps {
    messageRegistry: any;
}

const PlanListPage: React.FC<PlanListPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const [plans, setPlans] = useState<Plan[]>([])

    const [title, setTitle] = useState("")
    const [date, setDate] = useState<moment.Moment>(moment())
    const [openSavePopUp, setOpenSavePopUp] = useState(false)

    const [deletable, setDeletable] = useState<Plan | null>(null)
    const [openDeletePopUp, setOpenDeletePopUp] = useState(false)

    const refresh = useCallback(() => {
        getAllPlans().then(setPlans)
    }, [])

    useEffect(() => { refresh() }, [refresh])

    const openSave = useCallback(() => {
        setTitle("")
        setDate(moment())
        setOpenSavePopUp(true)
    }, [])

    const closeSave = useCallback(() => setOpenSavePopUp(false), [])

    const save = useCallback(() => {
        createPlan({
            title,
            date: date.format(ApiDateFormatPattern),
        }).then(result => {
            if (result) {
                setOpenSavePopUp(false)
                refresh()
            }
        })
    }, [title, date, refresh])

    const openDelete = useCallback((plan: Plan) => {
        setDeletable(plan)
        setOpenDeletePopUp(true)
    }, [])

    const closeDelete = useCallback(() => setOpenDeletePopUp(false), [])

    const confirmDelete = useCallback(() => {
        if (!deletable) {
            return
        }
        deletePlan(deletable.id).then(response => {
            if (response.status === 204) {
                setOpenDeletePopUp(false)
                refresh()
            }
        })
    }, [deletable, refresh])

    const openDetail = useCallback((plan: Plan) => {
        window.location.href = `/plan/detail?id=${encodeURIComponent(plan.id)}`
    }, [])

    const planMessages = configMap.plan(messageRegistry)
    const common = configMap.common(messageRegistry)

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={planMessages.menuMessages} navBarItems={[]}>
                <OpenPopUpMenuItem
                    icon={<PlaylistAdd />}
                    openPopupHandler={openSave}
                    text={planMessages.menuMessages.newPlan} />
            </Menu>
            <Container>
                <SavePlanPopUp
                    open={openSavePopUp}
                    handleClose={closeSave}
                    plan={{
                        title,
                        date: date.format(FormDateFormatPattern),
                    }}
                    handlers={{
                        title: (event) => setTitle(event.target.value),
                        date: (value) => setDate(value),
                    }}
                    saveCallback={save}
                    modal={planMessages.savePlanModal}
                    formMessages={{ title: common.form.title, date: common.form.date }} />

                <DeletePlanConfirmationPopUp
                    open={openDeletePopUp}
                    handleClose={closeDelete}
                    saveCallback={confirmDelete}
                    modal={planMessages.deletePlanModal} />

                <PlanListContent
                    plans={plans}
                    openDetail={openDetail}
                    openDelete={openDelete}
                    messages={{
                        headers: {
                            date: common.table.date,
                            title: common.table.title,
                            todos: common.table.todos,
                            options: common.table.options
                        },
                        actions: { open: common.action.open, delete: common.action.delete }
                    }} />
            </Container>
        </Paper>
    </ThemeProvider>
}

export default PlanListPage;
