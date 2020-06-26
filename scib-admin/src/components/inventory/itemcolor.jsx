import React, {useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Fab from "@material-ui/core/Fab";
import DeleteIcon from "@material-ui/icons/Delete";


export const ItemColor = (props) => {

    const [itemColor, setItemColor] = useState(props.itemColor)
    const deleteColor = null


    return (
        <div className={"container"}>
            <img key={index} alt={itemColor.color_name} className={"colorImage"}
                 src={"http://localhost:8082/inventory/" + item.id + "/images/" + itemColor.image}
            />
            <div className={"middle"}>
                <Fab
                    color="primary"
                    size="small"
                    component="div"
                    aria-label="remove"
                    variant="extended"
                >
                    <DeleteIcon onClick={event => deleteColor(itemColor, index)}/>
                </Fab>
            </div>
        </div>
    )

}