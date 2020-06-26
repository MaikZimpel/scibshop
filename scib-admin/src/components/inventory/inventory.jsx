import React, {Component} from "react";
import InventoryCard from "./inventory_card";
import './inventory.scss'
import axios from 'axios'
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import Grid from '@material-ui/core/Grid';
import Button from '@material-ui/core/Button'

import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import PropTypes from 'prop-types';

class Inventory extends Component {

    state = {
        inventoryItems: [],
        selectedItem: 0,
        confirmationDlgOpen: false
    }

    handleDeleteConfirmationDlgClose = (yesOrNo) => {
        this.setState({
            confirmationDlgOpen: false
        })
        if (yesOrNo === 'yes') {
            this.deleteSelectedItem()
        }
    };

    deleteSelectedItem() {
        const itemId = this.state.inventoryItems[this.state.selectedItem].id
        const sIndex = this.state.selectedItem
        console.log('delete item: ' + itemId)
        axios.delete("http://localhost:8082/inventory/"+itemId)
            .then(response => {
                if(response.status === 204) {
                    /*
                        let imageArray = item.images
                        imageArray.splice(ix, 1)
                        setItem({...item, images: imageArray})
                     */
                    let itemArray = this.state.inventoryItems
                    itemArray.splice(sIndex, 1)
                    this.setState({
                        inventoryItems: itemArray,
                        selectedItem: 0
                    })
                } else {
                    console.log(response.statusText)
                }
            })
            .catch(ex => console.log(ex))
    }

    componentDidMount() {
        let req = {
            url: "http://localhost:8082/inventory/?stockableOnly=false",
            method: 'GET',
            mode: 'no-cors'
        };
        axios(req)
            .then(res => res.data)
            .then((data) => {
                if (data) {this.setState({inventoryItems: data})}
            })
            .catch(console.log)
    }

    newItem = () => {
        const itemArr = this.state.inventoryItems
        const l = itemArr.push({
            id: "",
            upc: "",
            name: "",
            description: "",
            categories: [],
            brand: "",
            images: [],
            supplier: "",
            sku: "sku",
            cnt: 0,
            stockable: true,
            available: true,
            price: {
                amount: 0,
                currency: "EUR"
            }
        })
        this.setState({
            inventoryItems: itemArr,
            selectedItem: l-1
        })
    }

    isSelected = (index) => {
        return this.state.selectedItem === index;
    }

    selectItem = (index) => {
        this.setState({
            selectedItem: index
        })
    }

    onClickDelete = () => {
        this.setState({
            confirmationDlgOpen: true
        })
    }

    render() {
        return (
            <div className={"inv-main"}>
                <div style={{margin: '10px'}}>
                    <Fab aria-label={"add"} style={{margin: '10px'}}>
                        <AddIcon onClick={this.newItem}/>
                    </Fab>
                    <Fab aria-label={"delete"}>
                        <DeleteIcon onClick={this.onClickDelete}/>
                    </Fab>
                </div>
                <ConfirmationDialog
                    onClose={this.handleDeleteConfirmationDlgClose}
                    keepMounted
                    open={this.state.confirmationDlgOpen}
                />
                <Grid container spacing={2}>
                {
                    this.state.inventoryItems.map((val, index) => {
                        return (
                            <Grid item key={index.toString()}>
                                <div onClick={this.selectItem.bind(this, index)} className={`${this.isSelected(index) ? "selected" : ""}`}><InventoryCard item={val}/></div>
                            </Grid>
                        );
                    })

                }
                </Grid>
            </div>
        )
    }
}

function ConfirmationDialog(props) {

    const { onClose, open } = props;


    const handleCancel = () => {
        onClose('no');
    };

    const handleOk = () => {
        onClose('yes');
    };

    return(
        <Dialog
            disableBackdropClick
            disableEscapeKeyDown
            maxWidth="xs"
            open={open}
            >
            <DialogTitle id={"confirm-delete-title"}>Confirm delete item</DialogTitle>
            <DialogContent dividers>
                <span>The item [placeholder] is going to disappear from your shop and your admin. Open orders are still going to be processed. Are you sure?</span>
            </DialogContent>
            <DialogActions>
                <Button autoFocus onClick={handleCancel} color="primary">
                    Yikes! Maybe Not!
                </Button>
                <Button onClick={handleOk} color="primary">
                    Yes, Go ahead!
                </Button>
            </DialogActions>
        </Dialog>

    );
}

ConfirmationDialog.propTypes = {
    onClose: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired
};

export default Inventory