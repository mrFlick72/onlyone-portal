import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { Box, Container, Paper, ThemeProvider, Typography } from "@mui/material";
import { ArrowBack, NoteAdd } from "@mui/icons-material";
import themeProvider from "../theme/ThemeProvider";
import Menu from "../components/menu/Menu";
import OpenPopUpMenuItem from "../components/menu/OpenPopUpMenuItem";
import MenuItem from "../components/menu/MenuItem";
import { ApiDateFormatPattern, FormDateFormatPattern } from "../components/form/FormDatePicker";
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap";
import Plan, { Todo, TodoStatus } from "./domain/Plan";
import { addTodo, changeTodoStatus, getPlan, removeTodo, updateTodo } from "./domain/PlanRepository";
import TodoContent from "./TodoContent";
import SaveTodoPopUp from "./SaveTodoPopUp";
import DeleteTodoConfirmationPopUp from "./DeleteTodoConfirmationPopUp";
import ChangeTodoStatusPopUp from "./ChangeTodoStatusPopUp";

interface PlanDetailPageProps {
    messageRegistry: any;
}

const PlanDetailPage: React.FC<PlanDetailPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const planId = new URLSearchParams(window.location.search).get("id") ?? ""

    const [plan, setPlan] = useState<Plan | null>(null)

    const [todoId, setTodoId] = useState("")
    const [todoContent, setTodoContent] = useState("")
    const [todoDate, setTodoDate] = useState<moment.Moment>(moment())
    const [openSavePopUp, setOpenSavePopUp] = useState(false)

    const [deletable, setDeletable] = useState<Todo | null>(null)
    const [openDeletePopUp, setOpenDeletePopUp] = useState(false)

    const [statusTarget, setStatusTarget] = useState<Todo | null>(null)
    const [openStatusPopUp, setOpenStatusPopUp] = useState(false)

    const refresh = useCallback(() => {
        if (!planId) {
            return
        }
        getPlan(planId).then(setPlan)
    }, [planId])

    useEffect(() => { refresh() }, [refresh])

    const openCreate = useCallback(() => {
        setTodoId("")
        setTodoContent("")
        setTodoDate(moment())
        setOpenSavePopUp(true)
    }, [])

    const openUpdate = useCallback((todo: Todo) => {
        setTodoId(todo.id)
        setTodoContent(todo.content)
        setTodoDate(moment(todo.date, ApiDateFormatPattern))
        setOpenSavePopUp(true)
    }, [])

    const closeSave = useCallback(() => setOpenSavePopUp(false), [])

    const save = useCallback(() => {
        const payload = {
            content: todoContent,
            date: todoDate.format(ApiDateFormatPattern),
        }
        const action = todoId
            ? updateTodo(planId, todoId, payload)
            : addTodo(planId, payload)
        action.then(response => {
            if (response.status === 201 || response.status === 204) {
                setOpenSavePopUp(false)
                refresh()
            }
        })
    }, [planId, todoId, todoContent, todoDate, refresh])

    const openDelete = useCallback((todo: Todo) => {
        setDeletable(todo)
        setOpenDeletePopUp(true)
    }, [])

    const closeDelete = useCallback(() => setOpenDeletePopUp(false), [])

    const confirmDelete = useCallback(() => {
        if (!deletable) {
            return
        }
        removeTodo(planId, deletable.id).then(response => {
            if (response.status === 204) {
                setOpenDeletePopUp(false)
                refresh()
            }
        })
    }, [planId, deletable, refresh])

    const openChangeStatus = useCallback((todo: Todo) => {
        setStatusTarget(todo)
        setOpenStatusPopUp(true)
    }, [])

    const closeChangeStatus = useCallback(() => setOpenStatusPopUp(false), [])

    const applyStatus = useCallback((status: TodoStatus) => {
        if (!statusTarget) {
            return
        }
        changeTodoStatus(planId, statusTarget.id, status).then(response => {
            if (response.status === 204) {
                setOpenStatusPopUp(false)
                refresh()
            }
        })
    }, [planId, statusTarget, refresh])

    const planMessages = configMap.plan(messageRegistry)
    const detailMessages = configMap.planDetail(messageRegistry)
    const common = configMap.common(messageRegistry)

    const formattedDate = plan
        ? moment(plan.date, ApiDateFormatPattern).format(FormDateFormatPattern)
        : ""

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={planMessages.menuMessages} navBarItems={[]}>
                <MenuItem
                    icon={<ArrowBack />}
                    text={detailMessages.menuMessages.backToList}
                    link="/plan/index" />
                <OpenPopUpMenuItem
                    icon={<NoteAdd />}
                    openPopupHandler={openCreate}
                    text={detailMessages.menuMessages.newTodo} />
            </Menu>
            <Container>
                <Box sx={{ my: 2 }}>
                    <Typography variant="h5">{plan?.title ?? ""}</Typography>
                    <Typography variant="subtitle1" color="text.secondary">{formattedDate}</Typography>
                </Box>

                <SaveTodoPopUp
                    open={openSavePopUp}
                    handleClose={closeSave}
                    todo={{
                        content: todoContent,
                        date: todoDate.format(FormDateFormatPattern),
                    }}
                    handlers={{
                        content: (event) => setTodoContent(event.target.value),
                        date: (value) => setTodoDate(value),
                    }}
                    saveCallback={save}
                    modal={todoId ? detailMessages.updateTodoModal : detailMessages.saveTodoModal}
                    formMessages={{ date: common.form.date, content: common.form.content }} />

                <DeleteTodoConfirmationPopUp
                    open={openDeletePopUp}
                    handleClose={closeDelete}
                    saveCallback={confirmDelete}
                    modal={detailMessages.deleteTodoModal} />

                <ChangeTodoStatusPopUp
                    open={openStatusPopUp}
                    todo={statusTarget}
                    handleClose={closeChangeStatus}
                    onSelect={applyStatus}
                    messages={detailMessages.changeStatusModal}
                    statusMessages={detailMessages.statusMessages} />

                <TodoContent
                    todos={plan?.todos ?? []}
                    openUpdate={openUpdate}
                    openDelete={openDelete}
                    openChangeStatus={openChangeStatus}
                    messages={{
                        headers: {
                            date: common.table.date,
                            content: common.table.content,
                            status: common.table.status,
                            options: common.table.options
                        },
                        actions: {
                            changeStatus: common.action.changeStatus,
                            edit: common.action.edit,
                            delete: common.action.delete
                        },
                        status: detailMessages.statusMessages
                    }} />
            </Container>
        </Paper>
    </ThemeProvider>
}

export default PlanDetailPage;
