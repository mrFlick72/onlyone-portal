import React, { useCallback, useEffect, useState } from "react"
import moment from "moment";
import { Container, Paper, ThemeProvider } from "@mui/material";
import { ArrowBack, Save } from "@mui/icons-material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import MenuItem from "../../components/menu/MenuItem";
import OpenPopUpMenuItem from "../../components/menu/OpenPopUpMenuItem";
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

    useEffect(() => {
        getSearchTagRegistry("expense").then(setSearchTagRegistry)
    }, [])

    const save = useCallback(() => {
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
                window.location.href = "/budget/scheduled-expense/index"
            }
        })
    }, [description, amount, note, searchTags, day, month, hasEndDate, endDate])

    const detailMessages = configMap.scheduledExpenseDetail(messageRegistry)

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={detailMessages.menuMessages} navBarItems={[]}>
                <MenuItem
                    icon={<ArrowBack />}
                    text={detailMessages.menuMessages.backToList}
                    link="/budget/scheduled-expense/index" />
                <OpenPopUpMenuItem
                    icon={<Save />}
                    openPopupHandler={save}
                    text={detailMessages.saveButtonLabel} />
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
            </Container>
        </Paper>
    </ThemeProvider>
}

export default ScheduledExpenseDetailPage;
