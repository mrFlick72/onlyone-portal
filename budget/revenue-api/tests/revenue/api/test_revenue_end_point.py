import pytest
from app.server import app
from pytest_mock import MockerFixture
from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue


@pytest.fixture
def client():
    with app.test_client() as client:
        yield client

    
""" 

@revenue_end_point.route('/budget/revenue', methods=['GET'])
def get_revenue():
    return {} 


@revenue_end_point.route('/budget/revenue', methods=['POST'])
def save_revenue():
    return {} 

@revenue_end_point.route('/budget/revenue/<id>', methods=['PUT'])
def update_revenue():
    return {} 


@revenue_end_point.route('/budget/revenue/<id>', methods=['DELETE'])
def delete_revenue():
    return {} 

    public void addANewBudgetRevenue() throws Exception {
    BudgetRevenue budgetRevenue = new BudgetRevenue(null, "USER", dateFor("10/10/2018"), Money.ONE, "A_NOTE");

    given(userRepository.currentLoggedUserName()).willReturn(new UserName("USER"));

    BudgetRevenueRepresentation budgetRevenueRepresentation = new BudgetRevenueRepresentation(null, "10/10/2018", "1.00", "A_NOTE");
    mockMvc.perform(post("/budget/revenue")
            .contentType(MediaType.APPLICATION_JSON)
            .with(csrf())
            .content(objectMapper.writeValueAsBytes(budgetRevenueRepresentation)))
            .andExpect(status().isCreated());

    verify(budgetRevenueRepository).save(budgetRevenue);
    verify(userRepository).currentLoggedUserName();
}
"""

@pytest.mark.skip(reason="to be fixed")
def test_add_new_revenue(client, mocker: MockerFixture):
    mocked_user_name_resolver_response = UserName("A_USER_NAME")

    mocked_save_revenue_service = mocker.Mock()
    mocked_save_revenue_service.save_revenue
    mocked_repository = mocker.Mock()

    mocked_user_name_resolver = mocker.Mock()
    mocked_user_name_resolver.get_user_name.return_value = (
        mocked_user_name_resolver_response
    )

    response = client.post(
        "/budget/revenue",
        json={"date": "10/10/2018", "amount": "1.00", "note": "A_NOTE"},
    )

    assert response.status_code == 201

    mocked_repository.get_user_name.assert_called()
    mocked_user_name_resolver.assert_called_once_with(
        Revenue(
            id=None,
            user_name=UserName("A_USER_NAME"),
            registration_date=Date.iso_date_for("2018-10-10"),
            amount=Money.money_for("1.00"),
            note="A_NOTE",
        )
    )
