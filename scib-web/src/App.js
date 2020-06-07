import React from 'react';
import './App.scss';
import ScibHeader from "./components/scib_header/scib-header";
import ScibFooter from "./components/scib-footer/scib-footer";
import ScibContent from "./components/scib-content/scib-content";

function App() {
  return (
    <div className="App">
      <ScibHeader/>
      <ScibContent/>
      <ScibFooter/>
    </div>
  );
}

export default App;
