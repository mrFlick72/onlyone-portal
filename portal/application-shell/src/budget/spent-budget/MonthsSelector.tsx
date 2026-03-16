import React from "react";
import FormSelect from "../../components/form/FormSelect";
import { Month } from "../../time/months";

type MonthsSelectorProps = {
    monthRegistry: Month[],
    month: string,
    handler: (event: any) => void
}


const MonthsSelector: React.FC<MonthsSelectorProps> = ({ monthRegistry, month, handler }) => {

    let valueLabel: string = "";
    let options = monthRegistry.map(item => {
        if (Number(item.monthValue) === Number(month)) {
            valueLabel = item.monthLabel
        }
        return {
            value: String(item.monthValue),
            label: item.monthLabel
        }
    });

    return <FormSelect id="monthsSelector"
        value={{ value: month, label: valueLabel }}
        label=""
        multi={false}
        options={options}
        onChangeHandler={handler} />
}

export default MonthsSelector