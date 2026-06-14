---
Title: Budget Expense Analytic
Projects:
    - budget/budget-api
    - budget/analytic-api
    - portal/application-shell
Scope:
    - Forntend in portal/application-shell
    - Backend in budget/budget-api and budget/analytic-api 
---
# Frontend

The Front is developed in portal/application-shell with a dedicated page. This page should be reachable by the Home via a dedicated Tile and in the global navigation menu.
The ui is composed by two diagram:
- budget expense by year  optionally filtered by one or more search tags
- budget expense within one or more year optionally filtered by one search tag. One at time

# Backend

The backend is composed by the budget/analytic-api that is the service responsible to aggregate user budget expense data useful for the graphical representation, while the budget/budget-api is the data owner.
To improve performance as soon as ona budget expense is created, updated or deleted successfully one event is fired in a kafka topic via the BudgetExpenseEventPublisher.go abstraction implemented by the KafkaBudgetExpenseEventPublisher.go implementation. 
The budget/analytic-api instead has a kafka listener that listen on events to create, update or delete data in a postgres database under the analytic-api scope responsible to store the projected data 
