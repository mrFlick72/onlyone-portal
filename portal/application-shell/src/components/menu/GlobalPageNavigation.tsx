import React from 'react';
import MenuItem from './MenuItem';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import { FormatListBulleted, Money, Person } from '@mui/icons-material';

type GlobalPageNavigationProps = {
    messages: {
        budget: string,
        revenue: string,
        plans: string,
        account: string
    }
}

const GlobalPageNavigation: React.FC<GlobalPageNavigationProps> = ({ messages }) => (
    <>
        <MenuItem icon={<ShoppingCartIcon />} text={messages.budget} link="/budget/expense/index" />
        <MenuItem icon={<Money />} text={messages.revenue} link="/budget/revenue/index" />
        <MenuItem icon={<FormatListBulleted />} text={messages.plans} link="/plan/index" />
        <MenuItem icon={<Person />} text={messages.account} link="/account/index" />
    </>
);

export default GlobalPageNavigation;
