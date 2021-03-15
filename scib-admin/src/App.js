import React, {useState} from 'react';
import './App.css';

import Appcontainer from "./components/appcontainer/appcontainer";
import Login from "./components/signin/Login"

export default function App() {

    const [token, setToken] = useState();

    if (!token) {
        return <Login setToken={setToken} />
    }

    return (
        <Appcontainer/>
    )
}

