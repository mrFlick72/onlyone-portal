package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.endpoint.BudgetSearchCriteriaRepresentation;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.time.Year;
import com.onlyoneportal.budget.user.UserRepository;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import static java.util.stream.Collectors.toList;

@RestController
@RequestMapping("/budget/revenue")
public class BudgetRevenueEndPoint {

    private final BudgetRevenueRepository budgetRevenueRepository;
    private final FindBudgetRevenue findBudgetRevenue;

    private final BudgetRevenueConverter budgetRevenueConverter;
    private final UserRepository userRepository;

    public BudgetRevenueEndPoint(BudgetRevenueRepository budgetRevenueRepository,
                                 FindBudgetRevenue findBudgetRevenue, BudgetRevenueConverter budgetRevenueConverter,
                                 UserRepository userRepository) {
        this.budgetRevenueRepository = budgetRevenueRepository;
        this.findBudgetRevenue = findBudgetRevenue;
        this.budgetRevenueConverter = budgetRevenueConverter;
        this.userRepository = userRepository;
    }

    @PostMapping
    public ResponseEntity newBudgetRevenue(@RequestBody BudgetRevenueRepresentation budgetRevenueRepresentation) {
        budgetRevenueRepository.save(new BudgetRevenue(null, userRepository.currentLoggedUserName().content(),
                Date.dateFor(budgetRevenueRepresentation.date()),
                Money.moneyFor(budgetRevenueRepresentation.amount()),
                budgetRevenueRepresentation.note()));

        return ResponseEntity.status(HttpStatus.CREATED).build();
    }

    @PutMapping("/{id}")
    public ResponseEntity updateBudgetRevenue(@PathVariable("id") String id,
                                              @RequestBody BudgetRevenueRepresentation representation) {
        BudgetRevenue budgetRevenue = budgetRevenueConverter.fromRepresentationToModel(new BudgetRevenueRepresentation(
                id,
                representation.date(),
                representation.amount(),
                representation.note()
        ));
        budgetRevenueRepository.update(budgetRevenue);

        return ResponseEntity.noContent().build();
    }


    @DeleteMapping("/{id}")
    public ResponseEntity deleteBudgetRevenue(@PathVariable("id") String id) {
        budgetRevenueRepository.delete(new BudgetRevenueId(id));
        return ResponseEntity.noContent().build();
    }

    @GetMapping
    public ResponseEntity findAllBudgetRevenue(@RequestParam("q") BudgetSearchCriteriaRepresentation searchCriteria) {
        Year year = Year.of(searchCriteria.year());

        return ResponseEntity.ok(findBudgetRevenue.findBy(year)
                .stream().map(budgetRevenueConverter::fromDomainToRepresentation)
                .collect(toList()));
    }

}
