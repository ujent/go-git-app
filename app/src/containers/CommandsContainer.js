import { connect } from 'react-redux';

import * as actions from '../actions';
import Commands from '../components/Commands';

const mapStateToProps = (state, ownProps) => {
  return {
    outputMsg: state.output
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    getRepositories: () => dispatch(actions.getRepositories()),
    getBranches: () => dispatch(actions.getBranches())
  };
};
const CommandsContainer = connect(mapStateToProps, mapDispatchToProps)(Commands);

export default CommandsContainer;
