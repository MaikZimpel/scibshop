import React, {Component} from 'react';
import logo from './scib-logo.jpg'
import './scib_header.scss'
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import EmailIcon from '@material-ui/icons/Email'
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faFacebook, faInstagram, faTwitter} from '@fortawesome/free-brands-svg-icons'

class ScibHeader extends Component {


    render() {
        return (
            <div>
                <div className={"clearfix"}>
                    <img className={"logo"} src={logo} alt={"SCIB Logo"}/>

                    <div className={"social-buttons"}>
                        <Tooltip title={"Show shopping cart"} className={"overflow"}>
                            <IconButton aria-label={"cart"}>
                                <ShoppingCartIcon/>
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Drop us a mail"} className={"overflow"}>
                            <IconButton aria-label={"mail"}>
                                <EmailIcon/>
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Follow us on Facebook"}>
                            <IconButton aria-label={"facebook"}>
                                <FontAwesomeIcon icon={faFacebook} size={"0.5x"}/>
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Follow us on Twitter"}>
                            <IconButton aria-label={"twitter"}>
                                <FontAwesomeIcon icon={faTwitter} size={"0.5x"}/>
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Follow us on Instagram"}>
                            <IconButton aria-label={"instagram"}>
                                <FontAwesomeIcon icon={faInstagram} size={"0.5x"}/>
                            </IconButton>
                        </Tooltip>
                    </div>
                </div>
            </div>
        )
    }
}

export default ScibHeader;