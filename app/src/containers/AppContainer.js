import { connect } from 'react-redux';

import * as actions from '../actions';
import App from '../App';
import { FileStatus } from '../constants'

const mapStateToProps = (state, ownProps) => {
  const changed = state.files.filter(e => e.fileStatus !== FileStatus.Unmodified);
  return {
    outputMsg: state.output,
    confirm: state.confirmPopup,
    isSpinnerVisible: state.isSpinnerVisible,
    hasUncommitted: changed.length > 0
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    getRepositories: () => dispatch(actions.getRepositories()),
    getBranches: () => dispatch(actions.getBranches()),
    handleCloseConfirm: () => dispatch(actions.closeConfirm()),
  };
};
const AppContainer = connect(mapStateToProps, mapDispatchToProps)(App);

export default AppContainer;
