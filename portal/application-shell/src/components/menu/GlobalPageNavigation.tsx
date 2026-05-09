import React from 'react';
import MenuItem from './MenuItem';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import { FormatListBulleted, Money, Person } from '@mui/icons-material';

const GlobalPageNavigation: React.FC = () => (
    <>
        <MenuItem icon={<ShoppingCartIcon />} text="Budget" link="/budget/expense/index" />
        <MenuItem icon={<Money />} text="Revenue" link="/budget/revenue/index" />
        <MenuItem icon={<FormatListBulleted />} text="Plans" link="/plan/index" />
        <MenuItem icon={<Person />} text="Account" link="/account/index" />
    </>
);

export default GlobalPageNavigation;
