import React from 'react';
import MenuItem from './MenuItem';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import { BarChart, FormatListBulleted, Money, Person } from '@mui/icons-material';
import { GlobalNavMessageBundle } from '../../messages/MessageBundles';

type GlobalPageNavigationProps = {
    messages: GlobalNavMessageBundle
}

const GlobalPageNavigation: React.FC<GlobalPageNavigationProps> = ({ messages }) => (
    <>
        <MenuItem icon={<ShoppingCartIcon />} text={messages.budget} link="/budget/expense/index" />
        <MenuItem icon={<Money />} text={messages.revenue} link="/budget/revenue/index" />
        <MenuItem icon={<FormatListBulleted />} text={messages.plans} link="/plan/index" />
        <MenuItem icon={<BarChart />} text={messages.analytics} link="/analytics/index" />
        <MenuItem icon={<Person />} text={messages.account} link="/account/index" />
    </>
);

export default GlobalPageNavigation;
