import React, {useContext, useState, useEffect} from "react";
import {InventoryCard} from "./inventory_card";
import './inventory.scss'
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import Grid from '@material-ui/core/Grid';
import Button from '@material-ui/core/Button'
import {v4 as uuid} from 'uuid';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import PropTypes from 'prop-types';
import {InventoryContext} from "./inventoryContext";
import * as api from './inventory_api';

export const Inventory = () => {

    const [confirmationDialogOpen, setConfirmationDialogOpen] = useState(false);
    const [selectedItemIndex, setSelectedItemIndex] = useState(0);
    const { items, selectedItem, actions } = useContext(InventoryContext);

    useEffect(() => {
        api.loadInventory().then(data => actions.loadInventory(data));
    }, [])

    const handleDeleteConfirmationDlgClose = (yesOrNo) => {
        setConfirmationDialogOpen(false);
        if (yesOrNo === 'yes') {
            deleteSelectedItem();
        }
    };

    const deleteSelectedItem = () => {
        api.removeItem(selectedItem).then(delOk => {
            console.log({delOk})
            if(delOk) {
                actions.removeItem(selectedItem);
            } else {
                console.log("item could not be removed");
            }
        })
    }

    const isSelected = (id) => {
        return selectedItemIndex === id;
    }

    const selectItem = (id, index) => {
        actions.selectItem(id);
        setSelectedItemIndex(index);
    }

    const onClickDelete = () => {
        setConfirmationDialogOpen(true)
    }

    const newItem = () => {
        const item = {
            id: uuid(),
            upc: "",
            name: "",
            description: "",
            categories: [],
            brand: "",
            price: 0.0,
            images: [],
            supplier: "",
            variants: [{
                sku: "",
                color: "",
                image: "",
                size: "",
                cnt: 1,
                stockable: true,
                available: true
            }],
            isPersistent: false
        }
        actions.addItem(item)
    }


    return (
        <div className={"inv-main"}>
            <div style={{margin: '10px'}}>
                <Fab aria-label={"add"} style={{margin: '10px'}}>
                    <AddIcon onClick={newItem}/>
                </Fab>
                <Fab aria-label={"delete"}>
                    <DeleteIcon onClick={onClickDelete}/>
                </Fab>
            </div>
            <ConfirmationDialog
                onClose={handleDeleteConfirmationDlgClose}
                keepMounted
                open={confirmationDialogOpen}
            />
            <Grid container spacing={2}>
                {
                    items ?
                    items.map((item, index) => {
                        return (
                            <Grid item key={item.id}>
                                <div onClick={() => selectItem(item.id, index)} className={`${isSelected(index) ? "selected" : ""}`}>
                                    <InventoryCard itemId={item.id}/>
                                </div>
                            </Grid>
                        );
                    }): <div/>
                }
            </Grid>
        </div>
    )

}

function ConfirmationDialog(props) {

    const {onClose, open} = props;


    const handleCancel = () => {
        onClose('no');
    };

    const handleOk = () => {
        onClose('yes');
    };

    return (
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