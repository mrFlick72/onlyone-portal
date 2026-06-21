import React, { FC, useState } from 'react'
import { OnlyonePortalPagesConfigMap } from "../../messages/OnlyonePortalPagesConfigMap";
import { Box, Container, Paper, Tab, Tabs, ThemeProvider, Typography } from "@mui/material";
import themeProvider from "../../theme/ThemeProvider";
import Menu from "../../components/menu/Menu";
import TabPanel from "../../components/layout/TabPanel";
import ScopedTagManager from "./ScopedTagManager";
import { TagScope } from "./domain/SearchTagRepository";
import { MessageBundle } from '../../messages/MessageRepository';

interface SearchTagsPageProps {
    messageRegistry: MessageBundle;
}

// Tabs are keyed by scope; the index is the position in this ordered list.
const SCOPES: TagScope[] = ["expense", "revenue"];

// Pick the initial tab from the ?scope= query param so the Expense / Revenue
// menu entries land on their matching tab. Unknown/absent scope -> Expense.
function initialTabFromQuery(): number {
    const scope = new URLSearchParams(window.location.search).get("scope");
    const index = SCOPES.indexOf(scope as TagScope);
    return index >= 0 ? index : 0;
}

const SearchTagsPage: FC<SearchTagsPageProps> = ({ messageRegistry }) => {
    const configMap = new OnlyonePortalPagesConfigMap()
    const messages = configMap.searchTags(messageRegistry)
    const tableMessages = {
        description: configMap.common(messageRegistry).table.description,
        operation: configMap.common(messageRegistry).table.operation
    }

    const [tabPanel, setTabPanel] = useState<number>(initialTabFromQuery)

    const handleTabPanelChange = (_event: React.SyntheticEvent, newValue: number) => {
        setTabPanel(newValue)
    }

    const tabLabels = [messages.tabs.expense, messages.tabs.revenue]

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={messages.menuMessages} navBarItems={[]} />

            <Container>
                <Typography variant="h5" sx={{ mt: 2, mb: 1 }}>{messages.heading}</Typography>

                <Tabs
                    value={tabPanel}
                    onChange={handleTabPanelChange}
                    textColor="secondary"
                    indicatorColor="secondary"
                    aria-label="tag scope tabs"
                >
                    {SCOPES.map((scope, index) => (
                        <Tab key={scope} value={index} label={tabLabels[index]} />
                    ))}
                </Tabs>

                <Box>
                    {SCOPES.map((scope, index) => (
                        <TabPanel key={scope} value={tabPanel} index={index}>
                            <ScopedTagManager
                                scope={scope}
                                formMessages={messages.form}
                                tableMessages={tableMessages} />
                        </TabPanel>
                    ))}
                </Box>
            </Container>
        </Paper>
    </ThemeProvider>
}
export default SearchTagsPage
