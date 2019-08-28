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
    handleCommit: (msg) => dispatch(actions.commit(msg)),
    handleCheckoutBranch: (name) => dispatch(actions.switchBranch(name)),
    handleClone: (url, authName, authPsw) => dispatch(actions.clone(url, authName, authPsw)),
    handleLog: () => dispatch(actions.log()),
    handlePull: (remote, authName, authPsw) => dispatch(actions.pull(remote, authName, authPsw)),
    handlePush: (remote, authName, authPsw) => dispatch(actions.push(remote, authName, authPsw)),
    handleRemoveBranch: (branch) => dispatch(actions.removeBranch(branch)),
    handleRemoveRepo: (repo) => dispatch(actions.removeRepo(repo)),
    handleMerge: () => dispatch(actions.merge())
  };
};
const CommandsContainer = connect(mapStateToProps, mapDispatchToProps)(Commands);

export default CommandsContainer;
