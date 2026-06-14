from typing import List, Optional

from pydantic import BaseModel, ConfigDict, Field, model_validator

MAX_YEAR_RANGE_SIZE = 20


class TotalByTagAnalysisRequestRepresentation(BaseModel):
    year: int
    month: Optional[int] = Field(default=None, ge=1, le=12)
    tags: List[str] = Field(default_factory=list)


class TotalByYearAnalysisRequestRepresentation(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    from_year: int = Field(alias="fromYear")
    to_year: int = Field(alias="toYear")
    tag: Optional[str] = None

    @model_validator(mode="after")
    def validate_year_range(self) -> "TotalByYearAnalysisRequestRepresentation":
        if self.from_year > self.to_year:
            raise ValueError("fromYear must not be greater than toYear")
        if self.to_year - self.from_year + 1 > MAX_YEAR_RANGE_SIZE:
            raise ValueError(f"year range must not exceed {MAX_YEAR_RANGE_SIZE} years")
        return self


class TagTotalRepresentation(BaseModel):
    tag: str
    total: str


class YearTotalRepresentation(BaseModel):
    year: int
    total: str
