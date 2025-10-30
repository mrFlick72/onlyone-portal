import pytest

from app.revenue.api.converter import QueryParamRepresentationConverter
from app.time.domain.year import Year


def test_form_raw_query_param_to_year():
    uut = QueryParamRepresentationConverter()

    query_param = "year=2025;other_param=value"

    representation = uut.from_query_param_to_year(query_param)

    assert isinstance(representation.year, Year)
    assert representation.year.content == 2025
