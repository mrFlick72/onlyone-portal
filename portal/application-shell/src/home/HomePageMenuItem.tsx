import React, { ReactNode } from 'react';
import { BarChart, FormatListBulleted, Money, Person } from "@mui/icons-material";
import MenuCard from "../components/menu/MenuCard";
import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";
import { OnlyonePortalPagesConfigMap } from '../messages/OnlyonePortalPagesConfigMap';
import { MessageBundle } from '../messages/MessageRepository';

const FONT_SIZE: number = 150

interface HomePageMenuItemWrapperProps {
    content: HomePageMenuItemProps
}

interface HomePageMenuItemProps {
    titleText: string
    titleIcon: ReactNode
    body: string
    link: string
}

export type HomePageTranslationMap = {
    account: HomePageMenuItemProps,
    budget: HomePageMenuItemProps,
    revenue: HomePageMenuItemProps,
    plan: HomePageMenuItemProps,
    analytics: HomePageMenuItemProps
}

const HomePageMenuItem: React.FC<HomePageMenuItemWrapperProps> = ({ content }) => {
    return (
        <MenuCard linkTo={content.link}
            content={
                <div>
                    <div style={{ textAlign: "center" }}>
                        <h1>{content.titleText}</h1>
                        {content.titleIcon}
                    </div>
                    <h3 style={{ textAlign: "justify" }}>
                        {content.body}
                    </h3>
                </div>
            } />
    )
}

export const homeMenuContentXXX = (messageRegistry: MessageBundle): HomePageTranslationMap => {
    const configMap = new OnlyonePortalPagesConfigMap()
    return {
        account: {
            titleText: configMap.home(messageRegistry).homePage.homeMenuContent.account.titleText,
            titleIcon: <Person style={{ fontSize: FONT_SIZE }} />,
            body: configMap.home(messageRegistry).homePage.homeMenuContent.account.body,
            link: "/account/index"
        },
        budget: {
            titleText: configMap.home(messageRegistry).homePage.homeMenuContent.budget.titleText,
            titleIcon: <ShoppingCartIcon style={{ fontSize: FONT_SIZE }} />,
            body: configMap.home(messageRegistry).homePage.homeMenuContent.budget.body,
            link: "/budget/expense/index"
        },
        revenue: {
            titleText: configMap.home(messageRegistry).homePage.homeMenuContent.revenue.titleText,
            titleIcon: <Money style={{ fontSize: FONT_SIZE }} />,
            body: configMap.home(messageRegistry).homePage.homeMenuContent.revenue.body,
            link: "/budget/revenue/index"
        },
        plan: {
            titleText: configMap.home(messageRegistry).homePage.homeMenuContent.plan.titleText,
            titleIcon: <FormatListBulleted style={{ fontSize: FONT_SIZE }} />,
            body: configMap.home(messageRegistry).homePage.homeMenuContent.plan.body,
            link: "/plan/index"
        },
        analytics: {
            titleText: configMap.home(messageRegistry).homePage.homeMenuContent.analytics.titleText,
            titleIcon: <BarChart style={{ fontSize: FONT_SIZE }} />,
            body: configMap.home(messageRegistry).homePage.homeMenuContent.analytics.body,
            link: "/analytics/index"
        },
    }

}

export default HomePageMenuItem