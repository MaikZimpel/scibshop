import {ADD_TO_CART, ADD_QUANTITY, ADD_SHIPPING, REMOVE_ITEM, SUB_QUANTITY} from "./action-types/cart-actions";

export const add = (id) => {
    return {
        type: ADD_TO_CART,
        id
    }
}

export const inc = (id) => {
    return {
        type : ADD_QUANTITY,
        id
    }
}

export const dec = (id) => {
    return {
        type: SUB_QUANTITY,
        id
    }
}

export const rmv = (id) => {
    return {
        type: REMOVE_ITEM,
        id
    }
}