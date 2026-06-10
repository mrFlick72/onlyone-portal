import React, { useEffect, useState } from 'react';
import themeProvider from "../theme/ThemeProvider";
import { Container, Paper, ThemeProvider } from "@mui/material";
import Menu from "../components/menu/Menu";
import { OnlyonePortalPagesConfigMap } from "../messages/OnlyonePortalPagesConfigMap";
import { isAuthenticated } from "../auth/Authenticator";
import { getAllMessageRegistry, MessageBundle } from "../messages/MessageRepository";
import HomePageMenuItem, { homeMenuContentXXX as homeMenuContentFor } from "./HomePageMenuItem";
import MenuCardContainer from "../components/menu/MenuCardContainer";
import BudgetPageMenuItem from '../components/menu/BudgetPageMenuItem';
import BudgetRevenuePageMenuItem from '../components/menu/BudgetRevenuePageMenuItem';
import AccountPageMenuItem from '../components/menu/AccountPageMenuItem';
import PlanPageMenuItem from '../components/menu/PlanPageMenuItem';
import AnalyticsPageMenuItem from '../components/menu/AnalyticsPageMenuItem';

const HomePage: React.FC = () => {
    const configMap = new OnlyonePortalPagesConfigMap();
    const [messageRegistry, setMessageRegistry] = useState<MessageBundle>({})

    useEffect(() => {
        isAuthenticated().then()
        setMessageRegistry(getAllMessageRegistry())
    }, []); // Or [] if effect doesn't need props or state

    const homeMenuContent = homeMenuContentFor(messageRegistry)

    return <ThemeProvider theme={themeProvider}>
        <Paper variant="outlined">
            <Menu messages={configMap.searchTags(messageRegistry).menuMessages}
                navBarItems={[]}
                showGlobalNav={false}>
                <BudgetPageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.budgetPageLabel} />
                <BudgetRevenuePageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.revenuePageLabel} />
                <PlanPageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.planPageLabel} />
                <AnalyticsPageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.analyticsPageLabel} />
                <AccountPageMenuItem text={configMap.budgetExpense(messageRegistry).menuMessages.userProfileLabel} />
            </Menu>

            <Container>
                <MenuCardContainer>
                    <HomePageMenuItem content={homeMenuContent.budget} />
                    <HomePageMenuItem content={homeMenuContent.revenue} />
                    <HomePageMenuItem content={homeMenuContent.plan} />
                    <HomePageMenuItem content={homeMenuContent.analytics} />
                    <HomePageMenuItem content={homeMenuContent.account} />
                </MenuCardContainer>
            </Container>
        </Paper>
    </ThemeProvider>
}

export default HomePage