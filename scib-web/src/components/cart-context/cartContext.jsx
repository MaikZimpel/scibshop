import React, {createContext, useReducer} from 'react'
import {cartReducer} from "./cartReducer";
import * as actions from './cartReducer';



/*
    state looks like this:
    state: {
        items: [{id, brand, name, description, price, images:[], variants: [{sku, qty, color, image}]}],
        cart: {
            items: [{itemId, sku, price, qty}],
            itemQty,
            total,
            show
        },
        actions: {}
    }
     */
const initialState = {items: [], cart: {items: [], itemQty: 0, total: 0, show: false}, actions: {}};
const localCart = JSON.parse(localStorage.getItem("cart"));

export const CartContext = createContext(initialState);

export const CartProvider = ({children}) => {

    const [state, dispatch] = useReducer(cartReducer, initialState);

    const loadInventory = items => {
        dispatch({
            type: actions.LOAD_INVENTORY,
            items
        })
    }

    const addToCart = (itemId, sku, price, qty) => {
        dispatch({
            type: actions.ADD_TO_CART,
            itemId,
            sku,
            price,
            qty
        })
    }

    const rmvFromCart = sku => {
        dispatch({
            type: actions.REMOVE_ITEM,
            sku
        })
    }

    const toggleCartDialog = () => {
        dispatch({
            type: actions.TOGGLE_CART_DIALOG
        })
    }

    return (
        <CartContext.Provider value={{
            items: state.items,
            cart: state.cart,
            actions: {
                loadInventory,
                addToCart,
                rmvFromCart,
                toggleCartDialog
            }
        }}>
            {children}
        </CartContext.Provider>
    )
}