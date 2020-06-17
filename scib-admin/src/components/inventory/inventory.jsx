import React, {Component} from "react";
import InventoryCard from "./inventory_card";
import './inventory.scss'

class Inventory extends Component {


    state = {
        inventoryItems: []
    }

    componentDidMount() {
        fetch("http://localhost:8082/inventory?stockableOnly=false")
            .then(res => res.json())
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