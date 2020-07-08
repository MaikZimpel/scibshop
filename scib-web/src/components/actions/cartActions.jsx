import {ADD_TO_CART, ADD_QUANTITY, ADD_SHIPPING, REMOVE_ITEM, SUB_QUANTITY} from "./action-types/cart-actions";

export const add = (sku) => {
    return {
        type: ADD_TO_CART,
        sku
    }
}

export const inc = (sku) => {
    return {
        type : ADD_QUANTITY,
        sku
    }
}

export const dec = (sku) => {
    return {
        type: SUB_QUANTITY,
        sku
    }
}

export const rmv = (sku) => {
    return {
        type: REMOVE_ITEM,
        sku
    }
}