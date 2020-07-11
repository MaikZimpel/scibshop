import React, {useContext, useEffect} from "react";
import './scib-content.scss'
import {makeStyles} from "@material-ui/core/styles";
import {ItemCard} from "./itemcard";
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import {CartContext} from "../cart-context/cartContext";
import * as api from '../cart-context/cartApi';
import {CartDialog} from "./cartDlg";

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow:1,
        control: {
            padding: theme.spacing(2),
        },
        '& .MuiTextField-root': {
            margin: theme.spacing(1),
            width: '25ch',
        },
        '& .MuiCardHeader-root': {
            'display': 'flex',
            'flex-direction': 'row',
            'align-items': 'flex-start',
        },
        '& .MuiCardHeader-subheader': {}
    },
}));

export const Shopfront = () => {

    useEffect(() => {
        api.loadInventory().then(data => actions.loadInventory(data));
    }, []);

    const classes = useStyles();
    const {items, actions } = useContext(CartContext);

    const itemList = items.map((item, index) => {
        return (
            <Grid key={index}>
                <Paper className={classes.paper}>
                    <ItemCard itemId={item.id}/>
                </Paper>
            </Grid>

        )
    })

    return (
        <Grid container className={"item-display"}>
            <Grid item xs={12}>
                <Grid container justify={"center"} spacing={1}>
                    {itemList}
                </Grid>
            </Grid>
            <CartDialog/>
        </Grid>
    )
}

