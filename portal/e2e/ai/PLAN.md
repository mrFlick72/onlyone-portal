Follow @prerequisite.md directive
Then verify the following AC

Login in the system like described in the @LOGIN.md

- AC1: Click the Plan Management tiles should go to the plan home page
  - We expect the url is http://local.onlyone-portal.com:8070/plan/index
  - We expect the UI is like in assertions/plan/empty_plan_list.png
  - Mark the AC failed otherwise
- AC2: Create a new Plan with the title "Smock Test" should create a new line in the plan list
  - We expect the plan list updated like in the assertions/plan/new_plan.png 
  - We should have only one "Smock Test" row
  - Mark the AC failed otherwise
- AC3: Clicking to the Open button we should go to the Plan details page
  - We expect the UI like assertions/plan/empty_plan.png
  - Mark the AC failed otherwise
- CLEANUP: Go back to the Plan list page and delete the all the plans

Give me a report on AC outcome