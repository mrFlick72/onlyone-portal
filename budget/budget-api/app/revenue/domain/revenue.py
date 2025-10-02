
# public record BudgetRevenue(BudgetRevenueId id, String userName, Date registrationDate, Money amount, String note) {

from money.domain import Money
class Revenue:
    
    def __init__(self, id, user_name:str, registration_date, amount : Money, note:str):
        pass 