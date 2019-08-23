import { connect } from 'react-redux';

import * as actions from '../actions';
import Settings from '../components/Settings';

const mapStateToProps = (state, ownProps) => {
    return {
        users: state.users,
        selectedUser: state.settings.currentUser,
        repos: state.repositories,
        selectedRepo: state.settings.currentRepo,
        branches: state.branches,
        selectedBranch: state.settings.currentBranch
    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        handleSelectUser: (user) => dispatch(actions.switchUser(user)),
        handleSelectRepo: (name) => dispatch(actions.switchRepo(name)),
        handleSelectBranch: (name) => dispatch(actions.switchBranch(name))
    };
};
const SettingsContainer = connect(mapStateToProps, mapDispatchToProps)(
    Settings
);

export default SettingsContainer;