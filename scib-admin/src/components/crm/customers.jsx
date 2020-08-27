import React, {useContext, useState, useEffect} from "react";
import {makeStyles} from "@material-ui/core/styles";
import MaterialTable from "material-table";

const useStyles = makeStyles({
    table: {
        minWidth: 500,
    },
});

export const Customers = () => {

    const classes = useStyles();
    const columns = [
        { title: 'Name', field: 'name' },
        { title: 'Email', field: 'email'},
        { title: 'Addresses', field: 'addresses'},
        { title: 'Phone', field: 'phone'}
    ]

    const [state, setState] = useState({
        data: [
            { name: 'Jimmmy Glitschy', email: 'jimmmy.glitschy@gmail.com', addresses: [{ Id: '42', country: 'Germany', city: 'Hamburg', street: 'Luruper Drift 62', code: '22587' }]}
        ]
    })

    return (
        <MaterialTable columns={columns} data={state.data} editable={
            {
                onRowAdd: (newData) => {
                    new Promise((resolve) => {
                        setTimeout(() => {
                            resolve();
                            setState((prevState) => {
                                const data = [...prevState.data];
                                data.push(newData);
                                return {...prevState, data};
                            });
                        }, 600);
                    })
                },
                onRowUpdate: (newData, oldData) => {

                },
                onRowDelete: (oldData) => {
                    
                }

            }
        }/>
    )
}