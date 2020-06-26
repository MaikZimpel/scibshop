import React from 'react';
import './App.scss';
import ScibHeader from "./components/scib_header/scib-header";
import ScibFooter from "./components/scib-footer/scib-footer";
import {Shopfront} from "./components/scib-content/shopfront";

function App() {
  return (
    <div className="App">
      <ScibHeader/>
      <Shopfront/>
      <ScibFooter/>
    </div>
  );
}

export default App;
