import React, { Component } from 'react';
import classNames from 'classnames';
import '../App.css';

export default class RepoContent extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedFile: '',
            fileNameInput: {
                isVisible: false,
                value: ''
            }
        }
    }

    render() {
        const onFileClick = (name) => {
            this.setState({
                selectedFile: name
            })
        }

        const onFileNameChange = (e) => {

            this.setState({
                fileNameInput: {
                    value: e.target.value,
                    isVisible: true
                }
            })
        }

        const onPlusBtnClick = () => {
            const isVisible = this.state.fileNameInput.isVisible

            this.setState({
                fileNameInput: {
                    isVisible: !isVisible,
                    value: ''
                }
            })
        }

        const onBlurFileName = () => {
            this.setState({
                fileNameInput: {
                    isVisible: false,
                    value: ''
                }
            })
        }

        const files = this.props.files.map(el => {
            const fileClass = classNames({
                'repo-file': true,
                'repo-file-selected': el.path === this.state.selectedFile,
                'conflict-file': el.isConflict
            })

            return <li key={el.path} className="repo-file-wrapper">
                <div className={fileClass} onClick={() => onFileClick(el.path)}>{el.path}</div>
                <button className="add-remove-file minus-btn" onClick={onPlusBtnClick}>-</button>
            </li>
        });

        return (
            <section className="repo-content">
                <h2 className="visually-hidden">Repository content</h2>
                <div className="repo-files">
                    <header className="add-file-wrapper">
                        <h3>Files</h3>
                        <button className="add-remove-file" onClick={onPlusBtnClick}>+</button>
                    </header>
                    {
                        this.state.fileNameInput.isVisible ?
                            <input autoFocus type="text" value={this.state.fileNameInput.value} onBlur={onBlurFileName} onChange={onFileNameChange} />
                            : null
                    }
                    <ul className="files-list">
                        {files}
                    </ul>
                </div>
                {
                    this.state.selectedFile ?
                        <div className="file-content ">
                            <h3>{this.state.selectedFile}</h3>
                            <div className="file-content-buttons">
                                <button type="button" className="button">Save</button>
                            </div>
                            <textarea rows="32"></textarea>
                        </div> : null
                }

            </section>
        );
    };

}