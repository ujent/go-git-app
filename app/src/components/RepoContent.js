import React, { Component } from 'react';
import classNames from 'classnames';
import ReactTooltip from 'react-tooltip'
import '../App.css';
import { FileStatus } from '../constants';

export default class RepoContent extends Component {
    constructor(props) {
        super(props);

        this.state = {
            fileNameInput: {
                isVisible: false,
                value: ''
            },
        }
    }

    render() {

        if (!this.props.isVisible) {
            return null;
        }

        const files = this.props.files.map(el => {
            const fs = this.convertFileStatus(el.fileStatus, el.isConflict);

            if (el.fileStatus === FileStatus.Deleted) {

                return <li key={el.path} className="repo-file-wrapper">
                    <div className="repo-file"><p className="file-status-letter" data-tip={fs.tooltip}>{fs.letter}</p><p className="removed-file">{el.path}</p></div>
                </li>
            }

            const fileClass = classNames({
                'repo-file': true,
                'repo-file-selected': this.props.currentFile && el.path === this.props.currentFile.path,
                'conflict-file': el.isConflict,
            });

            const letterClass = classNames({
                'file-status-letter': true,
                'conflict-file-letter': el.isConflict,
            });

            return <li key={el.path} className="repo-file-wrapper">
                <div className={fileClass} onClick={() => this.onFileClick(el)}><p className={letterClass} data-tip={fs.tooltip}>{fs.letter}</p><p>{el.path}</p></div>
                <button className="add-remove-file remove-btn" onClick={() => this.onRemoveClick(el)}><span role="img" aria-label="remove">✕</span></button>
            </li>
        });

        return (
            <section className="repo-content">
                <h2 className="visually-hidden">Repository content</h2>
                <div className="repo-files" onKeyPress={e => this.handleEnterPress(e)}>
                    <header className="add-file-wrapper">
                        <h3>Files</h3>
                        <button className="add-remove-file" onClick={this.onAddClick}>+</button>
                    </header>
                    {
                        this.state.fileNameInput.isVisible ?
                            <input autoFocus type="text" value={this.state.fileNameInput.value} onBlur={this.onBlurFileName} onChange={this.onFileNameChange} />
                            : null
                    }
                    <ul className="files-list">
                        {files}
                    </ul>
                </div>
                {
                    this.props.currentFile ?
                        <div className="file-content ">
                            <h3>{this.props.currentFile.path}</h3>
                            <div className="file-content-buttons">
                                <button type="button" className="button" onClick={this.onSaveClick}>Save</button>
                            </div>
                            <textarea rows="34" value={this.props.currentFile.content} onChange={this.onContentChange}></textarea>
                        </div> : null
                }
                <ReactTooltip place="top" type="info" effect="solid" getContent={(c) => { if (c) return c; return null }} />
            </section>
        );
    };

    convertFileStatus = (fs, isConflict) => {
        if (isConflict) {
            return { letter: 'C', tooltip: 'Conflict' }
        }

        switch (fs) {
            case FileStatus.Unspecified:
                return { letter: '' }
            case FileStatus.Unmodified:
                return { letter: '' }
            case FileStatus.Modified:
                return { letter: 'M', tooltip: 'Modified' }
            case FileStatus.Added:
                return { letter: 'A', tooltip: 'Added' }
            case FileStatus.Deleted:
                return { letter: 'D', tooltip: 'Deleted' }
            case FileStatus.Untracked:
                return { letter: '?', tooltip: 'Untracked' }
            case FileStatus.Renamed:
                return { letter: 'R', tooltip: 'Renamed' }
            case FileStatus.Copied:
                return { letter: 'C', tooltip: 'Copied' }
            case FileStatus.UpdatedButUnmerged:
                return { letter: 'U', tooltip: 'UpdatedButUnmerged' }
            default:
                return { letter: '??', tooltip: `Wrong status: ${fs}` }
        }
    }

    onFileClick = (file) => {
        if (file) {
            this.props.handleGetFile(file.path, file.isConflict, file.fileStatus)
        }
    }

    onSaveClick = () => {
        const f = this.props.currentFile;

        if (f) {
            this.props.handleSaveFile(f.path, f.content)
        }
    }

    onContentChange = (e) => {
        const f = this.props.currentFile;

        if (f) {
            this.props.handleSetCurrentFile(f.path, e.target.value, f.isConflict, f.fileStatus)
        }
    }

    onFileNameChange = (e) => {

        this.setState({
            fileNameInput: {
                value: e.target.value,
                isVisible: true
            }
        });
    }

    onAddClick = () => {
        const isVisible = this.state.fileNameInput.isVisible

        this.setState({
            fileNameInput: {
                isVisible: !isVisible,
                value: ''
            }
        });
    }

    onBlurFileName = () => {

        const path = this.state.fileNameInput.value;

        this.setState({
            fileNameInput: {
                isVisible: false,
                value: ''
            }
        });

        if (this.state.fileNameInput.value) {
            this.props.handleAddFile(path, '');
        }
    }

    handleEnterPress = (e) => {

        if (e.key === 'Enter') {
            e.preventDefault();

            this.onBlurFileName();
        }
    }

    onRemoveClick = (file) => {
        if (file) {
            const msg = `Do you really want to remove ${file.path}?`;
            const onClose = function () { }
            const onConfirm = () => this.props.handleRemoveFile(file.path);

            this.props.handleOpenConfirm(msg, onConfirm, onClose);
        }
    }

}