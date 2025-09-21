const selectUiAdapterFor = (searchTagRegistry) => searchTagRegistry.map(searchTag => {
    return {value: searchTag.key, label: searchTag.value}
})

export default selectUiAdapterFor;