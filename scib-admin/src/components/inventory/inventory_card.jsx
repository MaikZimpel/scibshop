import React, {useContext, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@material-ui/core/CardActions';
import IconButton from '@material-ui/core/IconButton';
import TextField from '@material-ui/core/TextField';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import SaveIcon from '@material-ui/icons/Save';
import CancelIcon from '@material-ui/icons/Cancel'
import Collapse from '@material-ui/core/Collapse';
import clsx from "clsx";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import {InventoryContext} from "./inventoryContext";
import * as api from './inventory_api';
import {ItemVariant} from "./item_variant";

const useStyles = makeStyles((theme) => ({
    root: {
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
    expand: {
        transform: 'rotate(0deg)',
        marginLeft: 'auto',
        transition: theme.transitions.create('transform', {
            duration: theme.transitions.duration.shortest,
        }),
    },
    expandOpen: {
        transform: 'rotate(180deg)',
    },
}));

export const InventoryCard = (props) => {
    /*
    type Item struct {
	Id          string      `json:"id" bson:"_id, omitempty"`
	Upc         string      `json:"upc" bson:"upc, omitempty"`
	Name        string      `json:"name" bson:"name, omitempty"`
	Description string      `json:"description" bson:"description, omitempty"`
	Categories  []string    `json:"categories" bson:"categories, omitempty"`
	Brand       string      `json:"brand" bson:"brand, omitempty"`
	Price       float32     `json:"price" bson:"price, omitempty"`
	Images      []string    `json:"images" bson:"images, omitempty"`
	Supplier    string      `json:"supplier" bson:"supplier, omitempty"`
	Variants []ItemVariant `json:"variants" bson:"variants, omitempty"`
	*/

    const {items, actions} = useContext(InventoryContext);
    const [item, setItem] = useState(items.find(i => i.id === props.itemId));

    const [expanded, setExpanded] = useState(false)
    const [saveBtnEnabled, setSaveBtnEnabled] = useState(false);

    function handleUpdate (e)  {
        const name = e.target.name;
        const value = function() {
            switch (name) {
                case 'categories': return e.target.value ? e.target.value.split(" ")  : [];
                case 'price': return Number(e.target.value);
                default: return e.target.value;
            }
        }() ;
        setItem(item => ({...item, [name]: value}));
        actions.updateItem(item.id, name, value);
    }

    const itemCopy = {...item};

    const handleExpandClick = () => {
        setExpanded(!expanded);
    };

    const enableSaveBtn = () => {
        setSaveBtnEnabled(true);
    }

    const disableSaveBtn = () => {
        setSaveBtnEnabled(false);
    }

    const classes = useStyles();

    function addImage(axiosResponse) {
        actions.addItemImage(item.id, axiosResponse.data);
    }

    const onFileChange = event => {
        const formData = new FormData()
        formData.append("content_type", "image/jpeg")
        const file = event.target.files[0]
        formData.append('originalFile', file)
        api.upload(item.id, formData).then(addImage)
    }

    function deletePicture(imageName) {
        api.deleteImageFile(item.id, imageName).then(() => actions.removeItemImage(item.id, imageName)).catch(console.log)
    }

    const saveItem = () => {
        api.saveItem(items.find(i => i.id === props.itemId)).then(result => {
            switch (result) {
                case 201:
                case 204:
                    break;
                default:
                    console.log(result);
                    break;
            }
        })
    };

    const resetItem = () => {
        setItem({...item, name: itemCopy.name});
        setItem({...item, brand: itemCopy.brand});
        setItem({...item, supplier: itemCopy.supplier});
        setItem({...item, categories: itemCopy.categories});
        setItem({...item, description: itemCopy.description});
        setItem({...item, price: itemCopy.price});
        setItem({...item, upc: itemCopy.upc});
        enableSaveBtn()
    }

    const addItemVariant = () => {
        actions.addVariant(item.id, {
            sku: "",
            color: "",
            image: "",
            size: "",
            cnt: 0,
            stockable: true,
            available: true
        });
    }

    return (
        <Card className={classes.root} variant={"outlined"} elevation={2}>
            <CardHeader
                avatar={
                    item.images && item.images.length > 0 ?
                        <img alt={""} className={"avtr-pic"}
                             src={"http://localhost:8082/inventory/" + item.id + "/images/" + item.images[0]}/>
                        :
                        <ShoppingCartIcon fontSize={"large"}/>

                }
                title={item.name}
                subheader={
                    <div>
                        <div>
                            EUR: {item.price}
                        </div>
                    </div>
                }
            >
            </CardHeader>
            <IconButton
                className={clsx(classes.expand, {[classes.expandOpen]: expanded})}
                onClick={handleExpandClick}
                aria-expanded={expanded}
                aria-label="show more"
            >
                <ExpandMoreIcon/>
            </IconButton>

            <Collapse in={expanded} timeout={"auto"} unmountOnExit>
                <CardContent>
                    <div className={"horizontal-box"}>
                        {
                            item.images ?
                                item.images.map((val, index) => {
                                    return (
                                        <div key={index} className={"container"}>
                                            <img className={"image"} alt={""}
                                                 src={"http://localhost:8082/inventory/" + item.id + "/images/" + val}/>
                                            <div className={"middle"}>
                                                <Fab
                                                    color="primary"
                                                    size="small"
                                                    component="div"
                                                    aria-label="remove"
                                                    variant="extended"
                                                >
                                                    <DeleteIcon onClick={event => deletePicture(val, index)}/>
                                                </Fab>
                                            </div>
                                        </div>
                                    );
                                }) : <div/>
                        }
                        <div className={"vertical-center"}>
                            <form>
                                <label htmlFor={"upload-picture"}>
                                    <input style={{display: 'none'}} id={"upload-picture"} name={"upload-picture"}
                                           type={"file"} onChange={onFileChange} aria-label={"Add picture"}/>
                                    <IconButton size="small" component="div" aria-label="add" variant="extended">
                                        <AddIcon/>
                                    </IconButton>
                                </label>

                            </form>
                        </div>
                    </div>
                </CardContent>
                <CardContent>
                    <div className={classes.root} style={{width: '75vh'}}>
                        <Grid container spacing={1}>
                            <Grid item xs={12}>
                                <TextField required name={"description"} value={item.description} label={"Description"}
                                           onChange={handleUpdate}
                                           margin="dense" InputLabelProps={{shrink: true,}} multiline
                                           rows={4}
                                           variant={"outlined"}
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField disabled={item.isPersistent} required name={"name"} value={item.name}
                                           label={"Name"}
                                           variant={"outlined"}
                                           margin="dense"
                                           onChange={handleUpdate}
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField disabled={item.isPersistent} name={"brand"} required value={item.brand}
                                           label={"Brand"}
                                           variant={"outlined"}
                                           onChange={handleUpdate} margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField name={"price"} value={item.price} label={"Price"}
                                           variant={"outlined"}
                                           onChange={handleUpdate}
                                           margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <Paper variant={"outlined"} style={{margin: "5px", padding: "5px"}}>
                                    {
                                        item.variants ?
                                            item.variants.map((variant, index) => {
                                                return (
                                                    <ItemVariant index={index} itemId={item.id}/>)
                                            })
                                        : <div/>
                                    }
                                    <IconButton size="small" component="div" aria-label="add" variant="extended">
                                        <AddIcon onClick={addItemVariant}/>
                                    </IconButton>
                                </Paper>
                            </Grid>
                        </Grid>
                    </div>
                </CardContent>
                <CardFooter>
                    <IconButton id={"cancelBtn"} aria-label={"cancel"} onClick={resetItem} disabled={!saveBtnEnabled}>
                        <CancelIcon/>
                    </IconButton>
                    <IconButton id={"saveBtn"} aria-label="save" onClick={saveItem}>
                        <SaveIcon/>
                    </IconButton>
                </CardFooter>
            </Collapse>
        </Card>
    );
}
