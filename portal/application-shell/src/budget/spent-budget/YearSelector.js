import React from "react";
import FormInputTextField from "../../components/form/FormInputTextField";

export default ({year, handler}) => {
    return <FormInputTextField id="yearSelector"
                               value={year}
                               label=""
                               type="number"
                               handler={handler}/>
}