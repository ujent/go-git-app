import { connect } from 'react-redux';

import * as actions from '../actions';
import RepoContent from '../components/RepoContent';

const mapStateToProps = (state, ownProps) => {
    return {
        files: state.files,
        currentFile: state.currentFile,
    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        handleSaveFile: (path, content) => dispatch(actions.saveFile(path, content)),
        handleAddFile: (path, content) => dispatch(actions.addFile(path, content)),
        handleRemoveFile: (path) => dispatch(actions.removeFile(path)),
        handleGetFile: (path, isConflict) => dispatch(actions.getFile(path, isConflict)),
        handleSetCurrentFile: (path, content, isConflict) => dispatch(actions.setCurrentFile(path, content, isConflict)),
        handleOpenConfirm: (msg, onConfirm, onClose) => dispatch(actions.showConfirm(msg, onConfirm, onClose))
    };
};
const RepoContentContainer = connect(mapStateToProps, mapDispatchToProps)(
    RepoContent
);

export default RepoContentContainer;