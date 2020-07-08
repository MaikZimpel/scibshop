import React, {createContext, useReducer} from 'react'
import * as InvRed from './inventory_reducer'

const initialState = {items: [], selectedItem: {}, actions: {}};

export const InventoryContext = createContext(initialState);

export const InventoryProvider = ({children}) => {

    const [state, dispatch] = useReducer(InvRed.InventoryReducer, initialState);

    const loadInventory = (data) => {
        dispatch({
            type: InvRed.LOAD_INVENTORY,
            items: data
        });
    }

    const selectItem = id => {
        dispatch({
            type: InvRed.SELECT_ITEM,
            id
        })
    }

    const addItem = item => {
        dispatch({
            type: InvRed.ADD_ITEM,
            item
        })
    }

    const updateItem = (itemId, fieldName, fieldValue) => {
        dispatch({
            type: InvRed.UPD_ITEM,
            itemId,
            fieldName,
            fieldValue
        })
    }

    const removeItem = item => {
        dispatch({
            type: InvRed.RMV_ITEM,
            item
        })
    }

    const saveItem = item => {
        dispatch({
            type: InvRed.SVE_ITEM,
            item
        })
    }

    const addItemImage = (itemId, imageId) => {
        dispatch({
            type: InvRed.ADD_ITEM_IMAGE,
            itemId,
            imageId
        })
    }

    const removeItemImage = (itemId, imageId) => {
        dispatch({
            type: InvRed.RMV_ITEM_IMAGE,
            itemId,
            imageId
        })
    }

    const addVariant = (itemId, variant) => {
        dispatch({
            type: InvRed.ADD_VARIANT,
            itemId,
            variant
        })
    }

    const updateVariant = (itemId, variantIndex, name, value) => {
        dispatch({
            type: InvRed.UPD_VARIANT,
            itemId,
            variantIndex,
            name,
            value
        })
    }

    const removeVariant = (itemId, variantIndex) => {
        dispatch({
            type: InvRed.RMV_VARIANT,
            itemId,
            variantIndex
        })
    }

    return (
        <InventoryContext.Provider value={{
            items: state.items,
            selectedItem: state.selectedItem,
            actions: {
                loadInventory,
                selectItem,
                addItem,
                updateItem,
                removeItem,
                saveItem,
                addItemImage,
                removeItemImage,
                addVariant,
                updateVariant,
                removeVariant
            }
        }}>
            {children}
        </InventoryContext.Provider>
    )
}