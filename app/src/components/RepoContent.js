import React, { Component } from 'react';
import classNames from 'classnames';
import '../App.css';

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

        console.log(this.props.files.length)
        const files = this.props.files.map(el => {
            const fileClass = classNames({
                'repo-file': true,
                'repo-file-selected': this.props.currentFile && el.path === this.props.currentFile.path,
                'conflict-file': el.isConflict
            })

            return <li key={el.path} className="repo-file-wrapper">
                <div className={fileClass} onClick={() => this.onFileClick(el)}>{el.path}</div>
                <button className="add-remove-file minus-btn" onClick={() => this.onRemoveClick(el)}>❌</button>
            </li>
        });

        return (
            <section className="repo-content">
                <h2 className="visually-hidden">Repository content</h2>
                <div className="repo-files">
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
                            <textarea rows="32" value={this.props.currentFile.content} onChange={this.onContentChange}></textarea>
                        </div> : null
                }

            </section>
        );
    };

    onFileClick = (file) => {
        if (file) {
            this.props.handleGetFile(file.path, file.isConflict)
        }
    }

    onSaveClick = () => {
        const f = this.props.currentFile;

        if (f) {
            this.props.handleSaveFile(f.path, f.content, f.isConflict)
        }
    }

    onContentChange = (e) => {
        const f = this.props.currentFile;

        if (f) {
            this.props.handleSetCurrentFile(f.path, e.target.value, f.isConflict)
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

        this.props.handleAddFile(path, '');
    }

    onRemoveClick = (file) => {
        if (file) {
            this.props.handleRemoveFile(file.path, file.isConflict);
        }
    }

}