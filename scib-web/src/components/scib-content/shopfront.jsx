import React, {useState, useEffect} from "react";
import './scib-content.scss'
import { connect } from 'react-redux'
import {add} from "../actions/cartActions";
import axios from 'axios'
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@material-ui/core/CardActions';
import {makeStyles} from "@material-ui/core/styles";
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart'
import {ItemCard} from "./itemcard";
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';

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

    const classes = useStyles();
    const [items, setItems] = useState([]);
    const [cart, setCart] = useState();

    useEffect(() => {
        let req = {
            url: "http://localhost:8082/inventory/?stockableOnly=false",
            method: 'GET',
            mode: 'no-cors'
        };
        axios(req)
            .then(res => setItems(res.data))
            .catch(console.log)
    }, []);


    const handleClick = (id) => {
        this.props.add(id)
    }

    const itemList = items.map((item, index) => {
        return (
            <Grid key={index}>
                <Paper className={classes.paper}>
                    <ItemCard item={item}/>
                </Paper>
            </Grid>

        )
    })

    return (
        <Grid container className={"item-display"} spacing={2}>
            <Grid item xs={12}>
                <Grid container justify={"center"} spacing={2}>
                    {itemList}
                </Grid>
            </Grid>
        </Grid>
    )
}

const mapStateToProps = (state) => {
    return {
        items: state.items
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        addToCart: (id) => {dispatch(add(id))}
    }
}

