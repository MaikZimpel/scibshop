import React from 'react';
import './App.scss';
import ScibHeader from "./components/scib_header/scib-header";
import ScibFooter from "./components/scib-footer/scib-footer";
import {Shopfront} from "./components/scib-content/shopfront";
import {CartProvider} from "./components/cart-context/cartContext";

function App() {
    return (
        <div className="App">
            <CartProvider>
                <ScibHeader/>
                <Shopfront/>
                <ScibFooter/>
            </CartProvider>
        </div>
    );
}

export default App;
