import React, {Component, useContext} from 'react';
import logo from './scib-logo.jpg'
import './scib_header.scss'
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import EmailIcon from '@material-ui/icons/Email'
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faFacebook, faInstagram, faTwitter} from '@fortawesome/free-brands-svg-icons'
import {CartContext} from "../cart-context/cartContext";
import Badge from '@material-ui/core/Badge';
import { withStyles } from '@material-ui/core/styles';

const StyledBadge = withStyles((theme) => ({
    badge: {
        right: -3,
        top: 13,
        border: `2px solid ${theme.palette.background.paper}`,
        padding: '0 4px',
    },
}))(Badge);

const ScibHeader = () => {

    const {cart, actions} = useContext(CartContext);


        return (
            <div>
                <div className={"clearfix"}>
                    <img className={"logo"} src={logo} alt={"SCIB Logo"}/>
                    <div className={"social-buttons"}>
                        <Tooltip title={"Total: EUR " + cart.total} className={"overflow"}>
                            <IconButton aria-label={"cart"} onClick={actions.toggleCartDialog}>
                                <StyledBadge badgeContent={cart.itemQty} color="secondary">
                                    <ShoppingCartIcon/>
                                </StyledBadge>
                            </IconButton>
                        </Tooltip>
                        <Tooltip title={"Drop us a mail"} className={"overflow"}>
                            <IconButton aria-label={"mail"}>
                                <EmailIcon/>
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

export default ScibHeader;