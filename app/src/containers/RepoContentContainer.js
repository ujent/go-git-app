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
        handleSaveFile: (path, content, isConflict) => dispatch(actions.saveFile(path, content, isConflict)),
        handleAddFile: (path, content) => dispatch(actions.addFile(path, content)),
        handleRemoveFile: (path, isConflict) => dispatch(actions.removeFile(path, isConflict)),
        handleGetFile: (path, isConflict) => dispatch(actions.getFile(path, isConflict)),
        handleSetCurrentFile: (path, content, isConflict) => dispatch(actions.setCurrentFile(path, content, isConflict)),
    };
};
const RepoContentContainer = connect(mapStateToProps, mapDispatchToProps)(
    RepoContent
);

export default RepoContentContainer;