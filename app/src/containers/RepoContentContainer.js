import { connect } from 'react-redux';

import * as actions from '../actions';
import RepoContent from '../components/RepoContent';

const mapStateToProps = (state, ownProps) => {
    return {
        files: state.files,
        currentFile: state.currentFile
    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        //handleSelectUser: (user) => dispatch(actions.switchUser(user)),

    };
};
const RepoContentContainer = connect(mapStateToProps, mapDispatchToProps)(
    RepoContent
);

export default RepoContentContainer;