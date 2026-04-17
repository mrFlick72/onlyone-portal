import React from "react";
import FormInputTextField from "../../components/form/FormInputTextField";

type YearSelectorProps = {
    year: string,
    handler: (event: React.ChangeEvent<HTMLInputElement>) => void
}

const YearSelector: React.FC<YearSelectorProps> = ({ year, handler }) => {
    return <FormInputTextField id="yearSelector"
        value={year}
        label=""
        type="number"
        handler={handler} />
}

export default YearSelector