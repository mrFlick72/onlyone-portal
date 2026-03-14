import React from "react";
import FormSelect from "../../components/form/FormSelect";

type MonthsSelectorProps = {
    monthRegistry: { monthValue: string, monthLabel: string }[],
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
            value: item.monthValue,
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