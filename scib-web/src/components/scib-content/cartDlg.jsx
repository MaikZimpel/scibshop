import React, {useContext} from "react";
import {CartContext} from "../cart-context/cartContext";
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Button from '@material-ui/core/Button';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {Paper} from "@material-ui/core";
import {yellow} from "@material-ui/core/colors";
import { createMuiTheme, withStyles, makeStyles, ThemeProvider } from '@material-ui/core/styles';
import {CheckoutButton, CloseDlgButton} from "./buttons";

export const CartDialog = () => {


    const {items, cart, actions} = useContext(CartContext)

    const ccyFormat = (num) => {
        return `${num.toFixed(2)}`;
    }

    const cartContent = cart.items.map((cartItem, index) => {
        const inventoryItem = items.find(i => i.id === cartItem.itemId);
        const itemVariant = inventoryItem.variants.find(v => v.sku === cartItem.sku);
        return (
            <TableRow key={index}>
                <TableCell>{inventoryItem.brand} {inventoryItem.name} {itemVariant.color}</TableCell>
                <TableCell>{cartItem.qty}</TableCell>
                <TableCell align={"right"}>{ccyFormat(cartItem.qty * cartItem.price)}</TableCell>
            </TableRow>
        )
    })

    return(
        <Dialog open={cart.show} onClose={actions.toggleCartDialog}>
            <DialogTitle>Ihr Einkaufswagen</DialogTitle>
            <DialogContent>
                <TableContainer component={Paper}>
                    <Table size={"small"}>
                        <TableHead>
                            <TableRow>
                                <TableCell>Artikel</TableCell>
                                <TableCell>Menge</TableCell>
                                <TableCell align={"right"}>Preis</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {cartContent}
                            <TableRow key={"summary"}>
                                <TableCell align={"right"} colSpan={2}>Gesamt</TableCell>
                                <TableCell align={"right"}>{ccyFormat(cart.total)}</TableCell>
                            </TableRow>
                        </TableBody>
                    </Table>
                </TableContainer>
            </DialogContent>
            <DialogActions>
                <CheckoutButton onclick={""}>Zur Kasse</CheckoutButton>
                <CloseDlgButton onClick={actions.toggleCartDialog} color={"primary"}>Schliessen</CloseDlgButton>
            </DialogActions>
        </Dialog>
    )
}