export type RevenueTag = {
  tagKey: string;
  tagValue: string;
};

type BudgetRevenue = {
  id?: string;
  amount: string;
  date: string;
  note: string;
  tags: RevenueTag[];
};

export type BudgetRevenueList = {
  revenues: BudgetRevenue[];
  total: string;
};

export default BudgetRevenue;