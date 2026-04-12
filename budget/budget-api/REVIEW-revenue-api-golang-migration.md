# Code Review: `revenue-api-golang-migration`

## Overview

This branch migrates the revenue API from the standalone Python FastAPI service (`budget/revenue-api`) into the Go `budget-api` service. It includes:

- **Repackaging** existing expense code from `adapter/budget/dynamodb/` to `adapter/budget/expense/dynamodb/` and `web/` to `web/budget/expense/`
- **New revenue domain, adapter, web, and config layers** following the hexagonal architecture
- **Frontend update** to point revenue API calls at the unified `budget-api` endpoint
- Comprehensive tests at all layers

The architecture is clean and the DynamoDB key scheme correctly preserves Python-era data compatibility.

---

## Critical Issues

### 1. Committed binary: `budget/budget-api/app` (40 MB) DONE

A compiled Go binary was committed to the repo. There is no `.gitignore` for `budget-api/`. This bloats the repository permanently (even after removal, it stays in git history). Should be removed from tracking and a `.gitignore` added.

### 2. Committed log file: `budget/budget-api/test/logs.log` DONE

Contains runtime output from a local test run. Should be git-ignored.

### 3. `CreateRevenue.Execute` silently discards `GetCurrentUser` error DONE

**File:** `domain/budget/revenue/actions.go:16`

```go
user, _ := security.GetCurrentUser(ctx)
revenue.UserName = *user.UserName  // nil-pointer panic if error
```

`UpdateRevenue` and `DeleteRevenue` both check this error properly. This is a runtime crash if the context has no user.

**Fix:** handle the error like the other actions do:

```go
func (action *CreateRevenue) Execute(ctx context.Context, revenue *Revenue) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	revenue.UserName = *user.UserName
	return action.Repository.Save(ctx, revenue)
}
```

### 4. `RevenueRepresentationToDomainModel` silently discards parse errors DONE

**File:** `web/budget/revenue/converter.go:12-14`

```go
d, _ := date.DateFor(rep.Date)
m, _ := money.MoneyFor(rep.Amount)
return &domainrevenue.Revenue{
    ...
    Date:   *d,   // nil-pointer panic on bad input
```

Malformed date or amount from user input will crash the server. Should return an error or the calling endpoint should validate before conversion.

**Fix:** return an error from the converter:

```go
func RevenueRepresentationToDomainModel(rep RevenueRepresentation) (*domainrevenue.Revenue, error) {
	d, err := date.DateFor(rep.Date)
	if err != nil {
		return nil, err
	}
	m, err := money.MoneyFor(rep.Amount)
	if err != nil {
		return nil, err
	}
	return &domainrevenue.Revenue{
		Id:     domainrevenue.RevenueId(rep.Id),
		Date:   *d,
		Amount: m,
		Note:   rep.Note,
	}, nil
}
```

Then update the endpoint handlers to check the error and return 400.

### 5. Endpoint handlers ignore facade errors

**File:** `web/budget/revenue/endpoint.go`

- Line 49: `facade.CreateRevenue(ctx, domainModel)` -- error discarded, returns `201` regardless
- Line 65: `facade.UpdateRevenue(ctx, domainModel)` -- error discarded, returns `204` regardless
- Line 71: `facade.DeleteRevenue(ctx, c.Param("id"))` -- error discarded, returns `204` regardless

Authorization failures, not-found errors, and DynamoDB errors all silently succeed from the client's perspective. Compare with the GET handler (line 29-34) which correctly checks errors.

**Fix for POST (create):**

```go
r.POST("/api/budget/revenue", func(c *gin.Context) {
    var representation RevenueRepresentation

    ctx := ContextFactoryConverter.CreateContextFromGin(c)
    if err := c.ShouldBindJSON(&representation); err != nil {
        logger.LogErrorfFor("Error binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    domainModel := RevenueRepresentationToDomainModel(representation)
    if err := facade.CreateRevenue(ctx, domainModel); err != nil {
        logger.LogErrorfFor("Error creating revenue: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusCreated)
})
```

Apply the same pattern to PUT and DELETE handlers.

---

## Medium Issues

### 6. `FindFor` and `FindByDateRange` use `context.TODO()` instead of the passed `ctx`

**File:** `adapter/budget/revenue/dynamodb/dynamo_db_revenue_repository.go:60,120`

```go
result, err := repository.Client.Query(context.TODO(), input)
```

The caller provides a `ctx` (carrying user info, cancellation), but the actual DynamoDB call ignores it. `Save` correctly uses `ctx` (line 182). This breaks context propagation for timeouts/tracing.

**Fix:** replace `context.TODO()` with `ctx` in both `FindFor` (line 60) and `FindByDateRange` (line 120).

### 7. `FindByDateRange` type-asserts to concrete `DynamoDbRevenueIdProvider`

**File:** `adapter/budget/revenue/dynamodb/dynamo_db_revenue_repository.go:104`

```go
idProvider, _ := repository.RevenueIdProvider.(*DynamoDbRevenueIdProvider)
```

The `ok` bool is discarded. If a different `RevenueIdProvider` is injected, `idProvider` is nil and the next line panics. This mirrors the expense adapter pattern so it is consistent, but worth noting as a latent crash.

---

## Low Priority / Style

### 8. `range_key` uses snake_case

Go convention is `rangeKey`. Used throughout `dynamo_db_revenue_repository.go`. Not a bug but inconsistent with Go idiom.

### 9. Trailing whitespace in `converter.go:20-21`

Blank line with trailing whitespace inside the struct literal in `RevenueRepresentationToDomainModel`.

---

## Positive Observations

- Hexagonal architecture is faithfully followed -- domain has no adapter/web imports
- DynamoDB key scheme preserves Python-era data (no migration needed)
- Dual-layer ownership enforcement (domain check + DynamoDB `ConditionExpression`) matches the expense pattern
- Test coverage is comprehensive: domain unit tests, adapter integration tests (LocalStack), web endpoint tests
- `parseYearQueryParam` correctly preserves the `?q=year=YYYY` wire format for frontend compat
- The expense repackaging renames are clean with no logic changes
