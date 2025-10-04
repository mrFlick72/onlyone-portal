from flask import Blueprint

revenue_end_point = Blueprint('revenue_end_point',__name__)

@revenue_end_point.route('/incomes', methods=['GET'])
def get_incomes():
    return {} 

