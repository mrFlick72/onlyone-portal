# Plan

Manages a user's Plans, each a titled, dated container of Todos.

## Language

**Plan**:
A user-owned, titled to-do list scoped to a single date. Owns zero or more Todos; deleting a Plan cascades to delete all its Todos.
_Avoid_: List, project

**Todo**:
A single actionable item inside a Plan, with its own content, date, and Status independent of the parent Plan's date.
_Avoid_: Task, item

**TodoStatus**:
The stage of a Todo's lifecycle — `TODO`, `IN_PROGRESS`, `DONE`, or `ABORTED` — governed by a fixed transition matrix enforced identically on backend and frontend; `DONE` and `ABORTED` are terminal.
_Avoid_: State
