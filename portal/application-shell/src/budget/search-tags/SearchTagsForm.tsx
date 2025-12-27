import { Save } from "@mui/icons-material";
import FormInputTextField from "../../components/form/FormInputTextField";
import FormButton from "../../components/form/FormButton";

interface SearchTagFormProps {
    searchTag: SearchTag;
    handler: any;
}
const SearchTagForm: React.FC<SearchTagFormProps> = ({
    searchTag,
    handler,
}) => {
    return (
        <div>
            <FormInputTextField
                handler={handler.valueHandler}
                autoFocus={true}
                value={searchTag.value}
                id="searchTagValue"
                label="Search Tag Value"
            />
            <FormButton
                type="button"
                labelPrefix={<Save />}
                label="Save"
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