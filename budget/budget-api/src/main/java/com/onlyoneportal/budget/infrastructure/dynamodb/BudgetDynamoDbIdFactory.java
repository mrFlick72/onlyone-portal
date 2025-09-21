package com.onlyoneportal.budget.infrastructure.dynamodb;

import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;

public interface BudgetDynamoDbIdFactory<ID, OUT> {

    ID budgetIdFrom(OUT budget);

    String partitionKeyFrom(Date date, UserName userName);

    String partitionKeyFrom(ID id);

    String rangeKeyFrom(ID id);


}