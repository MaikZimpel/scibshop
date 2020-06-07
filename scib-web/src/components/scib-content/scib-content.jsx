import React, {Component} from "react";
import './scib-content.scss'

class ScibContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            storeItems: []
        }
    }

    componentDidMount() {
        fetch('http://localhost:8082/store')
            .then(res => res.json())
            .then((data) => {
                this.setState({storeItems: data})
            })
            .catch(console.log)
    }

    render() {
        return (
            <div className={"content"}>

            </div>
        )
    }
}

export default ScibContent;

