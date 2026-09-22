import React, { useCallback, useEffect, useState } from "react"
import { Container, Paper, ThemeProvider } from "@mui/material";
import { EventRepeat } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { MessageBundle } from "../../messages/MessageRepository";
import ScheduledExpense from "./domain/ScheduledExpense";
import { getAllScheduledExpenses } from "./domain/ScheduledExpenseRepository";
import ScheduledExpenseListContent from "./ScheduledExpenseListContent";

type ScheduledExpenseListPageProps = {
    messageRegistry: MessageBundle;
}

// #51 ships list + link-to-details only. Delete and pause/resume (row
// actions) land in #53/#54 once their backend endpoints exist — see the
// parent issue and ADR 0005.
const ScheduledExpenseListPage: React.FC<ScheduledExpenseListPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const [scheduledExpenses, setScheduledExpenses] = useState<ScheduledExpense[]>([])

    const refresh = useCallback(() => {
        getAllScheduledExpenses().then(setScheduledExpenses)
    }, [])

    useEffect(() => { refresh() }, [refresh])

    const openDetail = useCallback((scheduledExpense: ScheduledExpense) => {
        window.location.href = `/budget/scheduled-expense/detail?id=${encodeURIComponent(scheduledExpense.id ?? "")}`
    }, [])

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
                <ScheduledExpenseListContent
                    scheduledExpenses={scheduledExpenses}
                    openDetail={openDetail}
                    messages={listMessages.content} />
            </Container>
        </Paper>
    </ThemeProvider>
}

export default ScheduledExpenseListPage;
