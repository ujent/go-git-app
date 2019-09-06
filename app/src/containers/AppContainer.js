import { connect } from 'react-redux';

import * as actions from '../actions';
import App from '../App';

const mapStateToProps = (state, ownProps) => {
  return {
    outputMsg: state.output,
    confirm: state.confirmPopup,
    isSpinnerVisible: state.isSpinnerVisible
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
