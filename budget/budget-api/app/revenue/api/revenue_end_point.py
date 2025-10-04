from flask import Blueprint

revenue_end_point = Blueprint('revenue_end_point',__name__)

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

