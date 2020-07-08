export const InventoryReducer = (state, action) => {

    function removeItem(item) {
        return removeFromContext(item.id);
    }

    function removeFromContext(id) {
        const updatedState = {...state};
        updatedState.items = state.items.filter(item => item.id !== id);
        updatedState.selectedItem = updatedState.items ? updatedState.items[0] : null
        return updatedState;
    }

    function saveItem(item) {
        const itemArray = state.items;
        itemArray.splice(itemArray.findIndex(i => i.id === item.id), 1, item);
        return {...state, items: itemArray};
    }

    function addToContext(item) {
        if (state.items) {
            return {...state, items: [...state.items, item]};
        } else {
            return {...state, items: [item]};
        }
    }

    function updateItem(itemId, fieldName, fieldValue) {
        const item = findItem(itemId);
        return saveItem({...item, [fieldName]: fieldValue});
    }

    function addItemImage(itemId, image){
        const item = findItem(itemId);
        item.images = item.images ? [...item.images, image] : [image];
        return saveItem(item);
    }

    function removeItemImage(itemId, imageId) {
        const item = findItem(itemId);
        item.images = item.images.filter(img => img !== imageId);
        return saveItem(item);
    }

    function addVariant(itemId, variant) {
        const item = findItem(itemId);
        item.variants = item.variants ? [...item.variants, variant] : [variant];
        return saveItem(item);
    }

    function removeVariant(itemId, variantIndex) {
        const item = findItem(itemId);
        item.variants = [
            ...item.variants.slice(0, variantIndex),
            ...item.variants.slice(variantIndex + 1)
        ]
        return saveItem(item);
    }

    function updateVariant(itemId, variantIndex, name, value) {
        const item = findItem(itemId);
        const updatedItem = {...item, variants: [...item.variants]};
        updatedItem.variants[variantIndex] = {...updatedItem.variants[variantIndex], [name]: value};
        return saveItem(updatedItem);
    }

    function findItem(itemId) {
        return state.items.find(item => item.id === itemId);
    }

    switch (action.type) {
        case LOAD_INVENTORY:
            return {...state, items: action.items, selectedItem: action.items ? action.items[0] : null};
        case SELECT_ITEM:
            return {...state, selectedItem: state.items ? findItem(action.id) : null};
        case ADD_ITEM:
            return addToContext(action.item);
        case UPD_ITEM:
            return updateItem(action.itemId, action.fieldName, action.fieldValue);
        case RMV_ITEM:
            return removeItem(action.item);
        case SVE_ITEM:
            return saveItem(action.item);
        case ADD_ITEM_IMAGE:
            return addItemImage(action.itemId, action.imageId);
        case RMV_ITEM_IMAGE:
            return removeItemImage(action.itemId, action.imageId);
        case ADD_VARIANT:
            return addVariant(action.itemId, action.variant);
        case UPD_VARIANT:
            return updateVariant(action.itemId, action.variantIndex, action.name, action.value);
        case RMV_VARIANT:
            return removeVariant(action.itemId, action.variantIndex);
        default:
            return state;
    }
}

export const LOAD_INVENTORY = 'LD_INV';
export const SELECT_ITEM = 'SLCT_ITEM';
export const ADD_ITEM = 'ADD_ITEM';
export const UPD_ITEM ='UPD_ITEM';
export const RMV_ITEM = 'RMV_ITEM';
export const SVE_ITEM = 'SVE_ITEM';
export const ADD_ITEM_IMAGE = 'ADD_ITEM_IMG';
export const RMV_ITEM_IMAGE = 'RMV_ITEM_IMG';
export const ADD_VARIANT = 'ADD_VARIANT';
export const UPD_VARIANT = 'UPD_VARIANT';
export const RMV_VARIANT = 'RMV_VARIANT';
