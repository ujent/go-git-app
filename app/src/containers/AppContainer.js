import { connect } from 'react-redux';

import * as actions from '../actions';
import App from '../App';

const mapStateToProps = (state, ownProps) => {
  return {
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    getRepositories: () => dispatch(actions.getRepositories()),
    getBranches: () => dispatch(actions.getBranches())
  };
};
const AppContainer = connect(mapStateToProps, mapDispatchToProps)(App);

export default AppContainer;
