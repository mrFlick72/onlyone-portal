type BudgetRevenue = {
  id?: string;
  amount: string;
  date: string;
  note: string;
};

export type BudgetRevenueList = {
  revenues: BudgetRevenue[];
  total: string;
};

export default BudgetRevenue;