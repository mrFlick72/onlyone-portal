const selectUiAdapterFor = (searchTagRegistry: SearchTag[]) => searchTagRegistry.map(searchTag => {
    return { value: searchTag.key, label: searchTag.value }
})

export default selectUiAdapterFor;