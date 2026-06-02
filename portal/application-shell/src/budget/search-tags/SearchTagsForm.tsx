import { Save } from "@mui/icons-material";
import FormInputTextField from "../../components/form/FormInputTextField";
import FormButton from "../../components/form/FormButton";

interface SearchTagFormProps {
    searchTag: SearchTag;
    handler: any;
    messages: { valueLabel: string; saveLabel: string };
}
const SearchTagForm: React.FC<SearchTagFormProps> = ({
    searchTag,
    handler,
    messages,
}) => {
    return (
        <div>
            <FormInputTextField
                handler={handler.valueHandler}
                autoFocus={true}
                value={searchTag.value}
                id="searchTagValue"
                label={messages.valueLabel}
            />
            <FormButton
                type="button"
                labelPrefix={<Save />}
                label={messages.saveLabel}
                onClickHandler={handler.submitHandler.bind(
                    this,
                    searchTag.key,
                    searchTag.value
                )}
            />
        </div>
    );
};

export default SearchTagForm;