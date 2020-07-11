export const cartReducer = (state, action) => {

    /*
    state looks like this:
    state: {
        items: [{id, name, brand, description, price, images:[], variants: [{sku, qty, color, image}]}],
        cart: {
            items: [{itemId, sku, price, qty}],
            itemQty,
            total
        },
        actions: {}
    }
     */

    const loadInventoryItems = items => {
        return {...state, items: items.map(item => {
                return {
                    id: item.id,
                    name: item.name,
                    brand: item.brand,
                    description: item.description,
                    price: item.price,
                    images: item.images,
                    variants: item.variants.map(variant => {
                        return {
                            sku: variant.sku,
                            qty: variant.cnt,
                            color: variant.color,
                            image: variant.image
                        }
                    })
                }
            })
        };
    }

    const addToCart = (itemId, sku, price, qty) => {
        const ceIndex = state.cart.items.findIndex(ce => ce.sku === sku);
        if (ceIndex === -1) {
            return {...state, cart: {...state.cart,
                    items: [...state.cart.items, {itemId, sku, price, qty}],
                    itemQty: state.cart.itemQty + qty,
                    total: state.cart.total + (price * qty)
            }};
        } else {
            const items = [...state.cart.items];
            items[ceIndex] = {...items[ceIndex], qty: state.cart.items[ceIndex].qty + qty};
            return {...state, cart: {...state.cart,
                    items: items,
                    itemQty: state.cart.itemQty + qty,
                    total: state.cart.total + (price * qty)
            }};
        }
    }

    const rmvFromCart = (sku, qty) => {
        const ceIndex = state.cart.findIndex(ce => ce.sku === sku);
        if (ceIndex === -1) {
            return state;
        } else {
            const newQty = state.cart.items[ceIndex].qty - qty;
            if (newQty > 0) {
                return {...state, cart: {...state.cart, ...state.cart.items[ceIndex], qty: newQty}};
            } else {
                return {...state, cart: [...state.cart.slice(0, ceIndex), ...state.cart.slice(ceIndex + 1)]};
            }
        }
    }

    switch (action.type) {
        case LOAD_INVENTORY: return loadInventoryItems(action.items);
        case ADD_TO_CART: return addToCart(action.itemId, action.sku, action.price, action.qty);
        case REMOVE_ITEM: return rmvFromCart(action.sku, action.qty);
        case TOGGLE_CART_DIALOG: return {...state, cart: {...state.cart, show: !state.cart.show}};
        default: return state
    }
}

export const LOAD_INVENTORY = 'LOAD_INVENTORY';
export const ADD_TO_CART = 'ADD_TO_CART';
export const REMOVE_ITEM = 'REMOVE_ITEM';
export const TOGGLE_CART_DIALOG = 'TOGGLE_CART_DIALOG';
export const ADD_SHIPPING = 'ADD_SHIPPING';