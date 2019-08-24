import { connect } from 'react-redux';

import * as actions from '../actions';
import Commands from '../components/Commands';

const mapStateToProps = (state, ownProps) => {
  return {
    repositories: state.repositories,
    branches: state.branches,
    currentUser: state.settings.currentUser,
    currentRepo: state.settings.currentRepo,
    currentBranch: state.settings.currentBranch
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    handleCommit: (msg) => dispatch(actions.commit()),
    handleCheckoutBranch: () => dispatch(actions.checkoutBranch()),
    handleClone: () => dispatch(actions.clone()),
    handleLog: () => dispatch(actions.log()),
    handleMerge: () => dispatch(actions.merge()),
    handlePull: () => dispatch(actions.pull()),
    handlePush: () => dispatch(actions.push()),
    handleRemoveBranch: () => dispatch(actions.removeBranch()),
    handleRemoveRepo: () => dispatch(actions.removeRepo()),
  };
};
const CommandsContainer = connect(mapStateToProps, mapDispatchToProps)(Commands);

export default CommandsContainer;
