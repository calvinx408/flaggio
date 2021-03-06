import React from 'react';
import PropTypes from 'prop-types';
import {
  Button,
  FormControl,
  Grid,
  Hidden,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Tooltip,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/styles';
import RemoveIcon from '@material-ui/icons/RemoveCircleOutline';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { VariantTypes } from '../../../../helpers';
import { BooleanType, VariantType } from '../../copy';

const useStyles = makeStyles(theme => ({
  formControl: {
    fullWidth: true,
    display: 'flex',
    wrap: 'nowrap',
  },
  sideButton: {
    minWidth: theme.spacing(0),
  },
}));

const VariantFields = ({ variant, showAddButton, onAddVariant, onUpdateVariant, onDeleteVariant }) => {
  const classes = useStyles();
  return (
    <Grid container spacing={1}>
      <Grid item sm={3} xs={5}>
        <FormControl
          className={classes.formControl}
          margin="dense"
          variant="outlined"
        >
          <InputLabel>Type</InputLabel>
          <Select
            value={variant.type}
            name="type"
            required
            onChange={e => {
              onUpdateVariant(e);
              onUpdateVariant({ target: { name: 'value', value: '' } });
            }}
            labelWidth={30}
          >
            {Object.keys(VariantTypes).map(type => (
              <MenuItem key={type} value={VariantTypes[type]}>{VariantType[type]}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item sm={4} xs={5}>
        {
          variant.type === VariantTypes.BOOLEAN ?
            (
              <FormControl
                className={classes.formControl}
                margin="dense"
                variant="outlined"
                required
              >
                <InputLabel>Value</InputLabel>
                <Select
                  value={variant.value}
                  name="value"
                  onChange={onUpdateVariant}
                  labelWidth={45}
                >
                  {[true, false].map(val => (
                    <MenuItem key={val} value={val}>{BooleanType[val]}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            ) :
            (
              <TextField
                label="Value"
                margin="dense"
                name="value"
                value={variant.value}
                required
                type={variant.type === VariantTypes.NUMBER ? 'number' : 'text'}
                onChange={onUpdateVariant}
                fullWidth
                variant="outlined"
              />
            )
        }
      </Grid>
      <Hidden xsDown>
        <Grid item xs={4}>
          <TextField
            label="Description"
            margin="dense"
            name="description"
            value={variant.description}
            onChange={onUpdateVariant}
            fullWidth
            variant="outlined"
          />
        </Grid>
      </Hidden>
      <Grid item sm={1} xs={2} style={{ display: 'flex' }}>
        <Tooltip title="Delete variant" placement="top">
          <Button
            size="small"
            color="secondary"
            className={classes.sideButton}
            onClick={onDeleteVariant}
          >
            <RemoveIcon/>
          </Button>
        </Tooltip>
        {
          showAddButton ? (
            <Tooltip title="New variant" placement="top">
              <Button
                size="small"
                color="primary"
                className={classes.sideButton}
                onClick={onAddVariant}
              >
                <AddIcon/>
              </Button>
            </Tooltip>
          ) : (
            <Button
              size="small"
              color="primary"
              className={classes.sideButton}
              disabled
            >
              <AddIcon style={{ visibility: 'hidden' }}/>
            </Button>
          )
        }
      </Grid>
    </Grid>
  )
};

VariantFields.propTypes = {
  variant: PropTypes.object.isRequired,
  showAddButton: PropTypes.bool.isRequired,
  onAddVariant: PropTypes.func.isRequired,
  onUpdateVariant: PropTypes.func.isRequired,
  onDeleteVariant: PropTypes.func.isRequired,
};

export default VariantFields;