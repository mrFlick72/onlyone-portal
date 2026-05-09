import React, {ReactNode} from 'react';
import {FormatListBulleted, Money, Person} from "@mui/icons-material";
import MenuCard from "../components/menu/MenuCard";
import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";

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

const HomePageMenuItem: React.FC<HomePageMenuItemWrapperProps> = ({content}) => {
    return (
        <MenuCard linkTo={content.link}
                  content={
                      <div>
                          <div style={{textAlign: "center"}}>
                              <h1>{content.titleText}</h1>
                              {content.titleIcon}
                          </div>
                          <h3 style={{textAlign: "justify"}}>
                              {content.body}
                          </h3>
                      </div>
                  }/>
    )
}

export const homeMenuContent = {
    account: {
        titleText: "Account Management Section",
        titleIcon: <Person style={{fontSize: FONT_SIZE}}/>,
        body: "In this section you can manage your account details",
        link: "/account/index"
    },
    budget: {
        titleText: "Budget Management Section",
        titleIcon: <ShoppingCartIcon style={{fontSize: FONT_SIZE}}/>,
        body: "In this section you can manage your budget",
        link: "/budget/expense/index"
    },
    revenue: {
        titleText: "Revenue Management Section",
        titleIcon: <Money style={{fontSize: FONT_SIZE}}/>,
        body: "In this section you can manage your revenue",
        link: "/budget/revenue/index"
    },
    plan: {
        titleText: "Plan Management Section",
        titleIcon: <FormatListBulleted style={{fontSize: FONT_SIZE}}/>,
        body: "In this section you can manage your plans and todos",
        link: "/plan/index"
    },

}
export default HomePageMenuItem