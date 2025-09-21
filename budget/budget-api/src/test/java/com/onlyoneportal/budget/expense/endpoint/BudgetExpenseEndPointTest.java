package com.onlyoneportal.budget.expense.endpoint;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.action.*;
import com.onlyoneportal.budget.expense.converter.BudgetExpenseConverter;
import com.onlyoneportal.budget.expense.converter.SpentBudgetConverter;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.expense.model.SpentBudget;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.time.Month;
import com.onlyoneportal.budget.time.Year;
import com.onlyoneportal.budget.user.UserName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.web.servlet.MockMvc;

import java.util.UUID;

import static com.onlyoneportal.budget.support.BudgetFixture.emptyBudgetExpenseId;
import static java.util.Arrays.asList;
import static org.mockito.BDDMockito.given;
import static org.mockito.Mockito.verify;
import static org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors.csrf;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(BudgetExpenseEndPoint.class) //todo https://github.com/spring-projects/spring-boot/issues/32195
public class BudgetExpenseEndPointTest {

    @Autowired
    private MockMvc mockMvc;

    private final ObjectMapper objectMapper = new ObjectMapper();
    private final SpentBudgetRepresentation spentBudgetRepresentation =
            new SpentBudgetRepresentation("21.00", asList(new DailyBudgetExpenseRepresentation(asList(new BudgetExpenseRepresentation("1", "20/01/2018", "10.50", "", "super-market", "Super Market")), "20/01/2018", "10.50")), asList());

    @MockBean
    private JwtDecoder decoder;

    @MockBean
    private FindSpentBudget findSpentBudget;

    @MockBean
    private BudgetExpenseRepository budgetExpenseRepository;

    @MockBean
    private SpentBudgetConverter spentBudgetConverter;

    @MockBean
    private BudgetExpenseConverter budgetExpenseConverter;

    @MockBean
    private CreateBudgetExpense createBudgetExpense;

    @MockBean
    private UpdateBudgetExpenseDetails updateBudgetExpenseDetails;

    @MockBean
    private DeleteBudgetExpense deleteBudgetExpense;

    @Test
    @WithMockUser
    public void getBudgetExpenseListOfJanuary() throws Exception {
        String expectedSpentBudgetRepresentation =
                objectMapper.writeValueAsString(spentBudgetRepresentation);

        SpentBudget toBeReturned = new SpentBudget(asList(new BudgetExpense(emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("20/01/2018"), Money.moneyFor("10.50"), "", "super-market")), null);

        given(findSpentBudget.findBy(Month.of(1), Year.of(2018), asList()))
                .willReturn(toBeReturned);

        given(spentBudgetConverter.domainToRepresentationModel(toBeReturned))
                .willReturn(spentBudgetRepresentation);

        mockMvc.perform(get("/budget/expense")
                        .with(csrf())
                        .param("q", "month=1;year=2018;searchTag=")
                        .accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isOk())
                .andExpect(content().json(expectedSpentBudgetRepresentation));
    }

    @Test
    @WithMockUser
    public void getBudgetExpenseListOfJanuaryV2() throws Exception {
        String expectedSpentBudgetRepresentation =
                objectMapper.writeValueAsString(spentBudgetRepresentation);

        SpentBudget toBeReturned = new SpentBudget(asList(new BudgetExpense(emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("20/01/2018"), Money.moneyFor("10.50"), "", "super-market")), null);

        given(findSpentBudget.findBy(Month.of(1), Year.of(2018), asList()))
                .willReturn(toBeReturned);

        given(spentBudgetConverter.domainToRepresentationModel(toBeReturned))
                .willReturn(spentBudgetRepresentation);

        mockMvc.perform(put("/budget/expense")
                        .with(csrf())
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new BudgetSearchCriteriaRepresentation(1, 2018, asList()))))
                .andExpect(status().isOk())
                .andExpect(content().json(expectedSpentBudgetRepresentation));
    }

    @Test
    @WithMockUser
    public void deleteBudgetExpense() throws Exception {
        String id = UUID.randomUUID().toString();

        mockMvc.perform(delete("/budget/expense/{id}", id)
                        .with(csrf()))
                .andExpect(status().isNoContent());

        verify(deleteBudgetExpense).delete(new BudgetExpenseId(id));
    }

    @Test
    @WithMockUser
    public void insertBudgetExpense() throws Exception {
        BudgetExpenseRepresentation budgetExpenseRepresentation =
                new BudgetExpenseRepresentation(null, "20/01/2018", "10.50", "", "super-market", null);

        String budgetExpenseRepresentationAsJsonString =
                objectMapper.writeValueAsString(budgetExpenseRepresentation);

        NewBudgetExpenseRequest newBudgetExpenseRequest = new NewBudgetExpenseRequest(Date.dateFor("20/01/2018"), Money.moneyFor("10.50"), "", "super-market");
        given(budgetExpenseConverter.newBudgetExpenseRequestFromRepresentation(budgetExpenseRepresentation))
                .willReturn(newBudgetExpenseRequest);

        given(createBudgetExpense.newBudgetExpense(newBudgetExpenseRequest))
                .willReturn(new BudgetExpense(emptyBudgetExpenseId(), new UserName("USER"), null, null, null, null));

        mockMvc.perform(post("/budget/expense")
                        .contentType(MediaType.APPLICATION_JSON)
                        .with(csrf())
                        .content(budgetExpenseRepresentationAsJsonString))
                .andExpect(status().isCreated());

        verify(budgetExpenseConverter).newBudgetExpenseRequestFromRepresentation(budgetExpenseRepresentation);
        verify(createBudgetExpense).newBudgetExpense(newBudgetExpenseRequest);
    }

    @Test
    @WithMockUser
    public void updateBudgetExpense() throws Exception {
        BudgetExpenseRepresentation budgetExpenseRepresentation =
                new BudgetExpenseRepresentation("ID", "20/01/2018", "10.50", "", "super-market", null);

        String budgetExpenseRepresentationAsJsonString = objectMapper.writeValueAsString(new BudgetExpenseRepresentation("ID123", "20/01/2018", "10.50", "", "super-market", null));

        BudgetExpense newBudgetExpenseRequest = new BudgetExpense(new BudgetExpenseId("ID"), new UserName("USER"), Date.dateFor("20/01/2018"), Money.moneyFor("10.50"), "", "super-market");
        given(budgetExpenseConverter.representationModelToDomainModel(budgetExpenseRepresentation))
                .willReturn(newBudgetExpenseRequest);


        mockMvc.perform(put("/budget/expense/ID")
                        .contentType(MediaType.APPLICATION_JSON)
                        .with(csrf())
                        .content(budgetExpenseRepresentationAsJsonString))
                .andExpect(status().isNoContent());

        verify(budgetExpenseConverter).representationModelToDomainModel(budgetExpenseRepresentation);
        verify(updateBudgetExpenseDetails).updateWithoutAttachment(newBudgetExpenseRequest);
    }
}
