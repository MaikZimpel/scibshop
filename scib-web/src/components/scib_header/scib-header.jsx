import React,  {Component} from 'react';
import logo from './scib-logo.jpg'
import './scib_header.scss'
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';

class ScibHeader extends Component {


    render() {
        return (
            <div className={"app-header"}>
                <div>
                    <img className={"logo"} src={logo} alt={"SCIB Logo"}/>
                </div>
                <div className={"punchline"}>
                    <span>AFRIKANISCHE SACHEN FÜR AFRIKANER UND ALLE ANDEREN</span>
                    <Tooltip title={"Show shopping cart"}>
                        <IconButton aria-label={"cart"}>
                            <ShoppingCartIcon/>
                        </IconButton>
                    </Tooltip>
                </div>
            </div>
        )
    }
}

export default ScibHeader;