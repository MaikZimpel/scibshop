import React, {useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@material-ui/core/CardActions';
import Avatar from '@material-ui/core/Avatar';
import IconButton from '@material-ui/core/IconButton';
import Checkbox from '@material-ui/core/Checkbox';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import axios from 'axios'
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import SaveIcon from '@material-ui/icons/Save';
import CancelIcon from '@material-ui/icons/Cancel'
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import clsx from "clsx";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import Grid from '@material-ui/core/Grid';

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
        }
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

export default function InventoryCard(props) {

    const [item, setItem] = useState(props.item)
    const [itemCopy, setItemCopy] = useState(props.item)
    const [expanded, setExpanded] = React.useState(false);
    const [saveBtnDisabled, setSaveBtnDisabled] = React.useState(true)

    const handleExpandClick = () => {
        setExpanded(!expanded);
    };

    const enableSaveBtn = () => {
        setSaveBtnDisabled(false)
    }

    /* Backend Data Structure
    Id          string             `json:"id"`
	Upc         string             `json:"upc"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Categories  []string           `json:"categories"`
	Brand       string             `json:"brand"`
	Size        string             `json:"size"`
	Color       string             `json:"color"`
	Price       int                `json:"price"`
	Images      []string           `json:"images"`
	Supplier    string             `json:"suppliers"`
	Sku         string             `json:"sku"`
	Cnt         int                `json:"cnt"`
	Stockable   bool               `json:"stockable"`
	Available   bool               `json:"available"`
     */

    const classes = useStyles();

    function addImage(axiosResponse) {
        const imagePath = axiosResponse.data
        console.log(imagePath)
        let imageArray = item.images
        if (imageArray == null) {
            imageArray = []
        }
        imageArray.push(imagePath)
        setItem({...item, images: imageArray})
    }

    function removeImage(ix) {
        let imageArray = item.images
        imageArray.splice(ix, 1)
        setItem({...item, images: imageArray})
    }

    const onFileChange = event => {
        // get signed url from backend
        const formData = new FormData()
        formData.append("content_type", "image/jpeg")
        const file = event.target.files[0]
        formData.append('originalFile', file)
        axios.post("http://localhost/inventory/" + item.id + "/images", formData)
            .then(addImage)
            .catch((ex) => console.log(ex))
    }

    function deletePicture(imageName, imgIndex) {
        axios.post("http://localhost/inventory/" + item.id + "/images/" + imageName)
            .then(removeImage(imgIndex))
            .catch((ex) => console.log(ex))
    }

    const saveItem = async () => {
        await axios.put("http://localhost/inventory/" + item.id, item)
            .then(() => setSaveBtnDisabled(true))
            .catch(x => console.log(x))
    };

    const resetItem = () => {
        setItem(itemCopy);
        setSaveBtnDisabled(true)
    }

    return (
        <Card className={classes.root} variant={"outlined"} elevation={2} width={'50vh'}>
            <CardHeader
                avatar={
                    item.images ?
                        <img alt={""} className={"avtr-pic"}
                             src={"http://localhost/inventory/" + item.id + "/images/" + item.images[0]}/>
                        :
                        <Avatar aria-label="inventory-item" className={"_avatar"}>
                            <ShoppingCartIcon/>
                        </Avatar>
                }
                title={item.name}
                subheader={
                    <div style={{overflow: "hidden", textOverflow: "ellipsis", width: '50vh'}}>
                        <Typography display={"inline"} variant="body2" color="textSecondary" component="span">
                            {item.description}
                        </Typography>
                    </div>
                }
            >
            </CardHeader>
            <IconButton
                className={clsx(classes.expand, {
                    [classes.expandOpen]: expanded,
                })}
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
                                                 src={"http://localhost/inventory/" + item.id + "/images/" + val}/>
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
                    <div className={classes.root} style={{width: '50vh'}}>
                        <Grid container spacing={1}>
                            <Grid item xs={12}>
                                <TextField required id={"description"} value={item.description} label={"Description"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, description: e.target.value});
                                           }}
                                           margin="dense" InputLabelProps={{shrink: true,}} multiline
                                           rows={4}
                                           variant={"outlined"}
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField required id={"name"} value={item.name} label={"Name"} variant={"outlined"}
                                           margin="dense"
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, name: e.target.value});
                                           }}
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField id={"upc"} value={item.upc} label={"Upc"} variant={"outlined"}
                                           onChange={e =>{
                                               enableSaveBtn();
                                               setItem({...item, upc: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField required id={"sku"} value={item.sku} label={"Sku"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, sku: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField id={"brand"} value={item.brand} label={"Brand"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn()
                                               setItem({...item, brand: e.target.value})
                                           }} margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField id={"size"} value={item.size} label={"Size"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, size: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField required id={"color"} value={item.color} label={"Color"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, color: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField id={"price"} type={"number"} value={item.price} label={"Price"}
                                           variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, price: Number(e.target.value)});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField id={"supplier"} value={item.suppliers} label={"Supplier"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, suppliers: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField type={"number"} required id={"cnt"} label={"Quantity"}
                                           variant={"outlined"} value={item.cnt}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, cnt: Number(e.target.value)});
                                           }}
                                           margin="dense"
                                           style={{width:'100%'}}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <FormControlLabel control={
                                    <Checkbox
                                        color="primary"
                                        checked={item.stockable}
                                        onChange={(e, value) => {
                                            enableSaveBtn();
                                            setItem({...item, stockable: value});
                                        }}
                                    />
                                } label={"Can be stocked"}/>
                            </Grid>
                            <Grid item xs={12}>
                                <FormControlLabel control={
                                    <Checkbox
                                        color="primary"
                                        checked={item.available}
                                        onChange={(e, value) => {
                                            enableSaveBtn();
                                            setItem({...item, available: value});
                                        }}
                                    />
                                } label={"Can be ordered"}/>
                            </Grid>
                        </Grid>
                    </div>
                </CardContent>
                <CardFooter>
                    <IconButton id={"cancelBtn"} aria-label={"cancel"} onClick={resetItem} disabled={saveBtnDisabled}>
                        <CancelIcon />
                    </IconButton>
                    <IconButton id={"saveBtn"} aria-label="save" onClick={saveItem} disabled={saveBtnDisabled}>
                        <SaveIcon />
                    </IconButton>
                </CardFooter>
            </Collapse>
        </Card>
    );
}
