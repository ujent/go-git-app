import { connect } from 'react-redux';

import * as actions from '../actions';
import App from '../App';

const mapStateToProps = (state, ownProps) => {
  return {
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    getUsers: () => dispatch(actions.getUsers()),
    getRepositories: () => dispatch(actions.getRepositories()),
    getBranches: () => dispatch(actions.getBranches())
  };
};
const AppContainer = connect(mapStateToProps, mapDispatchToProps)(App);

export default AppContainer;
