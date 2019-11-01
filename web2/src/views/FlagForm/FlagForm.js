import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import { Redirect, useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { FlagDetails } from './components';
import {
  CREATE_FLAG_RULE_QUERY,
  CREATE_VARIANT_QUERY,
  DELETE_FLAG_QUERY,
  DELETE_FLAG_RULE_QUERY,
  DELETE_VARIANT_QUERY,
  FLAG_QUERY,
  UPDATE_FLAG_QUERY,
  UPDATE_FLAG_RULE_QUERY,
  UPDATE_VARIANT_QUERY,
} from './queries';
import { formatFlag, formatRule, formatVariant, newFlag } from './models';
import { FLAGS_QUERY } from '../FlagList/queries';
import { reject } from 'lodash';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
}));

const FlagForm = () => {
  const { id } = useParams();
  const [toFlagsPage, setToFlagsPage] = React.useState(false);
  const { loading, error, data } = useQuery(FLAG_QUERY, { variables: { id } });
  const [deleteFlag] = useMutation(DELETE_FLAG_QUERY, {
    update(cache, { data: { deleteFlag: id } }) {
      const { flags } = cache.readQuery({ query: FLAGS_QUERY });
      cache.writeQuery({
        query: FLAGS_QUERY,
        data: { flags: reject(flags, { id }) },
      });
    },
  });
  const [updateFlag] = useMutation(UPDATE_FLAG_QUERY);
  const [createVariant] = useMutation(CREATE_VARIANT_QUERY);
  const [updateVariant] = useMutation(UPDATE_VARIANT_QUERY);
  const [deleteVariant] = useMutation(DELETE_VARIANT_QUERY);
  const [createRule] = useMutation(CREATE_FLAG_RULE_QUERY);
  const [updateRule] = useMutation(UPDATE_FLAG_RULE_QUERY);
  const [deleteRule] = useMutation(DELETE_FLAG_RULE_QUERY);
  const handleSaveFlag = async (flag, deletedItems) => {
    if (flag.__changed) {
      await updateFlag({ variables: { id: flag.id, input: formatFlag(flag) } });
    }
    //TODO: fix scenario where rule references a new variant
    await Promise.all([
      ...flag.variants.map(variant => {
        const variables = {
          id: variant.id,
          flagId: flag.id,
          input: formatVariant(variant),
        };
        if (variant.__new) {
          return createVariant({ variables });
        }
        if (variant.__changed) {
          return updateVariant({ variables });
        }
      }),
      ...flag.rules.map(rule => {
        const variables = {
          id: rule.id,
          flagId: flag.id,
          input: formatRule(rule),
        };
        if (rule.__new) {
          return createRule({ variables });
        }
        if (rule.__changed) {
          return updateRule({ variables });
        }
      }),
      ...deletedItems.map(item => {
        switch (item.type) {
          case 'variant':
            return deleteVariant({ variables: item });
          case 'rule':
            return deleteRule({ variables: item });
        }
      }),
    ]);
    setToFlagsPage(true);
  };
  useEffect(() => {
    const handleEsc = (event) => {
      if (event.key === 'Escape') setToFlagsPage(true);
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);
  const classes = useStyles();
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;
  const handleDeleteFlag = id => {
    deleteFlag({ variables: { id } })
      .then(() => setToFlagsPage(true));
  };

  return (
    <div className={classes.root}>
      {toFlagsPage && <Redirect to='/flags'/>}
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <FlagDetails
            flag={newFlag(data.flag)}
            operations={data.operations.enumValues.map(v => v.name)}
            segments={data.segments}
            onSaveFlag={handleSaveFlag}
            onDeleteFlag={handleDeleteFlag}
          />
        </Grid>
      </Grid>
    </div>
  );
};

export default FlagForm;
