package expense

import (

	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

/*
	 public final class SpentBudget {

	    private final List<BudgetExpense> budgetExpenseList;
	    private final Map<String, String> searchTags;


	    @Override
	    public boolean equals(Object o) {
	        if (this == o) return true;
	        if (o == null || getClass() != o.getClass()) return false;
	        SpentBudget that = (SpentBudget) o;
	        return Objects.equals(budgetExpenseList, that.budgetExpenseList) && Objects.equals(searchTags, that.searchTags);
	    }

	    @Override
	    public int hashCode() {
	        return Objects.hash(budgetExpenseList, searchTags);
	    }
	}
*/
type SpentBudget struct {
	BudgetExpenseList *[]BudgetExpense
	SearchTags        *map[tags.SearchTagKey]tags.SearchTagValue
}

/*
public SpentBudget(List<BudgetExpense> budgetExpenseList,

	                   List<SearchTag> searchTags) {
	    this.budgetExpenseList = budgetExpenseList;
	    this.searchTags = adaptSearchTag(searchTags);
	}
*/
func NewSpentBudget(
	budgetExpenseList *[]BudgetExpense,
	searchTags *[]tags.SearchTag) *SpentBudget {

	return &SpentBudget{
		BudgetExpenseList: budgetExpenseList,
		SearchTags:        adaptSearchTagFormListToMap(searchTags),
	}
}

/*
private Map adaptSearchTag(List<SearchTag> searchTags) {  // todo

	    return Optional.ofNullable(searchTags)
	            .map(tags ->
	                    tags.stream()
	                            .map(search -> {
	                                HashMap<String, String> map = new HashMap<>();
	                                map.put(search.key(), search.value());
	                                return map;
	                            })
	                            .reduce(new HashMap<>(), (m1, m2) -> {
	                                m1.putAll(m2);
	                                return m1;
	                            })
	            ).orElse(new HashMap<>());
	}
*/
func adaptSearchTagFormListToMap(searchTags *[]tags.SearchTag) *map[tags.SearchTagKey]tags.SearchTagValue {
	result := make(map[tags.SearchTagKey]tags.SearchTagValue)

	for _, searchTag := range *searchTags {
		result[searchTag.Key] = searchTag.Value
	}

	return &result
}

/*
public Money total() { // todo dote

	    return budgetExpenseList.stream()
	            .map(BudgetExpense::amount)
	            .reduce(Money.ZERO, Money::plus);
	}
*/
func (spentBudget *SpentBudget) Total() money.Money {
	total := money.Zero()
	for _, budgetExpense := range *spentBudget.BudgetExpenseList {
		total = total.Plus(budgetExpense.Amount)
	}
	return total
}

/*
public Map<SearchTag, Money> totalForSearchTags() { // todo done

	    return budgetExpenseList.stream()
	            .collect(groupingBy(classifier -> findSearchTagFor(classifier.userName(), classifier.tag()),
	                    Collectors.mapping(BudgetExpense::amount, Collectors.reducing(Money.ZERO, Money::plus))));
	}
*/

func (spentBudget *SpentBudget) TotalForSearchTags() map[tags.SearchTag]money.Money {
	result := make(map[tags.SearchTag]money.Money)
	for _, budgetExpense := range *spentBudget.BudgetExpenseList {
		searchTag := spentBudget.findSearchTagFor(budgetExpense.Tag.Key)
		if searchTag != nil {
			currentTotal, exists := result[*searchTag]
			if !exists {
				currentTotal = money.Zero()
			}
			result[*searchTag] = currentTotal.Plus(budgetExpense.Amount)
		}
	}
	return result
}

/*
private SearchTag findSearchTagFor(UserName userName, String searchTag) {  // todo done

	    return Optional.ofNullable(this.searchTags.get(searchTag))
	            .map(searchTagValue -> new SearchTag(searchTag, searchTagValue))
	            .orElse(null);
	}
*/
func (spentBudget *SpentBudget) findSearchTagFor(searchTagKey string) *tags.SearchTag {
	value, ok := (*spentBudget.SearchTags)[searchTagKey]
	if !ok {
		return nil
	}
	return &tags.SearchTag{
		Key:   searchTagKey,
		Value: value,
	}
}

/*
public List<DailyBudgetExpense> dailyBudgetExpenseList() {  // todo done

	    LinkedHashMap<Date, List<BudgetExpense>> result = new LinkedHashMap<>();
	    for (BudgetExpense budgetExpense : budgetExpenseList) {
	        Date key = budgetExpense.date();
	        List<BudgetExpense> value = result.getOrDefault(key, new ArrayList<>());
	        value.add(budgetExpense);
	        result.put(key, value);
	    }
	    return result.entrySet().stream()
	            .map(entry -> new DailyBudgetExpense(entry.getValue(), entry.getKey(),
	                    entry.getValue()
	                            .stream().map(BudgetExpense::amount)
	                            .reduce(Money::plus)
	                            .orElse(Money.ZERO)))
	            .collect(toList());
	}
*/
func (spentBudget *SpentBudget) DailyBudgetExpenseList() []DailyBudgetExpense {
	result := make(map[date.Date][]BudgetExpense)
	for _, budgetExpense := range *spentBudget.BudgetExpenseList {
		key := budgetExpense.Date
		value, exists := result[key]
		if !exists {
			value = []BudgetExpense{}
		}
		result[key] = append(value, budgetExpense)
	}
	resultList := make([]DailyBudgetExpense, 0)
	for key, value := range result {
		total := money.Zero()
		for _, budgetExpense := range value {
			total = total.Plus(budgetExpense.Amount)
		}
		resultList = append(resultList, DailyBudgetExpense{
			BudgetExpenseList: &value,
			Date:              key,
			Total:             total,
		})
	}
	return resultList
}

type BudgetExpense struct {
	Id       BudgetExpenseId
	UserName UserName
	Date     date.Date
	Amount   money.Money
	Note     string
	Tag      tags.SearchTag
}

type BudgetExpenseId = string
type UserName = string

type BudgetExpenseIdProvider interface {
	GenerateIdFor(budgetExpense *BudgetExpense) BudgetExpenseId
}

// public record DailyBudgetExpense(List<BudgetExpense> budgetExpenseList, Date date, Money total) {

//     @Override
//     public boolean equals(Object o) {
//         if (this == o) return true;
//         if (o == null || getClass() != o.getClass()) return false;
//         DailyBudgetExpense that = (DailyBudgetExpense) o;
//         return Objects.equals(budgetExpenseList, that.budgetExpenseList) && Objects.equals(date, that.date) && Objects.equals(total, that.total);
//     }

// }

type DailyBudgetExpense struct {
	BudgetExpenseList *[]BudgetExpense
	Date              date.Date
	Total             money.Money
}
