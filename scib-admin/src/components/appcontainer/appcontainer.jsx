import React from 'react';
import PropTypes from 'prop-types';
import {makeStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import Inventory from "../inventory/inventory";
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import {InventoryProvider} from "../inventory/inventoryContext";


function TabPanel(props) {
    const {children, value, index, ...other} = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box p={3}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}

TabPanel.propTypes = {
    children: PropTypes.node,
    index: PropTypes.any.isRequired,
    value: PropTypes.any.isRequired,
};

function a11yProps(index) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
    },
}));

export default function Appcontainer() {

    const classes = useStyles();
    const [value, setValue] = React.useState(0);

    const handleChange = (event, newValue) => {
        setValue(newValue);
    };


    return (
        <div className={classes.root}>
            <AppBar position={"static"}>
                <Tabs value={value} onChange={handleChange} aria-label={"scib navigayion"}>
                    <Tab label="Inventory" {...a11yProps(0)}/>
                    <Tab label="Shipping" {...a11yProps(1)}/>
                    <Tab label="CRM" {...a11yProps(2)}/>
                    <Tab label="Suppliers" {...a11yProps(3)}/>
                    <Tab label="Reporting" {...a11yProps(4)}/>
                </Tabs>
            </AppBar>
            <TabPanel value={value} index={0}>
                <Typography component={"div"}>
                    <InventoryProvider>
                        <Inventory/>
                    </InventoryProvider>
                </Typography>
            </TabPanel>
            <TabPanel value={value} index={1}>
                <Shipping/>
            </TabPanel>
            <TabPanel value={value} index={2}>
                <Customers/>
            </TabPanel>
            <TabPanel value={value} index={3}>
                <Suppliers/>
            </TabPanel>
            <TabPanel value={value} index={4}>
                <Reporting/>
            </TabPanel>
        </div>
    )

}

function Shipping() {
    return <span>Shipments</span>
}

function Customers() {
    return <span>Customer Relationship Management</span>
}

function Suppliers() {
    return <span>Supplier Relations</span>
}

function Reporting() {
    return <span>Reports</span>
}