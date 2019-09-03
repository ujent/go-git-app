import { connect } from 'react-redux';

import * as actions from '../actions';
import ConfirmPopup from '../components/ConfirmPopup';

const mapStateToProps = (state, ownProps) => {
    const confirmPopup = state.confirmPopup;
    return {
        isOpen: confirmPopup.isOpen,
        message: confirmPopup.message,
        onClose: confirmPopup.onClose,
        onConfirm: confirmPopup.onConfirm
    };
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        closePopup: () => dispatch(actions.closeConfirm()),
    }
}
const ConfirmPopupContainer = connect(mapStateToProps, mapDispatchToProps)(ConfirmPopup);

export default ConfirmPopupContainer;