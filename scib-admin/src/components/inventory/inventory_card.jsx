import React, {useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@material-ui/core/CardActions';
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
import Collapse from '@material-ui/core/Collapse';
import clsx from "clsx";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import Grid from '@material-ui/core/Grid';
import EuroSymbolIcon from '@material-ui/icons/EuroSymbol';
import Paper from '@material-ui/core/Paper';

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

    type Item struct {
	Id          string      `json:"id"`
	Upc         string      `json:"upc"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Categories  []string    `json:"categories"`
	Brand       string      `json:"brand"`
	Sizes       []ItemSize  `json:"sizes"`
	Colors      []ItemColor `json:"colors"`
	Price       money.Money `json:"price"`
	Images      []string    `json:"images"`
	Supplier    string      `json:"supplier"`
	Sku         string      `json:"sku"`
	Cnt         int         `json:"cnt"`
	Stockable   bool        `json:"stockable"`
	Available   bool        `json:"available"`
}

type ItemColor struct {
	Sku       string `json:"sku"`
	Image     string `json:"image"`
	ColorName string `json:"color_name"`
	ColorCode string `json:"color_code"`
}

type ItemSize struct {
	Sku      string `json:"sku"`
	SizeName string `json:"size_name"`
}


     */

    const classes = useStyles();

    function addImage(axiosResponse) {
        const imagePath = axiosResponse.data
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

    function removeColor(ix) {
        let colorArray = item.colors
        colorArray.splice(ix, 1)
        setItem({...item, colors: colorArray})
    }

    function addColorImage(axiosResponse) {
        const imagePath = axiosResponse.data
        let colorImage = {
            image: imagePath,
            color_name: "",
            color_code: "",
            sku: ""
        }
        let colorImageArray = item.colors == null ? [] : item.colors
        console.log(colorImageArray)
        colorImageArray.push(colorImage)
        setItem({...item, colors: colorImageArray})
    }

    const onFileChange = event => {
        // get signed url from backend
        const formData = new FormData()
        formData.append("content_type", "image/jpeg")
        const file = event.target.files[0]
        formData.append('originalFile', file)
        const targetId = event.target.id
        if (targetId === 'upload-color-picture') {
            formData.append('isColorImage', 'true')
        }
        console.log(targetId)
        axios.post("http://localhost:8082/inventory/" + item.id + "/images", formData)
            .then((response) => {
                if (targetId === 'upload-picture') {
                    addImage(response)
                } else {
                    addColorImage(response)
                }
            })
            .catch((ex) => console.log(ex))
    }

    function deletePicture(imageName, imgIndex) {
        axios.delete("http://localhost:8082/inventory/" + item.id + "/images/" + imageName)
            .then(removeImage(imgIndex))
            .catch((ex) => console.log(ex))
    }

    function deleteColor(itemColor, itmCIndex) {
        console.log(itemColor.image)
        axios.delete("http://localhost:8082/inventory/" + item.id + "/images/" + itemColor.image)
            .then(removeColor(itmCIndex))
            .catch((ex) => console.log(ex))
    }

    const saveItem = async () => {
        if (item.id) {
            await axios.put("http://localhost:8082/inventory/" + item.id, item)
                .then(() => setSaveBtnDisabled(true))
                .catch(x => console.log(x))
        } else {
            await axios.post("http://localhost:8082/inventory", item)
                .then(response => {
                    if (response.status === 201) {
                        let re = /(http:\/\/localhost:8082\/inventory\/)(.*)/
                        item.id = re.exec(response.data)[2]
                        setSaveBtnDisabled(true)
                    } else {
                        console.log("got unexpected response status: " + response.status)
                    }
                })
        }
    };

    const resetItem = () => {
        setItem(itemCopy);
        setSaveBtnDisabled(true)
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
                        <div className={"horizontal-box"}>
                            <div>
                                <EuroSymbolIcon fontSize={"small"}/>:
                            </div>
                            <div>
                                {item.price.amount}
                            </div>
                        </div>
                        <div className={"horizontal-box"}>
                            <div>
                                Quantity:
                            </div>
                            <div>
                                {item.cnt}
                            </div>
                        </div>
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
                                <TextField required id={"description"} value={item.description} label={"Description"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, description: e.target.value});
                                           }}
                                           margin="dense" InputLabelProps={{shrink: true,}} multiline
                                           rows={4}
                                           variant={"outlined"}
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField disabled={item.id} required id={"name"} value={item.name} label={"Name"}
                                           variant={"outlined"}
                                           margin="dense"
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, name: e.target.value});
                                           }}
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            {/*<Grid item xs={4}>
                                <TextField id={"upc"} value={item.upc} label={"Upc"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, upc: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>*/}
                            {/*<Grid item xs={4}>
                                <TextField disabled id={"sku"} value={item.sku} label={"Sku"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, sku: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>*/}
                            <Grid item xs={4}>
                                <TextField disabled={item.id} id={"brand"} required value={item.brand} label={"Brand"}
                                           variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn()
                                               setItem({...item, brand: e.target.value})
                                           }} margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>
                            {/*<Grid item xs={4}>
                                <TextField id={"size"} value={item.size} label={"Sizes"} variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, size: e.target.value});
                                           }}
                                           margin="dense"
                                           style={{width: '100%'}}
                                />
                            </Grid>*/}

                            <Grid item xs={4}>
                                <TextField id={"price"} type={"number"} value={item.price.amount} label={"Price"}
                                           variant={"outlined"}
                                           onChange={e => {
                                               enableSaveBtn();
                                               setItem({...item, price: Number(e.target.value)});
                                           }}
                                           margin="dense"
                                           style={{width: '100%'}}
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
                                           style={{width: '100%'}}
                                />
                            </Grid><Grid item xs={4}>
                            <Paper variant={"outlined"} style={{margin: "5px", padding: "5px"}}>
                                {
                                    item.colors ?
                                        item.colors.map((itemColor, index) => {
                                            return (
                                                <div className={"container"}
                                                     style={{height: "25px", width: "25px"}}>
                                                    <div>
                                                        <img key={index} alt={itemColor.color_name}
                                                             className={"colorImage"}
                                                             src={"http://localhost:8082/inventory/" + item.id + "/images/" + itemColor.image}/>

                                                    </div>


                                                    <div className={"middle"}>
                                                        <Fab style={{height: "25px", width: "25px"}}
                                                             size="small"
                                                             component="div"
                                                             aria-label="remove"
                                                             variant="extended"
                                                        >
                                                            <DeleteIcon style={{height: "20px", width: "auto"}}
                                                                        onClick={event => deleteColor(itemColor, index)}/>
                                                        </Fab>
                                                    </div>

                                                </div>
                                            )
                                        }) : <div/>
                                }
                                <div>
                                    <form>
                                        <label htmlFor={"upload-color-picture"}>
                                            <input style={{display: 'none'}} id={"upload-color-picture"}
                                                   name={"upload-color-picture"}
                                                   type={"file"} onChange={onFileChange}
                                                   aria-label={"Add color picture"}/>
                                            <IconButton size="small" component="div" aria-label="add"
                                                        variant="extended">
                                                <AddIcon/>
                                            </IconButton>
                                        </label>

                                    </form>
                                </div>
                            </Paper>
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
                        <CancelIcon/>
                    </IconButton>
                    <IconButton id={"saveBtn"} aria-label="save" onClick={saveItem} disabled={saveBtnDisabled}>
                        <SaveIcon/>
                    </IconButton>
                </CardFooter>
            </Collapse>
        </Card>
    );
}
