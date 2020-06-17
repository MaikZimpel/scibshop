import React, {Component} from "react";
import './scib-content.scss'
import {connect} from "react-redux";

class ScibContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            storeItems: []
        }
    }

    render() {
        return (
            <div className={"content"}>

            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        storeItems: state.storeItems
    }
}

export default connect(mapStateToProps) (ScibContent);

