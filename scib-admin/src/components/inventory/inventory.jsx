import React, {Component} from "react";
import InventoryCard from "./inventory_card";
import './inventory.scss'
import axios from 'axios'

class Inventory extends Component {


    state = {
        inventoryItems: []
    }

    componentDidMount() {
        let req = {
            url: "http://localhost/inventory/?stockableOnly=false",
            method: 'GET',
            mode: 'no-cors'
        };
        axios(req)
            .then(res => res.data)
            .then((data) => {
                this.setState({inventoryItems: data})
            })
            .catch(console.log)
    }

    render() {
        return (
            <div className={"inv-main"}>
                {
                    this.state.inventoryItems.map((val, index) => {
                        return (
                            <div key={index.toString()} className={"inv-card-container"}>
                                <InventoryCard item={val}/>
                            </div>
                        );
                    })

                }
            </div>
        )
    }
}

export default Inventory