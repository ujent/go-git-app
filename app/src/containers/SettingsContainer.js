import { connect } from 'react-redux';

import * as actions from '../actions';
import Settings from '../components/Settings';

const mapStateToProps = (state, ownProps) => {
    return {

    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        handleSwitchUser: (user) => dispatch(actions.switchUser(user)),
        handleSelectRepo: (name) => dispatch(actions.switchRepo(name)),
        handleSelectBranch: (name) => dispatch(actions.switchBranch(name))
    };
};
const SettingsContainer = connect(mapStateToProps, mapDispatchToProps)(
    Settings
);

export default SettingsContainer;